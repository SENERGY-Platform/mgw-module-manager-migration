package old_impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mgw-module-manager-migration/pkg/old_impl/clients/cew_client"
	"mgw-module-manager-migration/pkg/old_impl/handlers/dep_hdl"
	"mgw-module-manager-migration/pkg/old_impl/handlers/mod_hdl"
	"mgw-module-manager-migration/pkg/old_impl/handlers/modfile_hdl"
	"mgw-module-manager-migration/pkg/old_impl/handlers/storage_hdl"
	"mgw-module-manager-migration/pkg/old_impl/libs/modfile_lib/modfile"
	"mgw-module-manager-migration/pkg/old_impl/libs/modfile_lib/v1/v1dec"
	"mgw-module-manager-migration/pkg/old_impl/libs/modfile_lib/v1/v1gen"
	"mgw-module-manager-migration/pkg/old_impl/model"
	"mgw-module-manager-migration/pkg/old_impl/model/pkg_model"
	"mgw-module-manager-migration/pkg/old_impl/util"
	"mgw-module-manager-migration/pkg/old_impl/util/naming_hdl"
	"net/http"
	"time"
)

type Config struct {
	ModHandlerWorkdirPath string
	DepHandlerWorkdirPath string
	ManagerIDPath         string
	CoreID                string
	DatabaseTimeout       time.Duration
	CewBaseUrl            string
	HttpTimeout           time.Duration
}

type Service struct {
	managerId          string
	modulesHandler     *mod_hdl.Handler
	deploymentsHandler *dep_hdl.Handler
}

func New(config Config, db *sql.DB, managerId string) (*Service, error) {
	managerId, err := util.GetManagerID(config.ManagerIDPath, managerId)
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
		managerId:          managerId,
		modulesHandler:     modHandler,
		deploymentsHandler: depHandler,
	}, nil
}

func (s *Service) GetManagerId() string {
	return s.managerId
}

func (s *Service) GetModules(ctx context.Context) (map[string]pkg_model.ModuleAndDeployment, error) {
	modules, err := s.modulesHandler.List(ctx)
	if err != nil {
		return nil, err
	}
	deployments, err := s.deploymentsHandler.List(ctx)
	if err != nil {
		return nil, err
	}
	depMap := make(map[string]model.Deployment)
	mulDep := make(map[string][]string)
	for id, deployment := range deployments {
		_, ok := depMap[deployment.Module.ID]
		if ok {
			mulDep[deployment.Module.ID] = append(mulDep[deployment.Module.ID], id)
		}
		depMap[deployment.Module.ID] = deployment
	}
	if len(mulDep) > 0 {
		return nil, errors.New(fmt.Sprintf("multiple deployments: %v", mulDep))
	}
	result := make(map[string]pkg_model.ModuleAndDeployment)
	for id, module := range modules {
		result[id] = pkg_model.ModuleAndDeployment{
			Module:     module,
			Deployment: depMap[id],
		}
	}
	return result, nil
}
