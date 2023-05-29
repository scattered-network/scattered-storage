package rbd

import (
	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/validators"
)

func (c *RadosBlockDeviceClient) listLocks(pool, name string) ([]*Lock, error) {
	var list []*Lock

	if !ValidatePool(pool) {
		return list, validators.ErrInvalidPoolName //nolint:exhaustruct
	}

	if !ValidateName(name) {
		return list, validators.ErrInvalidRBDName //nolint:exhaustruct
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("ListLocks")
	return c.executeListLocks(pool, name)
}

func (c *RadosBlockDeviceClient) addLock(pool, name string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("AddLock")
	return c.executeAddLock(pool, name)
}

func (c *RadosBlockDeviceClient) removeLock(pool, name string, lock *Lock) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Interface("Lock", lock).Msg("RemoveLock")

	return c.executeRemoveLock(pool, name, lock)
}
