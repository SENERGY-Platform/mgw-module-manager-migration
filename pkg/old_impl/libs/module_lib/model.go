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
	"io/fs"
	"time"
)

type ModuleType = string

type DeploymentType = string

type CPUArch = string

type Module struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	Description    string                  `json:"description"`
	Tags           Set[string]             `json:"tags"`
	License        string                  `json:"license"`
	Author         string                  `json:"author"`
	Version        string                  `json:"version"`
	Type           ModuleType              `json:"type"`
	DeploymentType DeploymentType          `json:"deployment_type"`
	Architectures  Set[CPUArch]            `json:"architectures"`
	Services       map[string]*Service     `json:"services"`       // {ref:Service}
	Volumes        Set[string]             `json:"volumes"`        // {volName}
	Dependencies   map[string]string       `json:"dependencies"`   // {moduleID:moduleVersion}
	HostResources  map[string]HostResource `json:"host_resources"` // {ref:{tag}}
	Secrets        map[string]Secret       `json:"secrets"`        // {ref:Secret}
	Configs        Configs                 `json:"configs"`        // {ref:ConfigValue}
	Inputs         Inputs                  `json:"inputs"`
	AuxServices    map[string]*AuxService  `json:"aux_services"` // {ref:AuxService}
	AuxImgSrc      Set[string]             `json:"aux_img_src"`
}

type Service struct {
	Name              string                         `json:"name"`
	Image             string                         `json:"image"`
	RunConfig         RunConfig                      `json:"run_config"`
	BindMounts        map[string]BindMount           `json:"bind_mounts"`      // {mntPoint:BindMount}
	Tmpfs             map[string]TmpfsMount          `json:"tmpfs"`            // {mntPoint:TmpfsMount}
	Volumes           map[string]string              `json:"volumes"`          // {mntPoint:volName}
	HostResources     map[string]HostResTarget       `json:"host_resources"`   // {mntPoint:HostResTarget}
	SecretMounts      map[string]SecretTarget        `json:"secret_mounts"`    // {mntPoint:SecretTarget}
	SecretVars        map[string]SecretTarget        `json:"secret_vars"`      // {refVar:SecretTarget}
	Configs           map[string]string              `json:"configs"`          // {refVar:ref}
	SrvReferences     map[string]SrvRefTarget        `json:"srv_references"`   // {refVar:SrvRefTarget}
	HttpEndpoints     map[string]HttpEndpoint        `json:"http_endpoints"`   // {externalPath:HttpEndpoint}
	RequiredSrv       Set[string]                    `json:"required_srv"`     // {ref}
	RequiredBySrv     Set[string]                    `json:"required_by_srv"`  // {ref}
	ExtDependencies   map[string]ExtDependencyTarget `json:"ext_dependencies"` // {refVar:ExtDependencyTarget}
	Ports             []Port                         `json:"ports"`
	DeviceCGroupRules []string                       `json:"device_cgroup_rules"`
}

type AuxService struct {
	Name            string                         `json:"name"`
	RunConfig       RunConfig                      `json:"run_config"`
	BindMounts      map[string]BindMount           `json:"bind_mounts"`      // {mntPoint:BindMount}
	Tmpfs           map[string]TmpfsMount          `json:"tmpfs"`            // {mntPoint:TmpfsMount}
	Volumes         map[string]string              `json:"volumes"`          // {mntPoint:volName}
	Configs         map[string]string              `json:"configs"`          // {refVar:ref}
	SrvReferences   map[string]SrvRefTarget        `json:"srv_references"`   // {refVar:SrvRefTarget}
	ExtDependencies map[string]ExtDependencyTarget `json:"ext_dependencies"` // {refVar:ExtDependencyTarget}
}

type RunConfig struct {
	MaxRetries  uint          `json:"max_retries"`
	RunOnce     bool          `json:"run_once"`
	StopTimeout time.Duration `json:"stop_timeout"`
	StopSignal  *string       `json:"stop_signal"`
	PseudoTTY   bool          `json:"pseudo_tty"`
	Command     []string      `json:"command"`
}

type BindMount struct {
	Source   string `json:"source"`
	ReadOnly bool   `json:"read_only"`
}

type TmpfsMount struct {
	Size uint64      `json:"size"`
	Mode fs.FileMode `json:"mode"`
}

type HttpEndpoint struct {
	Name      *string               `json:"name"`
	Port      *int                  `json:"port"`
	Path      *string               `json:"path"` // internal path
	ProxyConf HttpEndpointProxyConf `json:"proxy_conf"`
	StringSub HttpEndpointStrSub    `json:"string_sub"`
}

type HttpEndpointStrSub struct {
	ReplaceOnce bool              `json:"replace_once"`
	MimeTypes   []string          `json:"mime_types"`
	Filters     map[string]string `json:"filters"`
}

type HttpEndpointProxyConf struct {
	Headers     map[string]string `json:"headers"`
	WebSocket   bool              `json:"websocket"`
	ReadTimeout time.Duration     `json:"read_timeout"`
}

type PortProtocol = string

type Port struct {
	Name     *string      `json:"name"`
	Number   uint         `json:"number"`
	Protocol PortProtocol `json:"protocol"`
	Bindings []uint       `json:"bindings"`
}

type ExtDependencyTarget struct {
	ID       string  `json:"id"`
	Service  string  `json:"service"`
	Template *string `json:"template"`
}

type HostResTarget struct {
	Ref      string `json:"ref"`
	ReadOnly bool   `json:"read_only"`
}

type Resource struct {
	Tags     Set[string] `json:"tags"`
	Required bool        `json:"required"`
}

type HostResource struct {
	Resource
}

type SecretTarget struct {
	Ref  string  `json:"ref"`
	Item *string `json:"item"`
}

type Secret struct {
	Resource
	Type string `json:"type"`
}

type SrvRefTarget struct {
	Ref      string  `json:"ref"`
	Template *string `json:"template"`
}

type Configs map[string]configValue

type configValue struct {
	Default   any               `json:"default"`
	Options   any               `json:"options"`
	OptExt    bool              `json:"opt_ext"`
	Type      string            `json:"type"`
	TypeOpt   ConfigTypeOptions `json:"type_opt"`
	DataType  DataType          `json:"data_type"`
	IsSlice   bool              `json:"is_slice"`
	Delimiter string            `json:"delimiter"`
	Required  bool              `json:"required"`
}

type ConfigTypeOptions map[string]configTypeOption

type DataType = string

type configTypeOption struct {
	Value    any      `json:"value"`
	DataType DataType `json:"data_type"`
}

type Input struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Group       *string `json:"group"`
}

type InputGroup struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Group       *string `json:"group"`
}

type Inputs struct {
	Resources map[string]Input      `json:"resources"` // {ref:Input}
	Secrets   map[string]Input      `json:"secrets"`   // {ref:Input}
	Configs   map[string]Input      `json:"configs"`   // {ref:Input}
	Groups    map[string]InputGroup `json:"groups"`    // {ref:InputGroup}
}
