package old_impl

import (
	"database/sql"
	"mgw-module-manager-migration/pkg/old_impl/clients/cew_client"
	"mgw-module-manager-migration/pkg/old_impl/handlers/dep_hdl"
	"mgw-module-manager-migration/pkg/old_impl/handlers/mod_hdl"
	"mgw-module-manager-migration/pkg/old_impl/handlers/modfile_hdl"
	"mgw-module-manager-migration/pkg/old_impl/handlers/storage_hdl"
	"mgw-module-manager-migration/pkg/old_impl/libs/modfile_lib/modfile"
	"mgw-module-manager-migration/pkg/old_impl/libs/modfile_lib/v1/v1dec"
	"mgw-module-manager-migration/pkg/old_impl/libs/modfile_lib/v1/v1gen"
	"mgw-module-manager-migration/pkg/old_impl/util"
	"mgw-module-manager-migration/pkg/old_impl/util/naming_hdl"
	"net/http"
	"time"
)

type Config struct {
	ModHandlerWorkdirPath string
	DepHandlerWorkdirPath string
	ConfigDefsPath        string
	ManagerIDPath         string
	CoreID                string
	DatabaseTimeout       time.Duration
	CewBaseUrl            string
	HmBaseUrl             string
	HttpTimeout           time.Duration
}

type Service struct {
	ManagerId          string
	ModulesHandler     *mod_hdl.Handler
	DeploymentsHandler *dep_hdl.Handler
}

func New(config Config, db *sql.DB) (*Service, error) {
	managerId, err := util.GetManagerID(config.ManagerIDPath, "")
	if err != nil {
		return nil, err
	}

	naming_hdl.Init(config.CoreID, "mgw")

	storageHandler := storage_hdl.New(db)

	mfDecoders := make(modfile.Decoders)
	mfDecoders.Add(v1dec.GetDecoder)
	mfGenerators := make(modfile.Generators)
	mfGenerators.Add(v1gen.GetGenerator)

	modFileHandler := modfile_hdl.New(mfDecoders, mfGenerators)

	modHandler := mod_hdl.New(
		storageHandler,
		modFileHandler,
		config.DatabaseTimeout,
		config.ModHandlerWorkdirPath,
	)
	err = modHandler.Init(0770)
	if err != nil {
		return nil, err
	}

	cewClient := cew_client.New(http.DefaultClient, config.CewBaseUrl)

	depHandler := dep_hdl.New(
		storageHandler,
		cewClient,
		config.DatabaseTimeout,
		config.HttpTimeout,
		config.DepHandlerWorkdirPath,
		managerId,
	)
	err = depHandler.InitWorkspace(0770)
	if err != nil {
		return nil, err
	}

	return &Service{
		ManagerId:          managerId,
		ModulesHandler:     modHandler,
		DeploymentsHandler: depHandler,
	}, nil
}
