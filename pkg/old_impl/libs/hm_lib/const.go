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

package hm_lib

const (
	HeaderRequestID = "X-Request-ID"
	HeaderApiVer    = "X-Api-Version"
	HeaderSrvName   = "X-Service"
)

const (
	HostInfoPath      = "host-info"
	HostNetPath       = "network"
	HostOsPath        = "os"
	HostHwPath        = "hardware"
	HostResourcesPath = "host-resources"
	SrvInfoPath       = "info"
	RestrictedPath    = "restricted"
	HostAppsPath      = "applications"
	BlacklistsPath    = "blacklists"
	NetInterfacesPath = "net-interfaces"
	NetRangesPath     = "net-ranges"
	MDNSDiscoveryPath = "mdns-discovery"
)

const (
	SerialDevice ResourceType = "serial"
	Application  ResourceType = "app"
)
