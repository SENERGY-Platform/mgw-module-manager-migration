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
	"io/fs"
	"net"
	"time"
)

// Image ---------------------------------------------------------------------------------------

type Image struct {
	ID      string            `json:"id"`
	Created time.Time         `json:"created"`
	Size    int64             `json:"size"`
	Arch    string            `json:"arch"`
	Tags    []string          `json:"tags"`
	Digests []string          `json:"digests"`
	Labels  map[string]string `json:"labels"`
}

type ImageFilter struct {
	Name   string
	Tag    string
	Labels map[string]string
}

type ImageRequest struct {
	Image string `json:"image"`
}

// Network -------------------------------------------------------------------------------------

type IPAddr net.IP

type Subnet struct {
	Prefix IPAddr `json:"prefix"`
	Bits   int    `json:"bits"`
}

type NetworkType = string

type Network struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Type    NetworkType `json:"type"`
	Subnet  Subnet      `json:"subnet"`
	Gateway IPAddr      `json:"gateway"`
}

type PortType = string

type Port struct {
	Number   int           `json:"number"`
	Protocol PortType      `json:"protocol"`
	Bindings []PortBinding `json:"bindings"`
}

type PortBinding struct {
	Number    int    `json:"number"`
	Interface IPAddr `json:"interface"`
}

// Container -----------------------------------------------------------------------------------

type MountType = string

type Mount struct {
	Type     MountType         `json:"type"`
	Source   string            `json:"source"`
	Target   string            `json:"target"`
	ReadOnly bool              `json:"read_only"`
	Labels   map[string]string `json:"labels,omitempty"`
	Size     int64             `json:"size,omitempty"`
	Mode     fs.FileMode       `json:"mode,omitempty"`
}

type Device struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	ReadOnly bool   `json:"read_only"`
}

type RestartStrategy = string

type RunConfig struct {
	RestartStrategy RestartStrategy `json:"restart_strategy"`
	Retries         *int            `json:"retries"`
	RemoveAfterRun  bool            `json:"remove_after_run"`
	StopTimeout     *time.Duration  `json:"stop_timeout"`
	StopSignal      *string         `json:"stop_signal"`
	PseudoTTY       bool            `json:"pseudo_tty"`
	Command         []string        `json:"command"`
}

type ContainerState = string

type ContainerHealth = string

type Container struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	State             ContainerState    `json:"state"`
	Health            *ContainerHealth  `json:"health"`
	Created           time.Time         `json:"created"`
	Started           *time.Time        `json:"started"`
	Image             string            `json:"image"`
	ImageID           string            `json:"image_id"`
	EnvVars           map[string]string `json:"env_vars"`
	Labels            map[string]string `json:"labels"`
	Mounts            []Mount           `json:"mounts"`
	Devices           []Device          `json:"devices"`
	DeviceCGroupRules []string          `json:"device_cgroup_rules"`
	Ports             []Port            `json:"ports"`
	Networks          []ContainerNet    `json:"networks"`
	RunConfig         RunConfig         `json:"run_config"`
}

type ContainerNet struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	DomainNames []string `json:"domain_names"`
	Gateway     IPAddr   `json:"gateway"`
	IPAddress   IPAddr   `json:"ip_address"`
	MacAddress  string   `json:"mac_address"`
}

type ContainerFilter struct {
	Name   string
	State  ContainerState
	Labels map[string]string
}

type LogFilter struct {
	MaxLines int
	Since    time.Time
	Until    time.Time
}

type ExecConfig struct {
	Tty     bool
	EnvVars map[string]string
	WorkDir string
	Cmd     []string
}

// Volume --------------------------------------------------------------------------------------

type Volume struct {
	Name    string            `json:"name"`
	Created time.Time         `json:"created"`
	Labels  map[string]string `json:"labels"`
}

type VolumeFilter struct {
	Labels map[string]string
}

// Error -----------------------------------------------------------------------------------------

type cError struct {
	err error
}

type InternalError struct {
	cError
}

type NotFoundError struct {
	cError
}

type InvalidInputError struct {
	cError
}
