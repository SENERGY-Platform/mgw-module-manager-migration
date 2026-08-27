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

package model

import (
	"code.cloudfoundry.org/bytefmt"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
	"time"
)

const Version = "v1"

func (p *Port) UnmarshalYAML(yn *yaml.Node) error {
	var it any
	_ = yn.Decode(&it)
	switch v := it.(type) {
	case int:
		if v < 0 {
			return fmt.Errorf("invalid port: %d", v)
		}
		*p = Port(strconv.FormatInt(int64(v), 10))
	case string:
		parts := strings.Split(v, "-")
		if len(parts) > 2 {
			return fmt.Errorf("invalid port range: %s", v)
		}
		for i := 0; i < len(parts); i++ {
			n, err := strconv.ParseInt(parts[i], 10, 64)
			if err != nil || n < 0 {
				return fmt.Errorf("invalid port: %s", v)
			}
		}
		*p = Port(v)
	default:
		return fmt.Errorf("invlid port: %v", v)
	}
	return nil
}

func (p Port) Parse() ([]uint, error) {
	var r []uint
	parts := strings.Split(string(p), "-")
	pl := len(parts)
	if pl < 1 || pl > 2 {
		return nil, errors.New("invalid format")
	} else {
		var tmp []uint
		for _, part := range parts {
			n, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return nil, err
			}
			tmp = append(tmp, uint(n))
		}
		if len(tmp) > 1 {
			if tmp[0] >= tmp[1] {
				return nil, errors.New("invalid range")
			}
			for i := tmp[0]; i <= tmp[1]; i++ {
				r = append(r, i)
			}
		} else {
			r = tmp
		}
	}
	return r, nil
}

func (fb *ByteFmt) UnmarshalYAML(yn *yaml.Node) error {
	var it any
	_ = yn.Decode(&it)
	switch v := it.(type) {
	case int:
		*fb = ByteFmt(v)
	case string:
		bytes, err := bytefmt.ToBytes(v)
		if err != nil {
			return fmt.Errorf("invalid size: %s", err)
		}
		*fb = ByteFmt(bytes)
	default:
		return fmt.Errorf("invalid size: %v", v)
	}
	return nil
}

func (d *Duration) UnmarshalYAML(yn *yaml.Node) error {
	var s string
	if err := yn.Decode(&s); err != nil {
		return err
	}
	if dur, err := time.ParseDuration(s); err != nil {
		return err
	} else {
		*d = Duration(dur)
	}
	return nil
}

func (m *FileMode) UnmarshalYAML(yn *yaml.Node) error {
	var s string
	if err := yn.Decode(&s); err != nil {
		return err
	}
	i, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return err
	}
	*m = FileMode(i)
	return nil
}

func (r Resource) GetUserInput() *UserInput {
	return r.UserInput
}

func (c ConfigValue) GetUserInput() *UserInput {
	if c.UserInput != nil {
		return &c.UserInput.UserInput
	}
	return nil
}

func (t *StrOrSlice) UnmarshalYAML(yn *yaml.Node) error {
	var sl []string
	if err := yn.Decode(&sl); err != nil {
		var s string
		if err := yn.Decode(&s); err != nil {
			return err
		}
		sl = append(sl, s)
	}
	*t = sl
	return nil
}
