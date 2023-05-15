package rbd

import (
	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/validators"
)

func (c *RBDClient) listLocks(pool, name string) ([]*Lock, error) {
	var list []*Lock

	if !ValidatePool(pool) {
		return list, validators.ErrInvalidPoolName //nolint:exhaustruct
	}

	if !ValidateName(name) {
		return list, validators.ErrInvalidRBDName //nolint:exhaustruct
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("ListLocks")
	client := &RBDClient{}
	return client.executeListLocks(pool, name)
}

func (c *RBDClient) addLock(pool, name string) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Msg("AddLock")
	client := &RBDClient{}
	return client.executeAddLock(pool, name)
}

func (c *RBDClient) removeLock(pool, name string, lock *Lock) error {
	if !ValidatePool(pool) {
		return validators.ErrInvalidPoolName
	}

	if !ValidateName(name) {
		return validators.ErrInvalidRBDName
	}

	log.Trace().Str("Pool", pool).Str("Name", name).Interface("Lock", lock).Msg("RemoveLock")

	client := &RBDClient{}
	return client.executeRemoveLock(pool, name, lock)
}
