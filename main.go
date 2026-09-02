package main

import (
	"context"
	"fmt"
	"mgw-module-manager-migration/pkg/components/helper/mysql"
	"mgw-module-manager-migration/pkg/components/helper/os_signal"
	"mgw-module-manager-migration/pkg/components/migration"
	"mgw-module-manager-migration/pkg/components/new_impl"
	"mgw-module-manager-migration/pkg/components/old_impl"
	"mgw-module-manager-migration/pkg/configuration"
	"os"
	"syscall"
	"time"
)

func main() {
	ec := 0
	defer func() {
		os.Exit(ec)
	}()

	configuration.ParseFlags()
	config, err := configuration.New(configuration.ConfPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load configuration: %s", err)
		ec = 1
		return
	}

	db, err := mysql.NewSQLDatabase(mysql.Config{
		Address:            config.Database.Address,
		Database:           config.Database.Database,
		User:               config.Database.User,
		Password:           config.Database.Password.Value(),
		Timeout:            time.Duration(config.Database.Timeout),
		MaxOpenConnections: config.Database.MaxOpenConnections,
		MaxIdleConnections: config.Database.MaxIdleConnections,
		ConnMaxLifetime:    time.Duration(config.Database.ConnectionMaxLifetime),
	})
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		ec = 1
		return
	}
	defer db.Close()

	oldService, err := old_impl.New(
		old_impl.Config{
			ModHandlerWorkdirPath: config.OldImpl.ModHandlerWorkdirPath,
			DepHandlerWorkdirPath: config.OldImpl.DepHandlerWorkdirPath,
			ManagerIDPath:         config.OldImpl.ManagerIDPath,
			CoreID:                config.CoreId,
			CewBaseUrl:            config.OldImpl.CewBaseUrl,
			HttpTimeout:           time.Duration(config.OldImpl.HttpTimeout),
		},
		db,
		config.ManagerId,
	)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "init old service: %s\n", err)
		ec = 1
		return
	}

	newService := new_impl.New(
		new_impl.Config{
			ManagerIdPath: config.NewImpl.ManagerIDPath,
			ManagerId:     config.ManagerId,
		},
		db,
	)

	ctx, cf := context.WithCancel(context.Background())

	go func() {
		sig := os_signal.Wait(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		if sig != nil {
			_, _ = fmt.Fprintln(os.Stdout, "caught os signal", sig.String())
		}
		cf()
	}()

	err = migration.Run(ctx, oldService, newService, config.NewImpl.RepositorySource, config.NewImpl.RepositoryChannel)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "migration failed: %s\n", err)
		ec = 1
		return
	}
}
