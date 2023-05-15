package helpers

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/language"
)

type Executable struct {
	command  string
	args     []string
	timeout  time.Duration
	stdout   *bytes.Buffer
	stderr   *bytes.Buffer
	exitCode int
}

func NewExecutable(command string, args []string, timeout time.Duration) *Executable {
	return &Executable{command: command, args: args, timeout: timeout}
}

func (e *Executable) Execute() error {
	e.stdout = &bytes.Buffer{}
	e.stderr = &bytes.Buffer{}
	e.exitCode = 0

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.command, e.args...)
	log.Trace().Str("Command", cmd.String()).Interface("Executable", e).Msg(language.InfoExecutingCommand)

	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w", err)
	}

	log.Trace().Str("Command", cmd.String()).Str("Output", e.stdout.String()).Msg(language.InfoExecutionCompleted)

	return nil
}
