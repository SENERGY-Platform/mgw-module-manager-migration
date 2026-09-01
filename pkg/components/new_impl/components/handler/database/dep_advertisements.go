/*
 * Copyright 2026 InfAI (CC SES)
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
	"mgw-module-manager-migration/pkg/components/new_impl/models"
)

func (h *Handler) WriteDeploymentAdvertisements(
	ctx context.Context,
	tx *sql.Tx,
	deploymentId string,
	advertisements []models.DeploymentAdvertisement,
	incremental bool,
) error {
	if !incremental {
		_, err := tx.ExecContext(
			ctx,
			"DELETE FROM dep_advertisements WHERE dep_id = ?;",
			deploymentId,
		)
		if err != nil {
			return err
		}
	}
	for _, advertisement := range advertisements {
		if incremental {
			_, err := tx.ExecContext(
				ctx,
				"DELETE FROM dep_advertisements WHERE dep_id = ? AND ref = ?;",
				deploymentId,
				advertisement.Reference,
			)
			if err != nil {
				return err
			}
		}
		err := h.insertDeploymentAdvertisement(ctx, tx, deploymentId, advertisement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) insertDeploymentAdvertisement(
	ctx context.Context,
	tx *sql.Tx,
	deploymentId string,
	advertisement models.DeploymentAdvertisement,
) error {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO dep_advertisements (id, dep_id, mod_id, ref, timestamp) VALUES (?, ?, ?, ?, ?)",
		advertisement.Id,
		deploymentId,
		advertisement.ModuleId,
		advertisement.Reference,
		advertisement.Timestamp,
	)
	if err != nil {
		return err
	}
	for key, value := range advertisement.Items {
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO dep_adv_items (dep_adv_id, item_key, item_value) VALUES (?, ?, ?)",
			advertisement.Id,
			key,
			value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
