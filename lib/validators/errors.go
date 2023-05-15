package validators

import "errors"

var (
	ErrRBDExists                   = errors.New("rbd already exists")
	ErrInvalidRegex                = errors.New("invalid regex")
	ErrInvalidRBDName              = errors.New("invalid rbd name")
	ErrInvalidPoolName             = errors.New("invalid pool name")
	ErrInvalidSize                 = errors.New("invalid rbd size")
	ErrInvalidSuffix               = errors.New("invalid rbd size suffix")
	ErrInvalidDevicePath           = errors.New("invalid device path")
	ErrInvalidMakeOptions          = errors.New("invalid make options")
	ErrNotTaggedForRBD             = errors.New("pool does not have the 'rbd' application tag")
	ErrNotTaggedForRGW             = errors.New("pool does not have the 'rgw' application tag")
	ErrNotTaggedForMgrDevicehealth = errors.New("pool does not have the 'mgr_devicehealth' application tag")
)
