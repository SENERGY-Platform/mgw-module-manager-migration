package storage_hdl

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
)

var tables = []string{
	"modules",
	"mod_dependencies",
	"deployments",
	"dependencies",
	"containers",
	"host_resources",
	"secrets",
	"configs",
	"list_configs",
	"aux_deployments",
	"aux_labels",
	"aux_configs",
	"aux_containers",
	"aux_volumes",
	"dep_advertisements",
	"dep_adv_items",
}

func (h *Handler) RenameTables(ctx context.Context) error {
	currentTables, err := h.getTables(ctx)
	if err != nil {
		return err
	}
	done := true
	for _, table := range tables {
		if !slices.Contains(currentTables, "bk_"+table) {
			done = false
			break
		}
	}
	if done {
		return nil
	}
	var todo []string
	for _, table := range tables {
		if !slices.Contains(currentTables, table) {
			continue
		}
		_, _ = fmt.Fprintf(os.Stdout, "rename table '%s' to 'bk_%s'\n", table, table)
		todo = append(todo, fmt.Sprintf("`%s` TO `bk_%s`", table, table))
	}
	if len(todo) == 0 {
		return nil
	}
	_, err = h.db.ExecContext(ctx, "RENAME TABLE "+strings.Join(todo, ", "))
	return err
}

func (h *Handler) getTables(ctx context.Context) ([]string, error) {
	rows, err := h.db.QueryContext(ctx, "SELECT `table_name` FROM `information_schema`.`tables` WHERE `table_schema` = DATABASE()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var currentTables []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}
		currentTables = append(currentTables, name)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return currentTables, nil
}
