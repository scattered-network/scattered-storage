package rbd

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/language"
	"github.com/scattered-network/scattered-storage/lib/validators"
)

func (c *RadosBlockDeviceClient) Partprobe(device string) error {
	log.Trace().Msg("Partprobe")

	if !ValidateDevicePath(device) {
		return validators.ErrInvalidDevicePath
	}

	if probeError := c.executePartprobe(device); probeError != nil {
		return probeError
	}

	return nil
}

func (c *RadosBlockDeviceClient) executePartprobe(device string) error {
	log.Trace().Msg("executePartprobe")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "partprobe", device)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ERROR: partprobe failed: %w", err)
	}

	return nil
}
