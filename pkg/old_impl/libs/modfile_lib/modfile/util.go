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

package modfile

import (
	"errors"
	"fmt"
	"mgw-module-manager-migration/pkg/old_impl/libs/module_lib"

	"gopkg.in/yaml.v3"
)

func New(decoders Decoders, generators Generators) *MfWrapper {
	return &MfWrapper{decoders: decoders, generators: generators}
}

func (w *MfWrapper) UnmarshalYAML(yn *yaml.Node) error {
	if len(w.decoders) < 1 {
		return errors.New("no decoders")
	}
	var mfb Base
	if err := yn.Decode(&mfb); err != nil {
		return err
	}
	if mfb.Version == "" {
		return errors.New("no version")
	}
	d, ok := w.decoders[mfb.Version]
	if !ok {
		return fmt.Errorf("no decoder for version '%s'", mfb.Version)
	}
	m, err := d(yn)
	if err != nil {
		return err
	}
	w.version = mfb.Version
	w.modFile = m
	return nil
}

func (w *MfWrapper) GetModule() (*module_lib.Module, error) {
	if len(w.generators) < 1 {
		return nil, errors.New("no generators")
	}
	g, ok := w.generators[w.version]
	if !ok {
		return nil, fmt.Errorf("no generator for version '%s'", w.version)
	}
	return g(w.modFile)
}

func (d Decoders) Add(f func() (string, func(*yaml.Node) (any, error))) {
	key, df := f()
	d[key] = df
}

func (g Generators) Add(f func() (string, func(any) (*module_lib.Module, error))) {
	key, gf := f()
	g[key] = gf
}
