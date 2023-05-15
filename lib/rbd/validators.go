package rbd

import (
	"github.com/rs/zerolog/log"
	"github.com/scattered-network/scattered-storage/lib/validators"
	"github.com/spf13/cast"
)

func ValidateName(name string) bool {
	nameExpression := "^[a-zA-Z0-9-_.]+$"
	if nameCheck := validators.ValidateRegex(nameExpression); nameCheck != nil {
		return validators.ValidateInput(nameCheck, name)
	}

	return false
}

func ValidatePool(pool string) bool {
	poolExpression := "^[a-zA-Z0-9-_]+$"
	if poolCheck := validators.ValidateRegex(poolExpression); poolCheck != nil {
		return validators.ValidateInput(poolCheck, pool)
	}

	return false
}

func ValidateSize(size int) bool {
	if size == 0 { // size MUST be greater than 0
		return false
	}

	sizeExpression := "^[0-9]+"
	if sizeCheck := validators.ValidateRegex(sizeExpression); sizeCheck != nil {
		return validators.ValidateInput(sizeCheck, cast.ToString(size))
	}

	return false
}

func ValidateSuffix(suffix string) bool {
	suffixExpression := "^[bBkKmMgGtTpP]$"
	if suffixCheck := validators.ValidateRegex(suffixExpression); suffixCheck != nil {
		return validators.ValidateInput(suffixCheck, suffix)
	}

	return false
}

func ValidateDevicePath(device string) bool {
	deviceExpression := "^/dev/[[:alnum:]]+$"
	if deviceCheck := validators.ValidateRegex(deviceExpression); deviceCheck != nil {
		return validators.ValidateInput(deviceCheck, device)
	}

	return false
}

func ValidateMakeFilesystemOptions(fsOptions *MkfsOptions) bool {
	if fsType, ok := fsOptions.Options["fsType"]; ok {
		stringType := cast.ToString(fsType.Value)
		switch stringType {
		case "xfs":
			break
		case "ext4":
			break
		default:
			log.Error().Str("Filesystem", stringType).Msg("Filesystem is not supported")

			return false
		}
	}

	if noDiscard, ok := fsOptions.Options["noDiscard"]; !ok {
		log.Error().Interface("Value", noDiscard.Value).Msg("invalid noDiscard option value")

		return false
	}

	return true
}
