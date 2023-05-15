package rbd

import (
	"github.com/rs/zerolog/log"
)

type RBDClient struct{}

// RBD
/* rbd --pool rbd info test-image --format json
{
  "name": "test-image",
  "id": "979ba5a95620ef",
  "size": 10737418240,
  "objects": 2560,
  "order": 22,
  "object_size": 4194304,
  "snapshot_count": 0,
  "block_name_prefix": "rbd_data.979ba5a95620ef",
  "format": 2,
  "features": [
    "layering"
  ],
  "op_features": [],
  "flags": [],
  "create_timestamp": "Sat May 21 15:31:59 2022",
  "access_timestamp": "Sat May 21 15:31:59 2022",
  "modify_timestamp": "Sat May 21 15:31:59 2022"
}
RBD is used to gather information and manipulate an RBD image. */
type RBD struct {
	Name            string         `json:"name"`
	ID              string         `json:"id"`
	Size            int64          `json:"size"`
	Objects         int            `json:"objects"`
	Order           int            `json:"order"`
	ObjectSize      int            `json:"object_size"`       //nolint:tagliatelle
	SnapshotCount   int            `json:"snapshot_count"`    //nolint:tagliatelle
	BlockNamePrefix string         `json:"block_name_prefix"` //nolint:tagliatelle
	Format          int            `json:"format"`
	Features        []*string      `json:"features"`
	OpFeatures      []*interface{} `json:"op_features"` //nolint:tagliatelle
	Flags           []*interface{} `json:"flags"`
	CreateTimestamp string         `json:"create_timestamp"` //nolint:tagliatelle
	AccessTimestamp string         `json:"access_timestamp"` //nolint:tagliatelle
	ModifyTimestamp string         `json:"modify_timestamp"` //nolint:tagliatelle
}

// isMapped requires the 'pool' and 'name' of the rbd image to check and returns
// both the device path and a bool representing the mapped state.
func (c *RBDClient) isMapped(pool, name string) (string, bool) {
	if !ValidatePool(pool) {
		return "", false
	}

	if !ValidateName(name) {
		return "", false
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("isMapped")

	list, listError := c.executeShowMapped()
	if listError != nil {
		log.Error().Str("Error", listError.Error()).Msg("error listing mapped images")

		return "", false
	}

	for _, image := range *list {
		if image.Name == name && image.Pool == pool {
			return image.Device, true
		}
	}

	return "", false
}
