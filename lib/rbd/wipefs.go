package rbd

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/validators"
)

// WipeFS
/* wipefs -J -n /dev/rbd0

Example output follows:

{
   "signatures": [
      {"device":"rbd0", "offset":"0x0", "type":"xfs", "uuid":"976330da-8105-4514-b4c6-8914fcd8e6d3", "label":null},
      {"device":"rbd0", "offset":"0x200", "type":"gpt", "uuid":null, "label":null},
      {"device":"rbd0", "offset":"0x13ffffe00", "type":"gpt", "uuid":null, "label":null},
      {"device":"rbd0", "offset":"0x1fe", "type":"PMBR", "uuid":null, "label":null}
   ]
}
WipeFS is used to discover filesystem and signatures on a device. */
type WipeFS struct {
	Signatures []*struct {
		Device string      `json:"device"`
		Offset string      `json:"offset"`
		Type   string      `json:"type"`
		UUID   *string     `json:"uuid"`
		Label  interface{} `json:"label"`
	} `json:"signatures"`
}

// hasPartitions returns a boolean value representing success in finding any partitions.
func (c *RadosBlockDeviceClient) hasPartitions(device string) (bool, error) {
	if !ValidateDevicePath(device) {
		log.Trace().Str("Device", device).Msg("could not validate device path")

		return false, fmt.Errorf("%w", validators.ErrInvalidDevicePath)
	}

	log.Trace().Str("Device", device).Msg("hasPartitions")

	deviceInfo, listError := c.executeListBlock(device)

	if listError != nil {
		return false, listError
	}

	if len(deviceInfo.Blockdevices[0].Children) > 0 {
		log.Info().Msg("Partitions found")

		return true, nil
	}

	return false, nil
}

// hasSupportedFileSystem returns a boolean value representing success in finding a supported signature.
func (c *RadosBlockDeviceClient) hasSupportedFileSystem(device string) bool {
	if !ValidateDevicePath(device) {
		log.Trace().Str("Device", device).Msg("could not validate device path")

		return false
	}

	log.Trace().Str("Device", device).Msg("hasSupportedFileSystem")

	if deviceCheck, fsCheckError := c.executeWipeFSWithoutAction(device); fsCheckError != nil {
		for _, signature := range deviceCheck.Signatures {
			if c.isValidFilesystemType(signature.Type) {
				return true
			}
		}
	}

	return false
}

func (c *RadosBlockDeviceClient) isValidFilesystemType(fsType string) bool {
	switch fsType {
	case "xfs":
		return true
	case "ext4":
		return true
	}

	return false
}
