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
)

func (c *RadosBlockDeviceClient) DeleteRBD(pool string, name string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	client := &RadosBlockDeviceClient{}
	if deleteError := client.executeRBDDelete(pool, name); deleteError != nil {
		return fmt.Errorf("%w", deleteError)
	}

	return nil
}

func (c *RadosBlockDeviceClient) executeRBDDelete(pool, name string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("executeRBDDelete")

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rbd", "--pool", pool, "rm", name)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ERROR: rbd rm failed: %w", err)
	}

	return nil
}
