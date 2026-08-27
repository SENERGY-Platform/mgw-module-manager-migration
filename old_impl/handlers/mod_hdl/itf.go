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

package mod_hdl

import (
	"context"
	"io/fs"
	"mgw-module-manager-migration/old_impl/libs/module_lib"
	"mgw-module-manager-migration/old_impl/model/pkg_model"
	"mgw-module-manager-migration/old_impl/util/dir_fs"
)

type StorageHandler interface {
	ListMod(ctx context.Context, filter pkg_model.ModFilter, dependencyInfo bool) (map[string]pkg_model.Module, error)
	ReadMod(ctx context.Context, mID string, dependencyInfo bool) (pkg_model.Module, error)
}

type ModFileHandler interface {
	GetModule(file fs.File) (*module_lib.Module, error)
	GetModFile(dir dir_fs.DirFS) (fs.File, string, error)
}
