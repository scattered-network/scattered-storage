package rbd

import (
	"bytes"
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

// Unmount will execute the umount and unmap of a given RBD image.
func (c *RadosBlockDeviceClient) Unmount(pool, name string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("Unmount")

	unmountError := c.executeUnmount(pool, name)
	if unmountError != nil {
		return unmountError
	}

	return nil
}

// executeUnmount runs the umount -A command against the partition a device has mounted returns nil error on success.
func (c *RadosBlockDeviceClient) executeUnmount(pool, name string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("executeUnmount")

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	device := c.findDevicePath(pool, name)

	if device.Blockdevices != nil {
		partition := device.Blockdevices[0].Children[0].Path
		cmd := exec.CommandContext(ctx, "umount", "-A", partition)
		log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

		cmd.Stdout = &stdOut
		cmd.Stderr = &stdErr

		if err := cmd.Run(); err != nil {
			var exitError *exec.ExitError
			if ok := errors.Is(err, exitError); ok {
				if exitError.ExitCode() != 1 || exitError.ExitCode() != 0 {
					log.Trace().Str("Error", err.Error()).Int(
						"ExitCode", exitError.ExitCode(),
					).Msg("umount non-zero/non-one exit code")

					return fmt.Errorf("ERROR: umount failed: %w", err)
					// return fmt.Errorf("%w", err)
				}
			}
		}
	}

	return nil
}

// Unmap will find the device path for a given image and unmap it from the server.
func (c *RadosBlockDeviceClient) Unmap(pool, name string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidDevicePath
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("Unmap")

	deviceMountInfo := c.findDevicePath(pool, name)
	if len(deviceMountInfo.Blockdevices) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "no block device path found (for: %s), skipping umount", name)
		return nil
	}
	device := deviceMountInfo.Blockdevices[0].Path

	if unmountError := c.executeUnmap(device); unmountError != nil {
		return unmountError
	}

	return nil
}

// executeUnmap runs the rbd unmap command and returns nil error on success.
func (c *RadosBlockDeviceClient) executeUnmap(device string) error {
	if !ValidateDevicePath(device) {
		return validators.ErrInvalidDevicePath
	}

	log.Trace().Str("Device", device).Msg("executeUnmap")

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rbd", "unmap", device)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ERROR: rbd unmap failed: %w", err)
	}

	return nil
}
