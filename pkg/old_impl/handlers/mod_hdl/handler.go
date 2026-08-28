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
	"mgw-module-manager-migration/pkg/old_impl/libs/module_lib"
	"mgw-module-manager-migration/pkg/old_impl/model/pkg_model"
	"os"
	"path"
	"time"
)

type Handler struct {
	storageHandler storageHandler
	modFileHandler modFileHandler
	dbTimeout      time.Duration
	httpTimeout    time.Duration
	wrkSpcPath     string
}

func New(storageHandler storageHandler, modFileHandler modFileHandler, dbTimeout, httpTimeout time.Duration, workspacePath string) *Handler {
	return &Handler{
		storageHandler: storageHandler,
		modFileHandler: modFileHandler,
		dbTimeout:      dbTimeout,
		httpTimeout:    httpTimeout,
		wrkSpcPath:     workspacePath,
	}
}

func (h *Handler) Init(perm fs.FileMode) error {
	if !path.IsAbs(h.wrkSpcPath) {
		return fmt.Errorf("workspace path must be absolute")
	}
	if err := os.MkdirAll(h.wrkSpcPath, perm); err != nil {
		return err
	}
	return nil
}

func (h *Handler) List(ctx context.Context) (map[string]pkg_model.Module, error) {
	ctxWt, cf := context.WithTimeout(ctx, h.dbTimeout)
	defer cf()
	modMap, err := h.storageHandler.ListMod(ctxWt, pkg_model.ModFilter{}, false)
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
		modules[mod.ID] = mod
	}
	return modules, nil
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
