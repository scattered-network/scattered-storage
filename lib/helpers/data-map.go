package helpers

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

const (
	TypeString      = "string"
	TypeBoolean     = "bool"
	TypeInt         = "int"
	TypeUint        = "uint"
	TypeUint32      = "uint32"
	TypeUint64      = "uint64"
	TypeInt64       = "int64"
	TypeFloat64     = "float64"
	TypeTime        = "time"
	TypeDuration    = "duration"
	TypeStringSlice = "stringSlice"
	TypeIntSlice    = "intSlice"
)

type Data struct {
	Value    interface{}
	DataType string
	Metadata map[string]*Metadata
}

func (d *Data) SetValue(value interface{}, dataType string) {
	switch dataType {
	case TypeBoolean:
		d.Value = cast.ToBool(value)
	case TypeString:
		d.Value = cast.ToString(value)
	case TypeInt:
		d.Value = cast.ToInt(value)
	case TypeUint:
		d.Value = cast.ToUint(value)
	case TypeUint32:
		d.Value = cast.ToUint32(value)
	case TypeUint64:
		d.Value = cast.ToUint64(value)
	case TypeInt64:
		d.Value = cast.ToInt64(value)
	case TypeFloat64:
		d.Value = cast.ToFloat64(value)
	case TypeTime:
		d.Value = cast.ToTime(value)
	case TypeDuration:
		d.Value = cast.ToDuration(value)
	case TypeStringSlice:
		d.Value = cast.ToStringSlice(value)
	case TypeIntSlice:
		d.Value = cast.ToIntSlice(value)
	default:
		log.Error().Str("Type", dataType).Interface("Value", value).Msg("unknown type")
	}

	d.DataType = dataType
}

func (d *Data) GetValue() interface{} {
	switch d.DataType {
	case TypeBoolean:
		return cast.ToBool(d.Value)
	case TypeString:
		return cast.ToString(d.Value)
	case TypeInt:
		return cast.ToInt(d.Value)
	case TypeUint:
		return cast.ToUint(d.Value)
	case TypeUint32:
		return cast.ToUint32(d.Value)
	case TypeUint64:
		return cast.ToUint64(d.Value)
	case TypeInt64:
		return cast.ToInt64(d.Value)
	case TypeFloat64:
		return cast.ToFloat64(d.Value)
	case TypeTime:
		return cast.ToTime(d.Value)
	case TypeDuration:
		return cast.ToDuration(d.Value)
	case TypeStringSlice:
		return cast.ToStringSlice(d.Value)
	case TypeIntSlice:
		return cast.ToIntSlice(d.Value)
	}

	return nil
}
