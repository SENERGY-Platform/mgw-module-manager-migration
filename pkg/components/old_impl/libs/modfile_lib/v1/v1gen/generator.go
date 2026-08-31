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

package v1gen

import (
	"errors"
	model2 "mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/v1/model"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/v1/v1gen/configs"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/v1/v1gen/generic"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/v1/v1gen/inputs"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/v1/v1gen/mounts"
	services2 "mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/v1/v1gen/services"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/module_lib"
)

func generator(f any) (*module_lib.Module, error) {
	mf, ok := f.(*model2.ModFile)
	if !ok {
		return nil, errors.New("invalid type")
	}
	mCs, err := configs.GenConfigs(mf.Configs)
	if err != nil {
		return nil, err
	}
	mSs, err := services2.GenServices(mf.Services)
	if err != nil {
		return nil, err
	}
	mAs, err := services2.GenAuxServices(mf.AuxServices)
	if err != nil {
		return nil, err
	}
	err = services2.SetSrvReferences(mf.ServiceReferences, mSs)
	if err != nil {
		return nil, err
	}
	err = services2.SetAuxSrvReferences(mf.ServiceReferences, mAs)
	if err != nil {
		return nil, err
	}
	err = services2.SetVolumes(mf.Volumes, mSs)
	if err != nil {
		return nil, err
	}
	err = services2.SetAuxVolumes(mf.Volumes, mAs)
	if err != nil {
		return nil, err
	}
	err = services2.SetExtDependencies(mf.Dependencies, mSs)
	if err != nil {
		return nil, err
	}
	err = services2.SetAuxExtDependencies(mf.Dependencies, mAs)
	if err != nil {
		return nil, err
	}
	err = services2.SetHostResources(mf.HostResources, mSs)
	if err != nil {
		return nil, err
	}
	err = services2.SetSecrets(mf.Secrets, mSs)
	if err != nil {
		return nil, err
	}
	err = services2.SetConfigs(mf.Configs, mSs)
	if err != nil {
		return nil, err
	}
	err = services2.SetAuxConfigs(mf.Configs, mAs)
	if err != nil {
		return nil, err
	}
	return &module_lib.Module{
		ID:             mf.ID,
		Name:           mf.Name,
		Description:    mf.Description,
		Tags:           generic.GenStringSet(mf.Tags),
		License:        mf.License,
		Author:         mf.Author,
		Version:        mf.Version,
		Type:           mf.Type,
		DeploymentType: mf.DeploymentType,
		Architectures:  generic.GenStringSet(mf.Architectures),
		Services:       mSs,
		AuxServices:    mAs,
		AuxImgSrc:      generic.GenStringSet(mf.AuxImageSources),
		Volumes:        mounts.GenVolumes(mf.Volumes),
		Dependencies:   mounts.GenDependencies(mf.Dependencies),
		HostResources:  mounts.GenHostResources(mf.HostResources),
		Secrets:        mounts.GenSecrets(mf.Secrets),
		Configs:        mCs,
		Inputs: module_lib.Inputs{
			Resources: inputs.GenInputs(mf.HostResources),
			Secrets:   inputs.GenInputs(mf.Secrets),
			Configs:   inputs.GenInputs(mf.Configs),
			Groups:    inputs.GenInputGroups(mf.InputGroups),
		},
	}, nil
}

func GetGenerator() (string, func(any) (*module_lib.Module, error)) {
	return model2.Version, generator
}
