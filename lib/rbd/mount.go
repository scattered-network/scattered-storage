package rbd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/language"
	"github.com/scattered-network/scattered-storage/lib/validators"
)

var ErrMountFailed = errors.New("rbd could not be mounted")

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Mount will execute the mapping and mounting of a given RBD image.
func Mount(pool, name, path string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	if exists, err := PathExists(path); err != nil {
		return err
	} else {
		if exists {
			log.Trace().Str("Pool", pool).Str("Name", name).Msg("Mount")
			client := &RBDClient{}
			return client.executeMount(pool, name, path)
		}
	}
	return nil
}

// executeMount performs the mapping, formatting, and mounting of an RBD image on the server.
func (c *RBDClient) executeMount(pool, name, path string) error {
	log.Trace().Msg("executeMount")

	device, mapped := c.isMapped(pool, name)
	if !mapped {
		if err := c.executeRBDMap(pool, name); err != nil {
			return err
		}

		device, _ = c.isMapped(pool, name)
	}

	partitionsExist, partitionCheckError := c.hasPartitions(device)

	if partitionCheckError != nil {
		return partitionCheckError
	}

	partitionPath := device + "p1"

	if !partitionsExist {
		if partitionError := PartitionEntireDisk(device); partitionError != nil {
			return partitionError
		}

		if probeError := Partprobe(device); probeError != nil {
			return probeError
		}

		client := &RBDClient{}
		if makeFSError := c.makeFilesystem(partitionPath, client.getFilesystemOptionDefaults("xfs")); makeFSError != nil {
			return makeFSError
		}

		if probeError := Partprobe(device); probeError != nil {
			return probeError
		}
	}

	if err := os.MkdirAll(path, 0o701); err != nil {
		log.Error().Str("Path", path).Str("Error", err.Error()).Msg("could not create directory")

		return fmt.Errorf("%w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "mount", partitionPath, path)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	if err := cmd.Start(); err != nil {
		log.Error().Interface("Error", err).Msg("Error starting mkfs")
	}

	err := cmd.Wait()
	if err != nil {
		var exitError *exec.ExitError
		if ok := errors.Is(err, exitError); ok {
			if exitError.ExitCode() == 32 {
				log.Trace().Str("device", partitionPath).Str("mount", path).Msg("device is already mounted")
			} else {
				log.Error().Interface("command", cmd).Interface("error", err).Msg("error during mount")
			}
		}
	}

	return nil
}

// getMountPoint returns the path where a given RBD image is currently mounted.
func (c *RBDClient) getMountPoint(pool string, name string) (string, error) {
	if !ValidatePool(pool) {
		return "", validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return "", validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("getMountPoint")

	deviceMountInfo := c.findDevicePath(pool, name)

	log.Trace().Interface("deviceMountInfo", deviceMountInfo).Msg("Device Information Found")

	return c.findMount(deviceMountInfo), nil
}

// findMount returns the path where an RBD image has been mounted.
func (c *RBDClient) findMount(deviceMountInfo *ListBlock) string {
	log.Trace().Interface("deviceMountInfo", deviceMountInfo).Msg("Executing findMount")

	if deviceMountInfo == nil {
		log.Error().Msg("deviceMountInfo is nil")

		return ""
	}

	for _, partition := range deviceMountInfo.Blockdevices {
		if partition.Mountpoint != "" {
			log.Info().Str("mount", partition.Mountpoint).Msg("mount point found")

			return partition.Mountpoint
		} else {
			log.Trace().Interface("parent", partition.Mountpoint).Msg("not using this mountpount")
		}

		for _, child := range partition.Children {
			if child.Mountpoint != "" {
				log.Info().Str("child", child.Mountpoint).Msg("mount point found")

				return child.Mountpoint
			} else {
				log.Trace().Interface("child", child.Mountpoint).Msg("not using this mountpount")
			}
		}
	}

	return ""
}
