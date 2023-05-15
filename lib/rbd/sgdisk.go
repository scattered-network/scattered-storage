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

func PartitionEntireDisk(device string) error {
	if !ValidateDevicePath(device) {
		return validators.ErrInvalidDevicePath
	}

	log.Trace().Str("Device", device).Msg("PartitionEntireDisk")

	client := &RBDClient{}
	if clearError := client.executeClearPartitions(device); clearError != nil {
		return clearError
	}

	if partitionError := client.executePartitionEntireDisk(device); partitionError != nil {
		return partitionError
	}

	return nil
}

func (c *RBDClient) executeClearPartitions(device string) error {
	log.Trace().Str("Device", device).Msg("executeClearPartitions")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sgdisk", "-o", device)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ERROR: sgdisk clear failed: %w", err)
	}

	return nil
}

func (c *RBDClient) executeZapPartitions(device string) error {
	log.Trace().Str("Device", device).Msg("executeZapPartitions")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sgdisk", "--zap", device)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ERROR: sgdisk zap failed: %w", err)
	}

	return nil
}

func (c *RBDClient) executePartitionEntireDisk(device string) error {
	log.Trace().Str("Device", device).Msg("executePartitionEntireDisk")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sgdisk", "--new", "1::0", "--typecode", "1:8300", device)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ERROR: sgdisk new failed: %w", err)
	}

	return nil
}
