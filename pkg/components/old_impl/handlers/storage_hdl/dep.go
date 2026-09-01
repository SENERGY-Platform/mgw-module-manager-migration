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
	"mgw-module-manager-migration/pkg/components/old_impl/handlers/storage_hdl/dep_util"
	"mgw-module-manager-migration/pkg/components/old_impl/model"
	"strings"
	"time"
)

func (h *Handler) ListDep(ctx context.Context, filter model.DepFilter, dependencyInfo, assets, containers bool) (map[string]model.Deployment, error) {
	q := "SELECT `id`, `mod_id`, `mod_ver`, `name`, `dir`, `enabled`, `indirect`, `created`, `updated` FROM `bk_deployments`"
	fc, val := genDepFilter(filter)
	if fc != "" {
		q += fc
	}
	rows, err := h.db.QueryContext(ctx, q, val...)
	if err != nil {
		return nil, model.NewInternalError(err)
	}
	defer rows.Close()
	deployments := make(map[string]model.Deployment)
	for rows.Next() {
		var deployment model.Deployment
		var depModule model.DepModule
		var ct, ut []uint8
		if err = rows.Scan(&deployment.ID, &depModule.ID, &depModule.Version, &deployment.Name, &deployment.Dir, &deployment.Enabled, &deployment.Indirect, &ct, &ut); err != nil {
			return nil, model.NewInternalError(err)
		}
		deployment.Module = depModule
		if deployment.Created, err = time.Parse(tLayout, string(ct)); err != nil {
			return nil, model.NewInternalError(err)
		}
		if deployment.Updated, err = time.Parse(tLayout, string(ut)); err != nil {
			return nil, model.NewInternalError(err)
		}
		if dependencyInfo {
			if deployment.RequiredDep, err = dep_util.SelectRequiredDep(ctx, h.db, deployment.ID); err != nil {
				return nil, model.NewInternalError(err)
			}
			if deployment.DepRequiring, err = dep_util.SelectDepRequiring(ctx, h.db, deployment.ID); err != nil {
				return nil, model.NewInternalError(err)
			}
		}
		if assets {
			if deployment.DepAssets, err = readDepAssets(ctx, h.db, deployment.ID); err != nil {
				return nil, model.NewInternalError(err)
			}
		}
		if containers {
			if deployment.Containers, err = dep_util.SelectDepContainers(ctx, h.db, deployment.ID); err != nil {
				return nil, model.NewInternalError(err)
			}
		}
		deployments[deployment.ID] = deployment
	}
	if err = rows.Err(); err != nil {
		return nil, model.NewInternalError(err)
	}
	return deployments, nil
}

func readDepAssets(ctx context.Context, db *sql.DB, id string) (model.DepAssets, error) {
	hostRes, err := dep_util.SelectHostResources(ctx, db, id)
	if err != nil {
		return model.DepAssets{}, err
	}
	secrets, err := dep_util.SelectSecrets(ctx, db, id)
	if err != nil {
		return model.DepAssets{}, err
	}
	configs := make(map[string]model.DepConfig)
	err = dep_util.SelectConfigs(ctx, db, id, configs)
	if err != nil {
		return model.DepAssets{}, err
	}
	err = dep_util.SelectListConfigs(ctx, db, id, configs)
	if err != nil {
		return model.DepAssets{}, err
	}
	return model.DepAssets{
		HostResources: hostRes,
		Secrets:       secrets,
		Configs:       configs,
	}, nil
}

func genDepFilter(filter model.DepFilter) (string, []any) {
	var fc []string
	var val []any
	if len(filter.IDs) > 0 {
		ids := removeDuplicates(filter.IDs)
		fc = append(fc, "`id` IN ("+strings.Repeat("?, ", len(ids)-1)+"?)")
		for _, id := range ids {
			val = append(val, id)
		}
	}
	if filter.ModuleID != "" {
		fc = append(fc, "`mod_id` = ?")
		val = append(val, filter.ModuleID)
	}
	if filter.Name != "" {
		fc = append(fc, "`name` = ?")
		val = append(val, filter.Name)
	}
	if filter.Enabled != 0 {
		fc = append(fc, "`enabled` = ?")
		if filter.Enabled > 0 {
			val = append(val, true)
		} else {
			val = append(val, false)
		}
	}
	if filter.Indirect != 0 {
		fc = append(fc, "`indirect` = ?")
		if filter.Indirect > 0 {
			val = append(val, true)
		} else {
			val = append(val, false)
		}
	}
	if len(fc) > 0 {
		return " WHERE " + strings.Join(fc, " AND "), val
	}
	return "", nil
}

func removeDuplicates(sl []string) []string {
	if len(sl) < 2 {
		return sl
	}
	set := make(map[string]struct{})
	var sl2 []string
	for _, s := range sl {
		if _, ok := set[s]; !ok {
			sl2 = append(sl2, s)
		}
		set[s] = struct{}{}
	}
	return sl2
}
