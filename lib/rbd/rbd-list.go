package rbd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/helpers"
	"github.com/scattered-network/scattered-storage/lib/language"
	"github.com/scattered-network/scattered-storage/lib/validators"
)

func (c *RadosBlockDeviceClient) GetRBDList(pool string) ([]string, error) {
	if !ValidatePool(pool) {
		return nil, validators.ErrInvalidPoolName
	}

	rbdList, listError := c.executeRBDList(pool)
	if listError != nil {
		return nil, fmt.Errorf("%w", listError)
	}

	return rbdList, nil
}

func (c *RadosBlockDeviceClient) executeRBDList(pool string) (helpers.List, error) {
	if !ValidatePool(pool) {
		return nil, validators.ErrInvalidPoolName
	}

	log.Trace().Str("Pool", pool).Msg("executeRBDList")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rbd", "--pool", pool, "list", "--format", "json")
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ERROR: rbd list failed: %w", err)
	}

	var RBDList helpers.List

	if err := json.Unmarshal(stdOut.Bytes(), &RBDList); err != nil {
		return nil, fmt.Errorf("ERROR: json for rbd list could not unmarshal:\n%w\n%s", err, stdOut.String())
	}

	return RBDList, nil
}
