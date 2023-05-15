package ceph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/language"
)

// ApplicationTag
/* ceph --format json osd pool application get <pool_name>

{"cephfs":{"data":"mushroomfs"}}
{"cephfs":{"metadata":"mushroomfs"}}
{"mgr_devicehealth":{}}
{"rbd":{}}
{"rgw":{}}

ApplicationTag is used to determine a pool's application tag(s). */
type ApplicationTag struct {
	RBD             *struct{} `json:"rbd"`
	MgrDevicehealth *struct{} `json:"mgr_devicehealth"` //nolint:tagliatelle
	RGW             *struct{} `json:"rgw"`
	Cephfs          *struct {
		Data     *string `json:"data"`
		Metadata *string `json:"metadata"`
	} `json:"cephfs"`
}

func (c *CephCLI) GetApplicationTag(pool string) (*ApplicationTag, error) {
	var stdOut, stdErr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ceph", "osd", "pool", "application", "get", pool, "--format", "json")
	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutingCommand)

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	result := &ApplicationTag{}

	if err := json.Unmarshal(stdOut.Bytes(), &result); err != nil {
		log.Trace().Str("Response", stdOut.String()).Str("Error", err.Error()).Msg(language.ErrUnmarshalling)

		return nil, fmt.Errorf("%w", err)
	}

	log.Trace().Str("Command", cmd.String()).Msg(language.InfoExecutionCompleted)

	return result, nil
}

func (c *CephCLI) IsRBDPool(pool string) (bool, error) {
	if tag, err := c.GetApplicationTag(pool); err != nil {
		return false, err
	} else {
		if tag.RBD != nil {
			return true, nil
		}
	}
	return false, nil
}

func (c *CephCLI) IsRGWPool(pool string) (bool, error) {
	if tag, err := c.GetApplicationTag(pool); err != nil {
		return false, err
	} else {
		if tag.RGW != nil {
			return true, nil
		}
	}
	return false, nil
}

func (c *CephCLI) IsMgrDevicehealthPool(pool string) (bool, error) {
	if tag, err := c.GetApplicationTag(pool); err != nil {
		return false, err
	} else {
		if tag.MgrDevicehealth != nil {
			return true, nil
		}
	}
	return false, nil
}
