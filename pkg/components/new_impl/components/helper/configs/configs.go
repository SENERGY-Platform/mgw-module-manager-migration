package configs

import (
	"errors"
	"fmt"
	"math"
	"mgw-module-manager-migration/pkg/components/new_impl/models"
	"mgw-module-manager-migration/pkg/components/new_impl/models/constants"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/module_lib"
	"slices"
)

func ValueIsEqual(a, b models.Value) bool {
	if a.DataType != b.DataType {
		return false
	}
	if a.IsSlice != b.IsSlice {
		return false
	}
	switch a.DataType {
	case constants.ValueDataTypeString:
		if a.IsSlice {
			return slices.Equal(a.StringSlice, b.StringSlice)
		}
		return a.String == b.String
	case constants.ValueDataTypeInt64:
		if a.IsSlice {
			return slices.Equal(a.Int64Slice, b.Int64Slice)
		}
		return a.Int64 == b.Int64
	case constants.ValueDataTypeFloat64:
		if a.IsSlice {
			return slices.Equal(a.Float64Slice, b.Float64Slice)
		}
		return a.Float64 == b.Float64
	case constants.ValueDataTypeBool:
		if a.IsSlice {
			return slices.Equal(a.BoolSlice, b.BoolSlice)
		}
		return a.Bool == b.Bool
	}
	return false
}

func GetValue(val any, dataType int, isSlice bool) (models.Value, error) {
	config := models.Value{
		DataType: dataType,
		IsSlice:  isSlice,
	}
	if isSlice {
		switch dataType {
		case constants.ValueDataTypeString:
			v, err := toTypedSlice(val, toString)
			if err != nil {
				return models.Value{}, err
			}
			config.StringSlice = v
		case constants.ValueDataTypeBool:
			v, err := toTypedSlice(val, toBool)
			if err != nil {
				return models.Value{}, err
			}
			config.BoolSlice = v
		case constants.ValueDataTypeInt64:
			v, err := toTypedSlice(val, toInt64)
			if err != nil {
				return models.Value{}, err
			}
			config.Int64Slice = v
		case constants.ValueDataTypeFloat64:
			v, err := toTypedSlice(val, toFloat64)
			if err != nil {
				return models.Value{}, err
			}
			config.Float64Slice = v
		default:
			return models.Value{}, errors.New(fmt.Sprintf("unsupported data type: '%d'", dataType))
		}
	} else {
		switch dataType {
		case constants.ValueDataTypeString:
			v, err := toString(val)
			if err != nil {
				return models.Value{}, err
			}
			config.String = v
		case constants.ValueDataTypeBool:
			v, err := toBool(val)
			if err != nil {
				return models.Value{}, err
			}
			config.Bool = v
		case constants.ValueDataTypeInt64:
			v, err := toInt64(val)
			if err != nil {
				return models.Value{}, err
			}
			config.Int64 = v
		case constants.ValueDataTypeFloat64:
			v, err := toFloat64(val)
			if err != nil {
				return models.Value{}, err
			}
			config.Float64 = v
		default:
			return models.Value{}, errors.New(fmt.Sprintf("unsupported data type: '%d'", dataType))
		}
	}
	return config, nil
}

func GetDataType(moduleDataType string) int {
	return moduleDataTypeMap[moduleDataType]
}

func toTypedSlice[T any](val any, converter func(any) (T, error)) ([]T, error) {
	tSlice, ok := val.([]T)
	if !ok {
		anySlice, ok := val.([]any)
		if !ok {
			return nil, errors.New(fmt.Sprintf("invalid data type: '%T'", val))
		}
		for _, item := range anySlice {
			v, err := converter(item)
			if err != nil {
				return nil, err
			}
			tSlice = append(tSlice, v)
		}
	}
	return tSlice, nil
}

func toString(val any) (string, error) {
	v, ok := val.(string)
	if !ok {
		return "", errors.New("invalid data type: 'string' required")
	}
	return v, nil
}

func toBool(val any) (bool, error) {
	v, ok := val.(bool)
	if !ok {
		return false, errors.New("invalid data type: 'boolean' required")
	}
	return v, nil
}

func float64ToInt64(val float64) (int64, bool) {
	i, fr := math.Modf(val)
	if fr > 0 {
		return 0, false
	}
	return int64(i), true
}

func toInt64(val any) (int64, error) {
	var i int64
	switch v := val.(type) {
	case int:
		i = int64(v)
	case int8:
		i = int64(v)
	case int16:
		i = int64(v)
	case int32:
		i = int64(v)
	case int64:
		i = v
	case float32:
		var ok bool
		i, ok = float64ToInt64(float64(v))
		if !ok {
			return 0, errors.New("invalid data type: 'integer' required")
		}
	case float64:
		var ok bool
		i, ok = float64ToInt64(v)
		if !ok {
			return 0, errors.New("invalid data type: 'integer' required")
		}
	default:
		return 0, errors.New("invalid data type: 'integer' required")
	}
	return i, nil
}

func toFloat64(val any) (float64, error) {
	var f float64
	switch v := val.(type) {
	case float32:
		f = float64(v)
	case float64:
		f = v
	default:
		return 0, errors.New("invalid data type: 'float' required")
	}
	return f, nil
}

var moduleDataTypeMap = map[string]int{
	module_lib.StringType:  constants.ValueDataTypeString,
	module_lib.Int64Type:   constants.ValueDataTypeInt64,
	module_lib.Float64Type: constants.ValueDataTypeFloat64,
	module_lib.BoolType:    constants.ValueDataTypeBool,
}
