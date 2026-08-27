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
	"mgw-module-manager-migration/old_impl/libs/module_lib"

	"gopkg.in/yaml.v3"
)

type Base struct {
	Version string `yaml:"modfileVersion" json:"modfileVersion"`
}

type MfWrapper struct {
	version    string
	modFile    any
	decoders   Decoders
	generators Generators
}

type Decoders map[string]func(*yaml.Node) (any, error)

type Generators map[string]func(any) (*module_lib.Module, error)
