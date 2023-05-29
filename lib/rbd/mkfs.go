package rbd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/language"
	"github.com/scattered-network/scattered-storage/lib/validators"
	"github.com/spf13/cast"
)

type MkfsOptions struct {
	Options map[string]*MkfsOption
}

type MkfsOption struct {
	Value interface{}
}

const (
	TagXfs                        = "xfs"
	XFSDefaultFilesystemBlockSize = 4096
	TagExt4                       = "ext4"
	fsTypeKey                     = "fsType"
	noDiscardKey                  = "noDiscard"
	filesystemBlockSizeKey        = "filesystemBlockSize"
)

// makeFilesystem takes a device path and a set of *MkfsOptions to execute the mkfs command.
func (c *RadosBlockDeviceClient) makeFilesystem(device string, fsOptions *MkfsOptions) error {
	if !ValidateDevicePath(device) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateMakeFilesystemOptions(fsOptions) {
		return validators.ErrInvalidMakeOptions
	}

	log.Trace().Str("Device", device).Interface("Options", fsOptions).Msg("makeFilesystem")
	return c.executeMakeFilesystem(device, fsOptions)
}

// getFilesystemOptionDefaults returns the suggested default options for a supported filesystem.
// Currently, only XFS and EXT4 support is enabled.
func (c *RadosBlockDeviceClient) getFilesystemOptionDefaults(fsType string) *MkfsOptions {
	log.Trace().Str("FsType", fsType).Msg("getFilesystemOptionDefaults")

	xfs := &MkfsOptions{
		Options: map[string]*MkfsOption{
			fsTypeKey: {
				Value: TagXfs,
			},
			filesystemBlockSizeKey: {
				Value: XFSDefaultFilesystemBlockSize, // The default value of an XFS filesystem is 4096
			},
			noDiscardKey: {
				Value: true,
			},
		},
	}

	ext4 := &MkfsOptions{
		Options: map[string]*MkfsOption{
			fsTypeKey: {
				Value: TagExt4,
			},
			noDiscardKey: {
				Value: true,
			},
		},
	}

	switch fsType {
	case TagXfs:
		return xfs
	case TagExt4:
		return ext4
	}

	return nil
}

// executeMakeFilesystem executes the mkfs -t <filesystem> command and includes support for
// the filesystem type and the no discard option (for use with very large images).
func (c *RadosBlockDeviceClient) executeMakeFilesystem(device string, fsOptions *MkfsOptions) error {
	log.Info().Str("Device", device).Interface("FsOptions", fsOptions).Msg("executeMakeFilesystem")

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	fsType := cast.ToString(fsOptions.Options["fsType"].Value)

	filesystemBlockSize := cast.ToString(fsOptions.Options["filesystemBlockSize"].Value)
	blockSize := "" // will use xfs default block size if empty

	if fsType == TagXfs {
		if filesystemBlockSize != "" {
			blockSize = "size=" + filesystemBlockSize
		}
	}

	noDiscard := cast.ToBool(fsOptions.Options["noDiscard"].Value)
	discard := "" // discard will happen if empty

	if noDiscard {
		discard = "-K" // do not discard
	}

	cmd := exec.CommandContext(ctx, "mkfs."+TagXfs, "-b", blockSize, discard, device)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Start(); err != nil {
		log.Error().Interface("Error", err).Msg("Error starting mkfs")
	}

	err := cmd.Wait()
	if err != nil {
		log.Error().Interface("Command", cmd).Interface("Error", err).Msgf("Error During mkfs.%s", fsType)

		return fmt.Errorf("%w", err)
	}

	return nil
}
