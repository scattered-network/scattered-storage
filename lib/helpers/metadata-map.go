package helpers

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

type Metadata struct {
	Value       interface{}
	DataType    string
	Description string
}

func (m *Metadata) SetValue(value interface{}, dataType string) {
	switch dataType {
	case TypeBoolean:
		m.Value = cast.ToBool(value)
	case TypeString:
		m.Value = cast.ToString(value)
	case TypeInt:
		m.Value = cast.ToInt(value)
	case TypeUint:
		m.Value = cast.ToUint(value)
	case TypeUint32:
		m.Value = cast.ToUint32(value)
	case TypeUint64:
		m.Value = cast.ToUint64(value)
	case TypeInt64:
		m.Value = cast.ToInt64(value)
	case TypeFloat64:
		m.Value = cast.ToFloat64(value)
	case TypeTime:
		m.Value = cast.ToTime(value)
	case TypeDuration:
		m.Value = cast.ToDuration(value)
	case TypeStringSlice:
		m.Value = cast.ToStringSlice(value)
	case TypeIntSlice:
		m.Value = cast.ToIntSlice(value)
	default:
		log.Error().Str("Type", dataType).Interface("Value", value).Msg("unknown type")

		return
	}

	m.DataType = dataType
}

func (m *Metadata) GetValue() interface{} { //nolint:cyclop
	switch m.DataType {
	case "bool":
		return cast.ToBool(m.Value)
	case "string":
		return cast.ToString(m.Value)
	case "int":
		return cast.ToInt(m.Value)
	case "uint":
		return cast.ToUint(m.Value)
	case "uint32":
		return cast.ToUint32(m.Value)
	case "uint64":
		return cast.ToUint64(m.Value)
	case "int64":
		return cast.ToInt64(m.Value)
	case "float64":
		return cast.ToFloat64(m.Value)
	case "time":
		return cast.ToTime(m.Value)
	case "duration":
		return cast.ToDuration(m.Value)
	case "stringSlice":
		return cast.ToStringSlice(m.Value)
	case "intSlice":
		return cast.ToIntSlice(m.Value)
	}

	return nil
}
