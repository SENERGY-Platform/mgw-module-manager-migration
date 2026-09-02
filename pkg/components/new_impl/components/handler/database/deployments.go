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

package database

import (
	"context"
	"database/sql"
	"fmt"
	helper_slices "mgw-module-manager-migration/pkg/components/new_impl/components/helper/slices"
	"mgw-module-manager-migration/pkg/components/new_impl/models"
	"os"
	"strings"
	"time"
)

func (h *Handler) ReadModules(ctx context.Context, filter models.ModulesFilter) (map[string]models.DatabaseModule, error) {
	fc, val := genModulesFilter(filter)
	rows, err := h.sqlDB.QueryContext(
		ctx,
		"SELECT id, dir, source, channel, added, updated FROM modules"+fc+";",
		val...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	mods := make(map[string]models.DatabaseModule)
	for rows.Next() {
		var mod models.DatabaseModule
		var at, ut []uint8
		err = rows.Scan(&mod.Id, &mod.DirName, &mod.Source, &mod.Channel, &at, &ut)
		if err != nil {
			return nil, err
		}
		if mod.Added, err = time.Parse(timeLayout, string(at)); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "read modules", mod.Id, err)
		}
		if mod.Updated, err = time.Parse(timeLayout, string(ut)); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "read modules", mod.Id, err)
		}
		mods[mod.Id] = mod
	}
	return mods, nil
}

func (h *Handler) CreateDeployment(
	ctx context.Context,
	tx *sql.Tx,
	deployment models.DeploymentBase,
	hostResources []models.DeploymentHostResource,
	secrets []models.DeploymentSecret,
	userConfigs []models.DeploymentUserConfig,
	globalConfigs []models.DeploymentGlobalConfig,
	files []models.DeploymentFile,
	fileGroups []models.DeploymentFileGroup,
	volumes []models.DeploymentVolume,
	containers []models.DeploymentContainerBase,
) error {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO deployments (id, mod_id, mod_source, mod_channel, mod_ver, dir, files_dir, enabled, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		deployment.Id,
		deployment.ModuleId,
		deployment.ModuleSource,
		deployment.ModuleChannel,
		deployment.ModuleVersion,
		deployment.DirName,
		deployment.FilesDirName,
		deployment.Enabled,
		deployment.Created,
		deployment.Updated,
	)
	if err != nil {
		return err
	}
	err = h.createDeploymentResourcesAndConfigs(
		ctx,
		tx,
		deployment.Id,
		hostResources,
		secrets,
		userConfigs,
		globalConfigs,
		files,
		fileGroups,
		volumes,
		containers,
	)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) createDeploymentResourcesAndConfigs(
	ctx context.Context,
	tx *sql.Tx,
	deploymentId string,
	hostResources []models.DeploymentHostResource,
	secrets []models.DeploymentSecret,
	userConfigs []models.DeploymentUserConfig,
	globalConfigs []models.DeploymentGlobalConfig,
	files []models.DeploymentFile,
	fileGroups []models.DeploymentFileGroup,
	volumes []models.DeploymentVolume,
	containers []models.DeploymentContainerBase,
) (err error) {
	for _, hostResource := range hostResources {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO dep_host_resources (dep_id, ref, res_id) VALUES (?, ?, ?)",
			deploymentId,
			hostResource.Reference,
			hostResource.Id,
		)
		if err != nil {
			return
		}
	}
	for _, secret := range secrets {
		for _, item := range secret.Items {
			_, err = tx.ExecContext(
				ctx,
				"INSERT INTO dep_secrets (dep_id, ref, sec_id, item, as_mount, as_env) VALUES (?, ?, ?, ?, ?, ?)",
				deploymentId,
				secret.Reference,
				secret.Id,
				item.Name,
				item.AsMount,
				item.AsEnv,
			)
			if err != nil {
				return
			}
		}
	}
	for _, config := range userConfigs {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO dep_configs (id, dep_id, ref, data_type, is_list) VALUES (?, ?, ?, ?, ?)",
			config.Id,
			config.DeploymentId,
			config.Reference,
			config.DataType,
			config.IsSlice,
		)
		if err != nil {
			return
		}
		err = createConfigValues(ctx, tx, "dep_config_values", config.Id, config.Value)
		if err != nil {
			return
		}
	}
	for _, globalConfig := range globalConfigs {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO dep_global_configs (dep_id, ref, c_id) VALUES (?, ?, ?)",
			deploymentId,
			globalConfig.Reference,
			globalConfig.Id,
		)
		if err != nil {
			return
		}
	}
	for _, file := range files {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO dep_files (dep_id, ref, data) VALUES (?, ?, ?)",
			deploymentId,
			file.Reference,
			file.Data,
		)
		if err != nil {
			return
		}
	}
	for _, fileGroup := range fileGroups {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO dep_file_groups (id, dep_id, ref) VALUES (?, ?, ?)",
			fileGroup.Id,
			deploymentId,
			fileGroup.Reference,
		)
		if err != nil {
			return
		}
		err = createFileGroupFiles(ctx, tx, fileGroup.Id, fileGroup.Files)
		if err != nil {
			return
		}
	}
	for _, volume := range volumes {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO dep_volumes (dep_id, ref, name) VALUES (?, ?, ?)",
			deploymentId,
			volume.Reference,
			volume.Name,
		)
		if err != nil {
			return
		}
	}
	for _, container := range containers {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO dep_containers (dep_id, name, srv_ref, alias) VALUES (?, ?, ?, ?)",
			deploymentId,
			container.Name,
			container.Reference,
			container.Alias,
		)
		if err != nil {
			return
		}
	}
	return
}

func createFileGroupFiles(ctx context.Context, tx *sql.Tx, groupId string, files []models.DeploymentFileGroupFile) (err error) {
	for _, file := range files {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO dep_file_group_files (g_id, path, format, data) VALUES (?, ?, ?, ?)",
			groupId,
			file.Path,
			file.Format,
			file.Data,
		)
		if err != nil {
			return
		}
	}
	return
}

func genModulesFilter(filter models.ModulesFilter) (string, []any) {
	var fc []string
	var val []any
	if len(filter.Ids) > 0 {
		ids := helper_slices.RemoveDuplicates(filter.Ids)
		fc = append(fc, "id IN ("+genQuestionMarks(len(ids))+")")
		for _, id := range ids {
			val = append(val, id)
		}
	}
	if filter.Source != "" {
		fc = append(fc, "source = ?")
		val = append(val, filter.Source)
	}
	if filter.Channel != "" {
		fc = append(fc, "channel = ?")
		val = append(val, filter.Channel)
	}
	if len(fc) > 0 {
		return " WHERE " + strings.Join(fc, " AND "), val
	}
	return "", nil
}

func genQuestionMarks(numCol int) string {
	if numCol <= 0 {
		return ""
	}
	if numCol >= 2 {
		return strings.Repeat("?, ", numCol-1) + "?"
	}
	return "?"
}
