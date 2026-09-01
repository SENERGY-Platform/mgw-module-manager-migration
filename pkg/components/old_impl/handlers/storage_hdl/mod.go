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

package storage_hdl

import (
	"context"
	"database/sql"
	"errors"
	"mgw-module-manager-migration/pkg/components/old_impl/handlers/storage_hdl/dep_util"
	"mgw-module-manager-migration/pkg/components/old_impl/model"
	"mgw-module-manager-migration/pkg/components/old_impl/model/pkg_model"
	"strings"
	"time"
)

func (h *Handler) ListMod(ctx context.Context, filter pkg_model.ModFilter, dependencyInfo bool) (map[string]pkg_model.Module, error) {
	q := "SELECT `id`, `dir`, `modfile`, `added`, `updated` FROM `modules`"
	fc, val := genModFilter(filter)
	if fc != "" {
		q += fc
	}
	rows, err := h.db.QueryContext(ctx, q, val...)
	if err != nil {
		return nil, model.NewInternalError(err)
	}
	defer rows.Close()
	modules := make(map[string]pkg_model.Module)
	for rows.Next() {
		var id string
		var mod pkg_model.Module
		var at, ut []uint8
		if err = rows.Scan(&id, &mod.Dir, &mod.ModFile, &at, &ut); err != nil {
			return nil, model.NewInternalError(err)
		}
		if dependencyInfo {
			if mod.RequiredMod, err = dep_util.SelectRequiredMod(ctx, h.db, id); err != nil {
				return nil, model.NewInternalError(err)
			}
			if mod.ModRequiring, err = dep_util.SelectModRequiring(ctx, h.db, id); err != nil {
				return nil, model.NewInternalError(err)
			}
		}
		if mod.Added, err = time.Parse(tLayout, string(at)); err != nil {
			return nil, model.NewInternalError(err)
		}
		if mod.Updated, err = time.Parse(tLayout, string(ut)); err != nil {
			return nil, model.NewInternalError(err)
		}
		modules[id] = mod
	}
	return modules, nil
}

func genModFilter(filter pkg_model.ModFilter) (string, []any) {
	var fc []string
	var val []any
	if len(filter.IDs) > 0 {
		ids := removeDuplicates(filter.IDs)
		fc = append(fc, "`id` IN ("+strings.Repeat("?, ", len(ids)-1)+"?)")
		for _, id := range ids {
			val = append(val, id)
		}
	}
	if len(fc) > 0 {
		return " WHERE " + strings.Join(fc, " AND "), val
	}
	return "", nil
}
