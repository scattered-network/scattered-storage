package rbd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/language"
	"github.com/scattered-network/scattered-storage/lib/validators"
)

// ListBlock
/* lsblk -J -o NAME,PATH,MOUNTPOINT,FSTYPE

Example output follows.

{
   "blockdevices": [
      {"name":"rbd0p1", "path":"/dev/rbd0p1", "mountpoint":null, "fstype":null}
      {"name":"rbd1p1", "path":"/dev/rbd1p1", "mountpoint":null, "fstype":null},
      {"name":"rbd2p1", "path":"/dev/rbd2p1", "mountpoint":null, "fstype":null}
   ]
}

ListBlock is used to find mounted RBD images and related mount points. */
type ListBlock struct {
	Blockdevices []*struct {
		Name       string `json:"name"`
		Path       string `json:"path"`
		Mountpoint string `json:"mountpoint"`
		FSType     string `json:"fstype"`
		Children   []*struct {
			Name       string `json:"name"`
			Path       string `json:"path"`
			Mountpoint string `json:"mountpoint"`
			FSType     string `json:"fstype"`
		} `json:"children,omitempty"`
	} `json:"blockdevices"`
}

// findDevicePath Goes through the list of mapped RBD images to verify the mapping of the image.
// Once verified, the lsblk command is executed which returns the device mount information.
func (c *RBDClient) findDevicePath(pool string, name string) *ListBlock {
	log.Trace().Str("Pool", pool).Str("Name", name).Msg("findDevicePath")

	list, showMappedError := ListMappedImages()
	if showMappedError != nil {
		return &ListBlock{Blockdevices: nil}
	}

	var listError error

	deviceMountInfo := &ListBlock{Blockdevices: nil}

	for _, image := range *list {
		if image.Name != name || image.Pool != pool {
			log.Trace().Interface("Image", image).Msgf("Skipping executeListBlock(%s)", image.Device)

			continue
		}

		log.Info().Interface("Image", image).Msgf("Image %s/%s Matches Request", image.Pool, image.Name)

		deviceMountInfo, listError = c.executeListBlock(image.Device)
		log.Trace().Interface("deviceMountInfo", deviceMountInfo)

		if listError != nil {
			log.Trace().Str("Error", listError.Error()).Msgf("executeListBlock(%s)", image.Device)
		}

		if len(deviceMountInfo.Blockdevices[0].Children) > 0 {
			if deviceMountInfo.Blockdevices[0].Children[0].Mountpoint == "" {
				time.Sleep(5 * time.Second)
			} else {
				continue
			}
		} else {
			continue
		}
	}

	return deviceMountInfo
}

// executeListBlock returns a listing of block devices on the server generated using lsblk.
func (c *RBDClient) executeListBlock(device string) (*ListBlock, error) {
	log.Trace().Str("Device", device).Msg("executeListBlock")

	if !ValidateDevicePath(device) {
		log.Trace().Str("Device", device).Msg("Device Path Invalid")

		return &ListBlock{Blockdevices: nil}, validators.ErrInvalidDevicePath
	}

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "lsblk", "-J", device, "-o", "NAME,PATH,MOUNTPOINT,FSTYPE")
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return &ListBlock{Blockdevices: nil}, fmt.Errorf("ERROR: lsblk failed:\n%w", err)
	}

	var list *ListBlock

	if err := json.Unmarshal(stdOut.Bytes(), &list); err != nil {
		return &ListBlock{Blockdevices: nil}, fmt.Errorf(
			"ERROR: json for lsblk could not unmarshal:\n%w\n%s", err, stdOut.String(),
		)
	}

	log.Trace().Interface("list", list).Msg("List of Block Devices Found")

	return list, nil
}
