/*
 * Copyright 2022 InfAI (CC SES)
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

package module_lib

import (
	"encoding/json"
	"strings"
)

func newConfigValue[T any](def *T, opt []T, dType DataType, optExt bool, cType string, cTypeOpt ConfigTypeOptions, required bool) configValue {
	cv := configValue{
		OptExt:   optExt,
		Type:     cType,
		DataType: dType,
		Required: required,
	}
	if def != nil {
		cv.Default = *def
	}
	if opt != nil && len(opt) > 0 {
		cv.Options = opt
	}
	if cTypeOpt != nil && len(cTypeOpt) > 0 {
		cv.TypeOpt = cTypeOpt
	}
	return cv
}

func newConfigValueSlice[T any](def []T, opt []T, dType DataType, optExt bool, cType string, cTypeOpt ConfigTypeOptions, delimiter string, required bool) configValue {
	cv := configValue{
		OptExt:    optExt,
		Type:      cType,
		DataType:  dType,
		IsSlice:   true,
		Delimiter: delimiter,
		Required:  required,
	}
	if def != nil && len(def) > 0 {
		cv.Default = def
	}
	if opt != nil && len(opt) > 0 {
		cv.Options = opt
	}
	if cTypeOpt != nil && len(cTypeOpt) > 0 {
		cv.TypeOpt = cTypeOpt
	}
	return cv
}

func (c Configs) SetString(ref string, def *string, opt []string, optExt bool, cType string, cTypeOpt ConfigTypeOptions, required bool) {
	c[ref] = newConfigValue(def, opt, StringType, optExt, cType, cTypeOpt, required)
}

func (c Configs) SetBool(ref string, def *bool, opt []bool, optExt bool, cType string, cTypeOpt ConfigTypeOptions, required bool) {
	c[ref] = newConfigValue(def, opt, BoolType, optExt, cType, cTypeOpt, required)
}

func (c Configs) SetInt64(ref string, def *int64, opt []int64, optExt bool, cType string, cTypeOpt ConfigTypeOptions, required bool) {
	c[ref] = newConfigValue(def, opt, Int64Type, optExt, cType, cTypeOpt, required)
}

func (c Configs) SetFloat64(ref string, def *float64, opt []float64, optExt bool, cType string, cTypeOpt ConfigTypeOptions, required bool) {
	c[ref] = newConfigValue(def, opt, Float64Type, optExt, cType, cTypeOpt, required)
}

func (c Configs) SetStringSlice(ref string, def []string, opt []string, optExt bool, cType string, cTypeOpt ConfigTypeOptions, delimiter string, required bool) {
	c[ref] = newConfigValueSlice(def, opt, StringType, optExt, cType, cTypeOpt, delimiter, required)
}

func (c Configs) SetBoolSlice(ref string, def []bool, opt []bool, optExt bool, cType string, cTypeOpt ConfigTypeOptions, delimiter string, required bool) {
	c[ref] = newConfigValueSlice(def, opt, BoolType, optExt, cType, cTypeOpt, delimiter, required)
}

func (c Configs) SetInt64Slice(ref string, def []int64, opt []int64, optExt bool, cType string, cTypeOpt ConfigTypeOptions, delimiter string, required bool) {
	c[ref] = newConfigValueSlice(def, opt, Int64Type, optExt, cType, cTypeOpt, delimiter, required)
}

func (c Configs) SetFloat64Slice(ref string, def []float64, opt []float64, optExt bool, cType string, cTypeOpt ConfigTypeOptions, delimiter string, required bool) {
	c[ref] = newConfigValueSlice(def, opt, Float64Type, optExt, cType, cTypeOpt, delimiter, required)
}

func (o ConfigTypeOptions) SetString(ref string, val string) {
	o[ref] = configTypeOption{
		Value:    val,
		DataType: StringType,
	}
}

func (o ConfigTypeOptions) SetBool(ref string, val bool) {
	o[ref] = configTypeOption{
		Value:    val,
		DataType: BoolType,
	}
}

func (o ConfigTypeOptions) SetInt64(ref string, val int64) {
	o[ref] = configTypeOption{
		Value:    val,
		DataType: Int64Type,
	}
}

func (o ConfigTypeOptions) SetFloat64(ref string, val float64) {
	o[ref] = configTypeOption{
		Value:    val,
		DataType: Float64Type,
	}
}

func (v configValue) OptionsLen() (l int) {
	switch o := v.Options.(type) {
	case []string:
		l = len(o)
	case []bool:
		l = len(o)
	case []int64:
		l = len(o)
	case []float64:
		l = len(o)
	}
	return
}

func (t SrvRefTarget) FillTemplate(s string) string {
	if t.Template != nil {
		return strings.ReplaceAll(*t.Template, "{"+RefPlaceholder+"}", s)
	}
	return s
}

func (t ExtDependencyTarget) FillTemplate(s string) string {
	if t.Template != nil {
		return strings.ReplaceAll(*t.Template, "{"+RefPlaceholder+"}", s)
	}
	return s
}

type Set[T comparable] map[T]struct{}

func (s *Set[T]) UnmarshalJSON(b []byte) error {
	var sl []T
	if err := json.Unmarshal(b, &sl); err != nil {
		return err
	}
	set := make(Set[T])
	for _, item := range sl {
		set[item] = struct{}{}
	}
	*s = set
	return nil
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	var sl []T
	for item := range s {
		sl = append(sl, item)
	}
	return json.Marshal(sl)
}
