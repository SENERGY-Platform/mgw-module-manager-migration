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
	"mgw-module-manager-migration/pkg/old_impl/libs/cew_lib"
	"mgw-module-manager-migration/pkg/old_impl/model"
	"mgw-module-manager-migration/pkg/old_impl/util/naming_hdl"
	"os"
	"path"
	"time"
)

type Handler struct {
	storageHandler storageHandler
	cewClient      cewClient
	dbTimeout      time.Duration
	httpTimeout    time.Duration
	wrkSpcPath     string
	managerID      string
}

func New(
	storageHandler storageHandler,
	cewClient cewClient,
	dbTimeout time.Duration,
	httpTimeout time.Duration,
	workspacePath string,
	managerID string,
) *Handler {
	return &Handler{
		storageHandler: storageHandler,
		cewClient:      cewClient,
		dbTimeout:      dbTimeout,
		httpTimeout:    httpTimeout,
		wrkSpcPath:     workspacePath,
		managerID:      managerID,
	}
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

func (h *Handler) List(ctx context.Context) (map[string]model.Deployment, error) {
	ctxWt, cf := context.WithTimeout(ctx, h.dbTimeout)
	defer cf()
	deployments, err := h.storageHandler.ListDep(ctxWt, model.DepFilter{}, false, true, true)
	if err != nil {
		return nil, err
	}
	if len(deployments) > 0 {
		ctxWt2, cf2 := context.WithTimeout(ctx, h.httpTimeout)
		defer cf2()
		ctrList, err := h.cewClient.GetContainers(ctxWt2, cew_lib.ContainerFilter{Labels: map[string]string{naming_hdl.ManagerIDLabel: h.managerID}})
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not retrieve containers: %s\n", err.Error())
			return nil, err
		}
		ctrMap := make(map[string]cew_lib.Container)
		for _, ctr := range ctrList {
			ctrMap[ctr.ID] = ctr
		}
		ctxWt3, cf3 := context.WithTimeout(ctx, h.httpTimeout)
		defer cf3()
		volList, err := h.cewClient.GetVolumes(ctxWt3, cew_lib.VolumeFilter{Labels: map[string]string{naming_hdl.ManagerIDLabel: h.managerID}})
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not retrieve volumes: %s\n", err.Error())
			return nil, err
		}
		volMap := getVolMap(volList)
		withCtrAndVols := make(map[string]model.Deployment)
		for dID, deployment := range deployments {
			deployment.Containers = getDepCtrInfo(dID, deployment.Containers, ctrMap)
			deployment.Volumes = volMap[dID]
			withCtrAndVols[dID] = deployment
		}
		return withCtrAndVols, nil
	}
	return deployments, nil
}

func getVolMap(vols []cew_lib.Volume) map[string]map[string]string {
	volMap := make(map[string]map[string]string)
	for _, vol := range vols {
		depId, ok := vol.Labels[naming_hdl.DeploymentIDLabel]
		if !ok {
			fmt.Fprintf(os.Stderr, "volume without deployment id: %s\n", vol.Name)
			continue
		}
		ref, ok := vol.Labels[naming_hdl.VolumeRefLabel]
		if !ok {
			fmt.Fprintf(os.Stderr, "volume without reference: %s\n", vol.Name)
			continue
		}
		depVols, ok := volMap[depId]
		if !ok {
			depVols = make(map[string]string)
			volMap[depId] = depVols
		}
		depVols[ref] = vol.Name
	}
	return volMap
}

func getDepCtrInfo(dID string, depContainers map[string]model.DepContainer, ctrMap map[string]cew_lib.Container) map[string]model.DepContainer {
	withCtrInfo := make(map[string]model.DepContainer)
	for ref, depCtr := range depContainers {
		ctr, ok := ctrMap[depCtr.ID]
		if !ok {
			fmt.Fprintf(os.Stderr, "deployment '%s' missing container '%s'\n", dID, depCtr.ID)
		} else {
			depCtr.Info = &model.ContainerInfo{
				Name:    ctr.Name,
				ImageID: ctr.ImageID,
				State:   ctr.State,
			}
		}
		withCtrInfo[ref] = depCtr
	}
	return withCtrInfo
}
