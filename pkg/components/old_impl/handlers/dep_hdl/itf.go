/*
 * Copyright 2025 InfAI (CC SES)
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

package dep_hdl

import (
	"context"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/cew_lib"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/hm_lib"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/module_lib"
	"mgw-module-manager-migration/pkg/components/old_impl/model"
)

type storageHandler interface {
	ListDep(ctx context.Context, filter model.DepFilter, dependencyInfo, assets, containers bool) (map[string]model.Deployment, error)
}

type cfgValidationHandler interface {
	ValidateValue(cType string, cTypeOpt module_lib.ConfigTypeOptions, value any, isSlice bool, dataType module_lib.DataType) error
	ValidateValInOpt(cOpt any, value any, isSlice bool, dataType module_lib.DataType) error
}

type cewClient interface {
	GetContainers(ctx context.Context, filter cew_lib.ContainerFilter) ([]cew_lib.Container, error)
	GetVolumes(ctx context.Context, filter cew_lib.VolumeFilter) ([]cew_lib.Volume, error)
}

type hmClient interface {
	GetHostResource(ctx context.Context, rID string) (hm_lib.HostResource, error)
}
