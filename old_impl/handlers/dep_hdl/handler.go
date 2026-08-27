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

package dep_hdl

import (
	"context"
	"fmt"
	"io/fs"
	"mgw-module-manager-migration/old_impl/handlers/context_hdl"
	"mgw-module-manager-migration/old_impl/libs/cew_lib"
	"mgw-module-manager-migration/old_impl/libs/hm_lib"
	"mgw-module-manager-migration/old_impl/model"
	"mgw-module-manager-migration/old_impl/util"
	"mgw-module-manager-migration/old_impl/util/naming_hdl"
	"os"
	"path"
	"time"
)

type Handler struct {
	storageHandler StorageHandler
	cfgVltHandler  CfgValidationHandler
	cewClient      cew_lib.Api
	hmClient       hm_lib.Api
	dbTimeout      time.Duration
	httpTimeout    time.Duration
	wrkSpcPath     string
	depHostPath    string
	secHostPath    string
	managerID      string
	coreID         string
	moduleNet      string
}

func New(storageHandler StorageHandler, cfgVltHandler CfgValidationHandler, cewClient cew_lib.Api, hmClient hm_lib.Api, dbTimeout time.Duration, httpTimeout time.Duration, workspacePath, depHostPath, secHostPath, managerID, moduleNet, coreID string) *Handler {
	return &Handler{
		storageHandler: storageHandler,
		cfgVltHandler:  cfgVltHandler,
		cewClient:      cewClient,
		hmClient:       hmClient,
		dbTimeout:      dbTimeout,
		httpTimeout:    httpTimeout,
		wrkSpcPath:     workspacePath,
		depHostPath:    depHostPath,
		secHostPath:    secHostPath,
		managerID:      managerID,
		coreID:         coreID,
		moduleNet:      moduleNet,
	}
}

type secretVariant struct {
	Item  *string
	Path  string
	AsEnv bool
	Value string
}

type secret struct {
	ID       string
	Variants map[string]secretVariant
}

func (h *Handler) InitWorkspace(perm fs.FileMode) error {
	if !path.IsAbs(h.wrkSpcPath) {
		return fmt.Errorf("workspace path must be absolute")
	}
	if err := os.MkdirAll(h.wrkSpcPath, perm); err != nil {
		return err
	}
	return nil
}

func (h *Handler) List(ctx context.Context, filter model.DepFilter, dependencyInfo, assets, containers, containerInfo bool) (map[string]model.Deployment, error) {
	ctxWt, cf := context.WithTimeout(ctx, h.dbTimeout)
	defer cf()
	if containerInfo {
		containers = true
	}
	deployments, err := h.storageHandler.ListDep(ctxWt, filter, dependencyInfo, assets, containers)
	if err != nil {
		return nil, err
	}
	if containerInfo && len(deployments) > 0 {
		ctxWt2, cf2 := context.WithTimeout(ctx, h.dbTimeout)
		defer cf2()
		ctrList, err := h.cewClient.GetContainers(ctxWt2, cew_lib.ContainerFilter{Labels: map[string]string{naming_hdl.ManagerIDLabel: h.managerID}})
		if err != nil {
			util.Logger.Errorf("could not retrieve containers: %s", err.Error())
			return deployments, nil
		}
		ctrMap := make(map[string]cew_lib.Container)
		for _, ctr := range ctrList {
			ctrMap[ctr.ID] = ctr
		}
		withCtrInfo := make(map[string]model.Deployment)
		for dID, deployment := range deployments {
			if deployment.Enabled {
				deployment.State, deployment.Containers = getDepHealthAndCtrInfo(dID, deployment.Containers, ctrMap)
			}
			withCtrInfo[dID] = deployment
		}
		return withCtrInfo, nil
	}
	return deployments, nil
}

func (h *Handler) Get(ctx context.Context, id string, dependencyInfo, assets, containers, containerInfo bool) (model.Deployment, error) {
	ctxWt, cf := context.WithTimeout(ctx, h.dbTimeout)
	defer cf()
	if containerInfo {
		containers = true
	}
	deployment, err := h.storageHandler.ReadDep(ctxWt, id, dependencyInfo, assets, containers)
	if err != nil {
		return model.Deployment{}, err
	}
	if containerInfo && deployment.Enabled {
		ctxWt2, cf2 := context.WithTimeout(ctx, h.dbTimeout)
		defer cf2()
		ctrList, err := h.cewClient.GetContainers(ctxWt2, cew_lib.ContainerFilter{Labels: map[string]string{naming_hdl.ManagerIDLabel: h.managerID, naming_hdl.DeploymentIDLabel: id}})
		if err != nil {
			util.Logger.Errorf("could not retrieve containers: %s", err.Error())
			return deployment, nil
		}
		ctrMap := make(map[string]cew_lib.Container)
		for _, ctr := range ctrList {
			ctrMap[ctr.ID] = ctr
		}
		deployment.State, deployment.Containers = getDepHealthAndCtrInfo(deployment.ID, deployment.Containers, ctrMap)
	}
	return deployment, nil
}

func (h *Handler) getModDependencyDeployments(ctx context.Context, modDependencies map[string]string) (map[string]model.Deployment, error) {
	ch := context_hdl.New()
	defer ch.CancelAll()
	m := make(map[string]model.Deployment)
	for mID := range modDependencies {
		deployments, err := h.storageHandler.ListDep(ch.Add(context.WithTimeout(ctx, h.dbTimeout)), model.DepFilter{ModuleID: mID}, false, false, true)
		if err != nil {
			return nil, err
		}
		if len(deployments) == 0 {
			return nil, model.NewInternalError(fmt.Errorf("dependency '%s' not deployed", mID))
		}
		if len(deployments) > 1 {
			return nil, model.NewInternalError(fmt.Errorf("dependency '%s' has multiple deployments", mID))
		}
		for _, dep := range deployments {
			m[mID] = dep
			break
		}
	}
	return m, nil
}

func getDepHealthAndCtrInfo(dID string, depContainers map[string]model.DepContainer, ctrMap map[string]cew_lib.Container) (*model.HealthState, map[string]model.DepContainer) {
	var state model.HealthState
	withCtrInfo := make(map[string]model.DepContainer)
	for ref, depCtr := range depContainers {
		ctr, ok := ctrMap[depCtr.ID]
		if !ok {
			state = model.DepUnhealthy
			util.Logger.Warningf("deployment '%s' missing container '%s'", dID, depCtr.ID)
		} else {
			if state == "" {
				if ctr.Health != nil {
					switch *ctr.Health {
					case cew_lib.TransitionState:
						state = model.DepTrans
					case cew_lib.UnhealthyState:
						state = model.DepUnhealthy
					}
				} else {
					switch ctr.State {
					case cew_lib.InitState, cew_lib.RestartingState, cew_lib.RemovingState:
						state = model.DepTrans
					case cew_lib.StoppedState, cew_lib.DeadState, cew_lib.PausedState:
						state = model.DepUnhealthy
					}
				}
			}
			depCtr.Info = &model.ContainerInfo{
				ImageID: ctr.ImageID,
				State:   ctr.State,
			}
		}
		withCtrInfo[ref] = depCtr
	}
	if state == "" {
		state = model.DepHealthy
	}
	return &state, withCtrInfo
}
