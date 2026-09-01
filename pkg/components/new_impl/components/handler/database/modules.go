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
	"mgw-module-manager-migration/pkg/components/new_impl/models"
)

func (h *Handler) CreateModule(ctx context.Context, tx *sql.Tx, mod models.DatabaseModule) error {
	_, err := tx.ExecContext(
		ctx,
		"INSERT INTO modules (id, dir, source, channel, added, updated) VALUES (?, ?, ?, ?, ?, ?);",
		mod.Id,
		mod.DirName,
		mod.Source,
		mod.Channel,
		mod.Added,
		mod.Updated,
	)
	if err != nil {
		return err
	}
	return nil
}
