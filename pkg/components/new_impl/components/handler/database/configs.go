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
	"fmt"
	helper_slices "mgw-module-manager-migration/pkg/components/new_impl/components/helper/slices"
	"mgw-module-manager-migration/pkg/components/new_impl/models"
	"mgw-module-manager-migration/pkg/components/new_impl/models/constants"
)

func createConfigValues(ctx context.Context, tx *sql.Tx, tableName string, id string, value models.Value) error {
	if value.IsSlice {
		colName, itfValues := getListInterfaceValsAndCol(value)
		stmt := fmt.Sprintf("INSERT INTO %s (c_id, %s, ord) VALUES (?, ?, ?)", tableName, colName)
		for i, itfValue := range itfValues {
			_, err := tx.ExecContext(ctx, stmt, id, itfValue, i)
			if err != nil {
				return err
			}
		}
	} else {
		colName, itfValue := getInterfaceValAndCol(value)
		_, err := tx.ExecContext(
			ctx,
			fmt.Sprintf("INSERT INTO %s (c_id, %s, ord) VALUES (?, ?, ?)", tableName, colName),
			id,
			itfValue,
			0,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func getInterfaceValAndCol(v models.Value) (colName string, value interface{}) {
	switch v.DataType {
	case constants.ValueDataTypeString:
		colName = "v_string"
		value = v.String
	case constants.ValueDataTypeInt64:
		colName = "v_int"
		value = v.Int64
	case constants.ValueDataTypeFloat64:
		colName = "v_float"
		value = v.Float64
	case constants.ValueDataTypeBool:
		colName = "v_bool"
		value = v.Bool
	}
	return
}

func getListInterfaceValsAndCol(v models.Value) (colName string, values []interface{}) {
	switch v.DataType {
	case constants.ValueDataTypeString:
		colName = "v_string"
		values = helper_slices.ToAny(v.StringSlice)
	case constants.ValueDataTypeInt64:
		colName = "v_int"
		values = helper_slices.ToAny(v.Int64Slice)
	case constants.ValueDataTypeFloat64:
		colName = "v_float"
		values = helper_slices.ToAny(v.Float64Slice)
	case constants.ValueDataTypeBool:
		colName = "v_bool"
		values = helper_slices.ToAny(v.BoolSlice)
	}
	return
}
