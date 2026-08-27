/*
 * Copyright 2024 InfAI (CC SES)
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
	"mgw-module-manager-migration/old_impl/model"
	"strings"
	"time"

	"mgw-module-manager-migration/old_impl/model/pkg_model"
)

func (h *Handler) ListDepAdv(ctx context.Context, filter pkg_model.DepAdvFilter) (map[string]pkg_model.DepAdvertisement, error) {
	q := "SELECT `id`, `dep_id`, `mod_id`, `origin`, `ref`, `timestamp` FROM `dep_advertisements`"
	fc, val := genDepAdvFilter(filter)
	if fc != "" {
		q += fc
	}
	rows, err := h.db.QueryContext(ctx, q, val...)
	if err != nil {
		return nil, model.NewInternalError(err)
	}
	defer rows.Close()
	depAdvertisements := make(map[string]pkg_model.DepAdvertisement)
	for rows.Next() {
		var adv pkg_model.DepAdvertisement
		var ts []uint8
		if err = rows.Scan(&adv.ID, &adv.DepID, &adv.ModuleID, &adv.Origin, &adv.Ref, &ts); err != nil {
			return nil, model.NewInternalError(err)
		}
		adv.Timestamp, err = time.Parse(tLayout, string(ts))
		if err != nil {
			return nil, model.NewInternalError(err)
		}
		if adv.Items, err = selectDepAdvItems(ctx, h.db, adv.ID); err != nil {
			return nil, err
		}
		depAdvertisements[adv.ID] = adv
	}
	if err = rows.Err(); err != nil {
		return nil, model.NewInternalError(err)
	}
	return depAdvertisements, nil
}

func (h *Handler) ReadDepAdv(ctx context.Context, dID, ref string) (pkg_model.DepAdvertisement, error) {
	row := h.db.QueryRowContext(ctx, "SELECT `id`, `dep_id`, `mod_id`, `origin`, `ref`, `timestamp` FROM `dep_advertisements` WHERE `dep_id` = ? AND `ref` = ?", dID, ref)
	var adv pkg_model.DepAdvertisement
	var ts []uint8
	err := row.Scan(&adv.ID, &adv.DepID, &adv.ModuleID, &adv.Origin, &adv.Ref, &ts)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg_model.DepAdvertisement{}, model.NewNotFoundError(err)
		}
		return pkg_model.DepAdvertisement{}, model.NewInternalError(err)
	}
	adv.Timestamp, err = time.Parse(tLayout, string(ts))
	if err != nil {
		return pkg_model.DepAdvertisement{}, model.NewInternalError(err)
	}
	if adv.Items, err = selectDepAdvItems(ctx, h.db, adv.ID); err != nil {
		return pkg_model.DepAdvertisement{}, err
	}
	return adv, nil
}

func selectDepAdvItems(ctx context.Context, db *sql.DB, id string) (map[string]string, error) {
	return selectStrMap(ctx, db.QueryContext, "SELECT `key`, `value` FROM `dep_adv_items` WHERE `adv_id` = ?", id)
}

func genDepAdvFilter(filter pkg_model.DepAdvFilter) (string, []any) {
	var fc []string
	var val []any
	if filter.DepID != "" {
		fc = append(fc, "`dep_id` = ?")
		val = append(val, filter.DepID)
	}
	if filter.ModuleID != "" {
		fc = append(fc, "`mod_id` = ?")
		val = append(val, filter.ModuleID)
	}
	if filter.Origin != "" {
		fc = append(fc, "`origin` = ?")
		val = append(val, filter.Origin)
	}
	if filter.Ref != "" {
		fc = append(fc, "`ref` = ?")
		val = append(val, filter.Ref)
	}
	if len(fc) > 0 {
		return " WHERE " + strings.Join(fc, " AND "), val
	}
	return "", nil
}
