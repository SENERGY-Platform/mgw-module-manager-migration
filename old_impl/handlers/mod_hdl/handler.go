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

package mod_hdl

import (
	"context"
	"fmt"
	"io/fs"
	"mgw-module-manager-migration/old_impl/libs/cew_lib"
	"mgw-module-manager-migration/old_impl/libs/module_lib"
	"mgw-module-manager-migration/old_impl/model"
	"mgw-module-manager-migration/old_impl/model/pkg_model"
	"os"
	"path"
	"sync"
	"time"
)

type Handler struct {
	storageHandler StorageHandler
	modFileHandler ModFileHandler
	cewClient      cew_lib.Api
	dbTimeout      time.Duration
	httpTimeout    time.Duration
	wrkSpcPath     string
	mu             sync.RWMutex
}

func New(storageHandler StorageHandler, modFileHandler ModFileHandler, cewClient cew_lib.Api, dbTimeout, httpTimeout time.Duration, workspacePath string) *Handler {
	return &Handler{
		storageHandler: storageHandler,
		modFileHandler: modFileHandler,
		cewClient:      cewClient,
		dbTimeout:      dbTimeout,
		httpTimeout:    httpTimeout,
		wrkSpcPath:     workspacePath,
	}
}

func (h *Handler) Init(perm fs.FileMode) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if !path.IsAbs(h.wrkSpcPath) {
		return fmt.Errorf("workspace path must be absolute")
	}
	if err := os.MkdirAll(h.wrkSpcPath, perm); err != nil {
		return err
	}
	return nil
}

func (h *Handler) List(ctx context.Context, filter model.ModFilter, dependencyInfo bool) (map[string]pkg_model.Module, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ctxWt, cf := context.WithTimeout(ctx, h.dbTimeout)
	defer cf()
	modMap, err := h.storageHandler.ListMod(ctxWt, pkg_model.ModFilter{IDs: filter.IDs}, dependencyInfo)
	if err != nil {
		return nil, err
	}
	modules := make(map[string]pkg_model.Module)
	for _, mod := range modMap {
		mod.Module.Module, err = h.readModule(mod.Dir, mod.ModFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		mod.Path = h.wrkSpcPath
		if filterMod(filter, mod.Module.Module) {
			modules[mod.ID] = mod
		}
	}
	return modules, nil
}

func (h *Handler) Get(ctx context.Context, mID string, dependencyInfo bool) (pkg_model.Module, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ctxWt, cf := context.WithTimeout(ctx, h.dbTimeout)
	defer cf()
	mod, err := h.storageHandler.ReadMod(ctxWt, mID, dependencyInfo)
	if err != nil {
		return pkg_model.Module{}, err
	}
	mod.Module.Module, err = h.readModule(mod.Dir, mod.ModFile)
	if err != nil {
		return pkg_model.Module{}, model.NewInternalError(err)
	}
	mod.Path = h.wrkSpcPath
	return mod, nil
}

func (h *Handler) GetTree(ctx context.Context, mID string) (map[string]pkg_model.Module, error) {
	mod, err := h.storageHandler.ReadMod(ctx, mID, true)
	if err != nil {
		return nil, err
	}
	mod.Module.Module, err = h.readModule(mod.Dir, mod.ModFile)
	if err != nil {
		return nil, model.NewInternalError(err)
	}
	mod.Path = h.wrkSpcPath
	tree := map[string]pkg_model.Module{mod.ID: mod}
	if err = h.appendModTree(ctx, mod, tree); err != nil {
		return nil, err
	}
	return tree, nil
}

func (h *Handler) AppendModTree(ctx context.Context, tree map[string]pkg_model.Module) error {
	for _, mod := range tree {
		if err := h.appendModTree(ctx, mod, tree); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) appendModTree(ctx context.Context, mod pkg_model.Module, tree map[string]pkg_model.Module) error {
	for _, mID := range mod.RequiredMod {
		if _, ok := tree[mID]; !ok {
			m, err := h.storageHandler.ReadMod(ctx, mID, true)
			if err != nil {
				return err
			}
			m.Module.Module, err = h.readModule(m.Dir, m.ModFile)
			if err != nil {
				return err
			}
			m.Path = h.wrkSpcPath
			tree[mID] = m
			if err = h.appendModTree(ctx, m, tree); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Handler) readModule(dir, modFile string) (*module_lib.Module, error) {
	f, err := os.Open(path.Join(h.wrkSpcPath, dir, modFile))
	if err != nil {
		return nil, err
	}
	m, err := h.modFileHandler.GetModule(f)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func filterMod(filter model.ModFilter, m *module_lib.Module) bool {
	if filter.Name != "" {
		if m.Name != filter.Name {
			return false
		}
	}
	if filter.Type != "" {
		if m.Type != filter.Type {
			return false
		}
	}
	if filter.DeploymentType != "" {
		if m.DeploymentType != filter.DeploymentType {
			return false
		}
	}
	if filter.Author != "" {
		if m.Author != filter.Author {
			return false
		}
	}
	if len(filter.Tags) > 0 {
		var ok bool
		for tag := range filter.Tags {
			if _, ok = m.Tags[tag]; ok {
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}
