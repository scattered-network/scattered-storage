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

func (c *RadosBlockDeviceClient) executeRBDMap(pool, name string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rbd", "--exclusive", "--options", "lock_timeout=10", "--pool", pool, "map", name)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *RadosBlockDeviceClient) executeAddLock(pool, name string) error {
	log.Trace().Msg("starting executeAddLock")

	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	imageReference := pool + "/" + name

	cmd := exec.CommandContext(ctx, "rbd", "lock", "add", imageReference, "scattered-storage-lock")
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *RadosBlockDeviceClient) executeListLocks(pool, name string) ([]*Lock, error) {
	log.Trace().Msg("executeListLocks")

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rbd", "--format", "json", "-p", pool, "lock", "ls", name)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	var list []*Lock

	if err := cmd.Run(); err != nil {
		return list, fmt.Errorf("ERROR: rbd lock ls failed:\n%w", err)
	}

	if err := json.Unmarshal(stdOut.Bytes(), &list); err != nil {
		return list, fmt.Errorf(
			"ERROR: json for rbd lock ls could not unmarshal: %w\n%s", err, stdOut.String(),
		)
	}

	return list, nil
}

func (c *RadosBlockDeviceClient) executeRemoveLock(pool, name string, lock *Lock) error {
	log.Trace().Str("Pool", pool).Str("Name", name).Interface("Lock", lock).Msg("executeRemoveLock")

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	imageReference := pool + "/" + name

	cmd := exec.CommandContext(ctx, "rbd", "lock", "remove", imageReference, "'"+lock.ID+"'", lock.Locker)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ERROR: rbd lock failed: %w", err)
	}

	return nil
}

func (c *RadosBlockDeviceClient) executeWipeFSWithoutAction(device string) (*WipeFS, error) {
	log.Trace().Str("Device", device).Msg("executeWipeFSWithoutAction")

	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "wipefs", "-J", "-n", device)
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return &WipeFS{
			Signatures: nil,
		}, fmt.Errorf("ERROR: wipefs failed: %w", err)
	}

	var signatures *WipeFS

	if err := json.Unmarshal(stdOut.Bytes(), &signatures); err != nil {
		return &WipeFS{
				Signatures: nil,
			}, fmt.Errorf(
				"ERROR: json for wipefs could not unmarshal:\n%w\n%s", err, stdOut.String(),
			)
	}

	return signatures, nil
}
