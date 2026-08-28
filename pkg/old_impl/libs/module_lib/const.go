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

const (
	AddOnModule           ModuleType = "add-on"
	DeviceConnectorModule ModuleType = "device-connector"
)

var ModuleTypeMap = map[ModuleType]struct{}{
	AddOnModule:           {},
	DeviceConnectorModule: {},
}

const (
	SingleDeployment   DeploymentType = "single"
	MultipleDeployment DeploymentType = "multiple"
)

var DeploymentTypeMap = map[DeploymentType]struct{}{
	SingleDeployment:   {},
	MultipleDeployment: {},
}

const (
	TcpPort PortProtocol = "tcp"
	UdpPort PortProtocol = "udp"
)

var PortProtocolMap = map[PortProtocol]struct{}{
	TcpPort: {},
	UdpPort: {},
}

const (
	BoolType    DataType = "bool"
	Int64Type   DataType = "int"
	Float64Type DataType = "float"
	StringType  DataType = "string"
)

const (
	X86     CPUArch = "x86"
	I386    CPUArch = "i386"
	X86_64  CPUArch = "x86_64"
	AMD64   CPUArch = "amd64"
	AARCH32 CPUArch = "aarch32"
	ARM32V5 CPUArch = "arm32v5"
	ARM32V6 CPUArch = "arm32v6"
	ARM32V7 CPUArch = "arm32v7"
	AARCH64 CPUArch = "aarch64"
	ARM64V8 CPUArch = "arm64v8"
)

var CPUArchMap = map[CPUArch]struct{}{
	X86:     {},
	I386:    {},
	X86_64:  {},
	AMD64:   {},
	AARCH32: {},
	ARM32V5: {},
	ARM32V6: {},
	ARM32V7: {},
	AARCH64: {},
	ARM64V8: {},
}

const RefPlaceholder = "ref"
