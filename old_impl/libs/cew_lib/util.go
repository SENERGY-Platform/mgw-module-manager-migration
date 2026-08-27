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

package cew_lib

import (
	"encoding/json"
	"fmt"
	"net"
)

func (p *Port) KeyStr() string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

func (m *Mount) KeyStr() string {
	return fmt.Sprintf("%s:%s", m.Source, m.Target)
}

func (d *Device) KeyStr() string {
	return fmt.Sprintf("%s:%s", d.Source, d.Target)
}

func (s *Subnet) KeyStr() string {
	return fmt.Sprintf("%s/%d", net.IP(s.Prefix).String(), s.Bits)
}

func (i *IPAddr) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return
	}
	*i = IPAddr(net.ParseIP(s))
	return
}

func (i IPAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(net.IP(i))
}

func (e *cError) Error() string {
	return e.err.Error()
}

func (e *cError) Unwrap() error {
	return e.err
}

func NewInternalError(err error) error {
	return &InternalError{cError{err: err}}
}

func NewNotFoundError(err error) error {
	return &NotFoundError{cError{err: err}}
}

func NewInvalidInputError(err error) error {
	return &InvalidInputError{cError{err: err}}
}
