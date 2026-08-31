/*
 * Copyright 2023 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package configs

import (
	"fmt"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/v1/model"
	module_lib2 "mgw-module-manager-migration/pkg/components/old_impl/libs/module_lib"
	"strconv"
)

func SetSlice(ref string, mfCV model.ConfigValue, mCs module_lib2.Configs) error {
	dataType := module_lib2.StringType
	if mfCV.DataType != nil {
		dataType = *mfCV.DataType
	}
	var configType string
	var cTypeOption map[string]any
	if mfCV.UserInput != nil {
		configType = mfCV.UserInput.Type
		cTypeOption = mfCV.UserInput.TypeOptions
	}

	delimiter := ","
	if mfCV.Delimiter != nil {
		delimiter = *mfCV.Delimiter
	}
	switch dataType {
	case module_lib2.StringType:
		d, o, co, err := parseConfigSlice(mfCV.Value, mfCV.Options, cTypeOption, parseConfigValueString)
		if err != nil {
			return fmt.Errorf("error parsing config '%s': %s", ref, err)
		}
		mCs.SetStringSlice(ref, d, o, mfCV.OptionsExt, configType, co, delimiter, !mfCV.Optional)
	case module_lib2.BoolType:
		d, o, co, err := parseConfigSlice(mfCV.Value, mfCV.Options, cTypeOption, parseConfigValueBool)
		if err != nil {
			return fmt.Errorf("error parsing config '%s': %s", ref, err)
		}
		mCs.SetBoolSlice(ref, d, o, mfCV.OptionsExt, configType, co, delimiter, !mfCV.Optional)
	case module_lib2.Int64Type:
		d, o, co, err := parseConfigSlice(mfCV.Value, mfCV.Options, cTypeOption, parseConfigValueInt64)
		if err != nil {
			return fmt.Errorf("error parsing config '%s': %s", ref, err)
		}
		mCs.SetInt64Slice(ref, d, o, mfCV.OptionsExt, configType, co, delimiter, !mfCV.Optional)
	case module_lib2.Float64Type:
		d, o, co, err := parseConfigSlice(mfCV.Value, mfCV.Options, cTypeOption, parseConfigValueFloat64)
		if err != nil {
			return fmt.Errorf("error parsing config '%s': %s", ref, err)
		}
		mCs.SetFloat64Slice(ref, d, o, mfCV.OptionsExt, configType, co, delimiter, !mfCV.Optional)
	default:
		return fmt.Errorf("%s invalid data type '%s'", ref, dataType)
	}
	return nil
}

func SetValue(ref string, mfCV model.ConfigValue, mCs module_lib2.Configs) error {
	dataType := module_lib2.StringType
	if mfCV.DataType != nil {
		dataType = *mfCV.DataType
	}
	var configType string
	var cTypeOption map[string]any
	if mfCV.UserInput != nil {
		configType = mfCV.UserInput.Type
		cTypeOption = mfCV.UserInput.TypeOptions
	}
	switch dataType {
	case module_lib2.StringType:
		d, o, co, err := parseConfig(mfCV.Value, mfCV.Options, cTypeOption, parseConfigValueString)
		if err != nil {
			return fmt.Errorf("error parsing config '%s': %s", ref, err)
		}
		mCs.SetString(ref, d, o, mfCV.OptionsExt, configType, co, !mfCV.Optional)
	case module_lib2.BoolType:
		d, o, co, err := parseConfig(mfCV.Value, mfCV.Options, cTypeOption, parseConfigValueBool)
		if err != nil {
			return fmt.Errorf("error parsing config '%s': %s", ref, err)
		}
		mCs.SetBool(ref, d, o, mfCV.OptionsExt, configType, co, !mfCV.Optional)
	case module_lib2.Int64Type:
		d, o, co, err := parseConfig(mfCV.Value, mfCV.Options, cTypeOption, parseConfigValueInt64)
		if err != nil {
			return fmt.Errorf("error parsing config '%s': %s", ref, err)
		}
		mCs.SetInt64(ref, d, o, mfCV.OptionsExt, configType, co, !mfCV.Optional)
	case module_lib2.Float64Type:
		d, o, co, err := parseConfig(mfCV.Value, mfCV.Options, cTypeOption, parseConfigValueFloat64)
		if err != nil {
			return fmt.Errorf("error parsing config '%s': %s", ref, err)
		}
		mCs.SetFloat64(ref, d, o, mfCV.OptionsExt, configType, co, !mfCV.Optional)
	default:
		return fmt.Errorf("%s invalid data type '%s'", ref, dataType)
	}
	return nil
}

func parseConfig[T any](val any, opt []any, ctOpt map[string]any, valParser func(any) (T, error)) (p *T, o []T, to module_lib2.ConfigTypeOptions, err error) {
	if val != nil {
		v, er := valParser(val)
		if er != nil {
			err = er
			return
		}
		p = &v
	}
	if o, err = parseConfigOptions(opt, valParser); err != nil {
		return
	}
	to, err = parseConfigTypeOptions(ctOpt)
	return
}

func parseConfigSlice[T any](val any, opt []any, ctOpt map[string]any, valParser func(any) (T, error)) (sl []T, o []T, to module_lib2.ConfigTypeOptions, err error) {
	if val != nil {
		v, ok := val.([]any)
		if !ok {
			err = fmt.Errorf("type missmatch: %T != slice", val)
			return
		}
		for _, i := range v {
			pi, e := valParser(i)
			if e != nil {
				err = e
				return
			}
			sl = append(sl, pi)
		}
	}
	if o, err = parseConfigOptions(opt, valParser); err != nil {
		return
	}
	to, err = parseConfigTypeOptions(ctOpt)
	return
}

func parseConfigOptions[T any](opt []any, valParser func(any) (T, error)) ([]T, error) {
	var opts []T
	for _, o := range opt {
		v, err := valParser(o)
		if err != nil {
			return nil, err
		}
		opts = append(opts, v)
	}
	return opts, nil
}

func parseConfigValueString(val any) (string, error) {
	var sVal string
	switch v := val.(type) {
	case string:
		sVal = v
	case float32:
		sVal = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		sVal = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		sVal = strconv.FormatInt(int64(v), 10)
	case int8:
		sVal = strconv.FormatInt(int64(v), 10)
	case int16:
		sVal = strconv.FormatInt(int64(v), 10)
	case int32:
		sVal = strconv.FormatInt(int64(v), 10)
	case int64:
		sVal = strconv.FormatInt(v, 10)
	default:
		return "", fmt.Errorf("invalid data type '%T'", val)
	}
	return sVal, nil
}

func parseConfigValueBool(val any) (bool, error) {
	v, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("invalid data type '%T'", val)
	}
	return v, nil
}

func parseConfigValueInt64(val any) (int64, error) {
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
	default:
		return i, fmt.Errorf("invalid data type '%T'", val)
	}
	return i, nil
}

func parseConfigValueFloat64(val any) (float64, error) {
	var f float64
	switch v := val.(type) {
	case float32:
		f = float64(v)
	case float64:
		f = v
	default:
		return f, fmt.Errorf("invalid data type '%T'", val)
	}
	return f, nil
}

func parseConfigTypeOptions(opt map[string]any) (module_lib2.ConfigTypeOptions, error) {
	o := make(module_lib2.ConfigTypeOptions)
	for key, val := range opt {
		switch v := val.(type) {
		case string:
			o.SetString(key, v)
		case bool:
			o.SetBool(key, v)
		case int:
			o.SetInt64(key, int64(v))
		case int8:
			o.SetInt64(key, int64(v))
		case int16:
			o.SetInt64(key, int64(v))
		case int32:
			o.SetInt64(key, int64(v))
		case int64:
			o.SetInt64(key, v)
		case float32:
			o.SetFloat64(key, float64(v))
		case float64:
			o.SetFloat64(key, v)
		default:
			return nil, fmt.Errorf("unknown data type '%T'", val)
		}
	}
	return o, nil
}
