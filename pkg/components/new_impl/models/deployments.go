/*
 * Copyright 2026 InfAI (CC SES)
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

package models

import (
	"time"
)

type DeploymentBase struct {
	Id            string
	ModuleId      string
	ModuleSource  string
	ModuleChannel string
	ModuleVersion string
	DirName       string
	FilesDirName  string
	Enabled       bool
	Created       time.Time
	Updated       time.Time
	Err           error
}

type DeploymentContainerBase struct {
	Name         string
	DeploymentId string
	Reference    string
	Alias        string
}

type DeploymentVolume struct {
	DeploymentId string
	Reference    string
	Name         string
}

type DeploymentHostResource struct {
	Id           string
	DeploymentId string
	Reference    string
}

type DeploymentSecret struct {
	Id           string
	DeploymentId string
	Reference    string
	Items        []DeploymentSecretItem
}

type DeploymentUserConfig struct {
	Id           string
	DeploymentId string
	Reference    string
	Value
}

type DeploymentGlobalConfig struct {
	Id           string
	DeploymentId string
	Reference    string
}

type DeploymentFile struct {
	DeploymentId string
	Reference    string
	Data         []byte
}

type DeploymentFileGroup struct {
	Id           string
	DeploymentId string
	Reference    string
	Files        []DeploymentFileGroupFile
}

type DeploymentFileGroupFile struct {
	Path   string
	Format string
	Data   []byte
}

type DeploymentSecretItem struct {
	Name    string
	AsMount bool
	AsEnv   bool
}

type DeploymentAdvertisement struct {
	Id           string
	DeploymentId string
	ModuleId     string
	Reference    string
	Timestamp    time.Time
	Items        map[string]string
}
