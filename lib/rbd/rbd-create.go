package rbd

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/language"
	"github.com/scattered-network/scattered-storage/lib/validators"
	"github.com/spf13/cast"
)

// CreateRBD validates the creation options and triggers the rbd create command.
func CreateRBD(pool, name string, size int, suffix string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	if !ValidateSize(size) {
		return validators.ErrInvalidSize
	}

	if !ValidateSuffix(suffix) {
		return validators.ErrInvalidSuffix
	}

	client := &RBDClient{}
	if createError := client.executeRBDCreate(pool, name, size, suffix); createError != nil {
		return createError
	}

	return nil
}

// executeRBDCreate runs the rbd create command enabling the following features:
// layering, striping, exclusive-lock, object-map, and fast-diff.
func (c *RBDClient) executeRBDCreate(pool, name string, size int, suffix string) error {
	log.Trace().Msg("executeRBDCreate")

	sizeArgument := cast.ToString(size) + suffix
	toCreate := pool + "/" + name

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		ctx, "rbd", "create", "--image-feature", "layering", "--image-feature", "striping", "--image-feature",
		"exclusive-lock", "--image-feature", "object-map", "--image-feature", "fast-diff", "--size", sizeArgument,
		toCreate,
	)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.Is(err, exitError); ok {
			if exitError.ExitCode() == 17 {
				return fmt.Errorf("%w", validators.ErrRBDExists)
			}

			return fmt.Errorf("ERROR: rbd create failed: %w", err)
		}
	}

	return nil
}
