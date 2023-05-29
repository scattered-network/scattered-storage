package rbd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
)

// ShowMapped
/* rbd --format json showmapped
[
  {
    "id": "0",
    "pool": "rbd",
    "namespace": "",
    "name": "test1",
    "snap": "-",
    "device": "/dev/rbd0"
  },
  {
    "id": "1",
    "pool": "rbd",
    "namespace": "",
    "name": "test2",
    "snap": "-",
    "device": "/dev/rbd1"
  }
]

ShowMapped is used to determine which images have been mapped to the local node. */
type ShowMapped []struct {
	ID        string `json:"id"`
	Pool      string `json:"pool"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Snap      string `json:"snap"`
	Device    string `json:"device"`
}

// ListMappedImages returns the RBD images mapped to the host.
func (c *RadosBlockDeviceClient) ListMappedImages() (*ShowMapped, error) {
	list, listError := c.executeShowMapped()
	if listError != nil {
		log.Error().Str("Error", listError.Error()).Msg("Could not list mapped images")

		return nil, listError
	}

	return list, nil
}

// executeShowMapped runs the rbd showmapped --format json command and returns the results as *ShowMapped.
func (c *RadosBlockDeviceClient) executeShowMapped() (*ShowMapped, error) {
	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rbd", "showmapped", "--format", "json")
	log.Trace().Str("Command", cmd.String()).Msg("Executing command")

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	var list *ShowMapped

	if err := cmd.Run(); err != nil {
		log.Error().Str("Error", err.Error()).Msg("Error executing command")

		return list, fmt.Errorf("%w", err)
	}

	if err := json.Unmarshal(stdOut.Bytes(), &list); err != nil {
		log.Error().Str("Response", stdOut.String()).Str("Error", err.Error()).
			Msg("Encountered Error Unmarshalling Response")

		return nil, fmt.Errorf("%w", err)
	}

	return list, nil
}
