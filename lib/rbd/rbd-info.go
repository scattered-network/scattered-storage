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

// GetImageInfo Gathers *RBD image info for the '<pool>/<name>' image.
func (c *RadosBlockDeviceClient) GetImageInfo(pool, name string) (*RBD, error) {
	if !ValidatePool(pool) {
		return nil, validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return nil, validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("GetImageInfo")

	client := &RadosBlockDeviceClient{}
	image, infoError := client.executeRBDInfo(pool, name)
	if infoError != nil {
		return nil, infoError
	}

	return image, nil
}

// executeRBDInfo executes rbd info --format json for the given RBD image.
func (c *RadosBlockDeviceClient) executeRBDInfo(pool, name string) (*RBD, error) {
	log.Trace().Str("Pool", pool).Str("Name", name).Msg("executeRBDInfo")

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rbd", "--pool", pool, "info", name, "--format", "json")
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var image *RBD

	if err := json.Unmarshal(stdOut.Bytes(), &image); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return image, nil
}
