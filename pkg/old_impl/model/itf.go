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
	"context"
)

type Api interface {
	AddModule(ctx context.Context, mID, version string) (string, error)
	GetModules(ctx context.Context, filter ModFilter) (map[string]Module, error)
	GetModule(ctx context.Context, mID string) (Module, error)
	DeleteModule(ctx context.Context, mID string, force bool) (string, error)
	GetModuleDeployTemplate(ctx context.Context, mID string) (ModDeployTemplate, error)
	CheckModuleUpdates(ctx context.Context) (string, error)
	GetModuleUpdates(ctx context.Context) (map[string]ModUpdate, error)
	GetModuleUpdate(ctx context.Context, mID string) (ModUpdate, error)
	PrepareModuleUpdate(ctx context.Context, mID, version string) (string, error)
	CancelPendingModuleUpdate(ctx context.Context, mID string) error
	UpdateModule(ctx context.Context, mID string, depInput DepInput, dependencies map[string]DepInput) (string, error)
	GetModuleUpdateTemplate(ctx context.Context, id string) (ModUpdateTemplate, error)
	CreateDeployment(ctx context.Context, mID string, depInput DepInput, dependencies map[string]DepInput) (string, error)
	GetDeployments(ctx context.Context, filter DepFilter, assets, containerInfo bool) (map[string]Deployment, error)
	GetDeployment(ctx context.Context, dID string, assets, containerInfo bool) (Deployment, error)
	UpdateDeployment(ctx context.Context, dID string, depInput DepInput) (string, error)
	DeleteDeployment(ctx context.Context, dID string, force bool) (string, error)
	DeleteDeployments(ctx context.Context, filter DepFilter, force bool) (string, error)
	StartDeployment(ctx context.Context, dID string, dependencies bool) (string, error)
	StartDeployments(ctx context.Context, filter DepFilter, dependencies bool) (string, error)
	StopDeployment(ctx context.Context, dID string, force bool) (string, error)
	StopDeployments(ctx context.Context, filter DepFilter, force bool) (string, error)
	RestartDeployment(ctx context.Context, dID string) (string, error)
	RestartDeployments(ctx context.Context, filter DepFilter) (string, error)
	GetDeploymentUpdateTemplate(ctx context.Context, dID string) (DepUpdateTemplate, error)
	AuxDeploymentApi
	DepAdvertisementApi
}

type AuxDeploymentApi interface {
	GetAuxDeployments(ctx context.Context, dID string, filter AuxDepFilter, assets, containerInfo bool) (map[string]AuxDeployment, error)
	GetAuxDeployment(ctx context.Context, dID, aID string, assets, containerInfo bool) (AuxDeployment, error)
	CreateAuxDeployment(ctx context.Context, dID string, auxDepInput AuxDepReq, forcePullImg bool) (string, error)
	UpdateAuxDeployment(ctx context.Context, dID, aID string, auxDepInput AuxDepReq, incremental, forcePullImg bool) (string, error)
	DeleteAuxDeployment(ctx context.Context, dID, aID string, force bool) (string, error)
	DeleteAuxDeployments(ctx context.Context, dID string, filter AuxDepFilter, force bool) (string, error)
	StartAuxDeployment(ctx context.Context, dID, aID string) (string, error)
	StartAuxDeployments(ctx context.Context, dID string, filter AuxDepFilter) (string, error)
	StopAuxDeployment(ctx context.Context, dID, aID string) (string, error)
	StopAuxDeployments(ctx context.Context, dID string, filter AuxDepFilter) (string, error)
	RestartAuxDeployment(ctx context.Context, dID, aID string) (string, error)
	RestartAuxDeployments(ctx context.Context, dID string, filter AuxDepFilter) (string, error)
}

type DepAdvertisementApi interface {
	QueryDepAdvertisements(ctx context.Context, filter DepAdvFilter) ([]DepAdvertisement, error)
	GetDepAdvertisement(ctx context.Context, dID, ref string) (DepAdvertisement, error)
	GetDepAdvertisements(ctx context.Context, dID string) (map[string]DepAdvertisement, error)
	PutDepAdvertisement(ctx context.Context, dID string, adv DepAdvertisementBase) error
	PutDepAdvertisements(ctx context.Context, dID string, ads map[string]DepAdvertisementBase) error
	DeleteDepAdvertisement(ctx context.Context, dID, ref string) error
	DeleteDepAdvertisements(ctx context.Context, dID string) error
}
