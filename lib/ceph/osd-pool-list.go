package ceph

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
)

// OSDPoolList
/* ceph osd pool ls --format json

[
  "rbd",
  "device_health_metrics",
  "rbd-ssd",
  ".rgw.root",
  "default.rgw.control",
  "default.rgw.meta",
  "default.rgw.log",
  "default.rgw.buckets.index",
  "default.rgw.buckets.data",
  "default.rgw.buckets.non-ec",
  "docker-ssd",
  "docker-hdd",
  "cephfs.mushroomfs.meta",
  "cephfs.mushroomfs.data"
]

OSDPoolList is used process the pool list output. */
type OSDPoolList helpers.List

func (c *CephCLI) GetOSDPoolList() (*OSDPoolList, error) {
	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ceph", "osd", "pool", "ls", "--format", "json")
	log.Trace().Str("Command", cmd.String()).Msg("Executing command")

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		log.Error().Str("Error", err.Error()).Msg(language.ErrExecutingCommand)

		return nil, fmt.Errorf("%w", err)
	}

	result := &OSDPoolList{}

	if err := json.Unmarshal(stdOut.Bytes(), &result); err != nil {
		log.Error().Str("Response", stdOut.String()).Str("Error", err.Error()).
			Msg("Encountered Error Unmarshalling Response")

		return nil, fmt.Errorf("%w", err)
	}

	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutionCompleted)

	return result, nil
}

func (c *CephCLI) GetRBDPools() OSDPoolList {
	log.Trace().Msg("Getting list of RBD pools")

	found := map[int]string{}
	count := 0

	pools, err := c.GetOSDPoolList()
	if err != nil {
		log.Error().Str("Error", err.Error()).Msg(language.ErrGettingOSDPoolList)

		return nil
	}

	for _, pool := range *pools {
		valid, err := c.IsRBDPool(pool)
		if err != nil {
			log.Error().Str("Error", err.Error()).Msg(language.ErrCheckingApplicationTag)

			continue
		}

		if valid {
			found[count] = pool
			count++
		}
	}

	foundPools := make(OSDPoolList, len(found))
	for index, pool := range found {
		foundPools[index] = pool
	}

	log.Trace().Int("Count", len(foundPools)).Msg("Getting list of RBD pools completed")

	return foundPools
}
