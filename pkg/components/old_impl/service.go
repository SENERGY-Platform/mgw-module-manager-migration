package old_impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mgw-module-manager-migration/pkg/components/old_impl/clients/cew_client"
	"mgw-module-manager-migration/pkg/components/old_impl/handlers/dep_hdl"
	"mgw-module-manager-migration/pkg/components/old_impl/handlers/mod_hdl"
	"mgw-module-manager-migration/pkg/components/old_impl/handlers/modfile_hdl"
	"mgw-module-manager-migration/pkg/components/old_impl/handlers/storage_hdl"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/modfile"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/v1/v1dec"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/modfile_lib/v1/v1gen"
	"mgw-module-manager-migration/pkg/components/old_impl/model"
	"mgw-module-manager-migration/pkg/components/old_impl/model/pkg_model"
	"mgw-module-manager-migration/pkg/components/old_impl/util"
	"mgw-module-manager-migration/pkg/components/old_impl/util/naming_hdl"
	"net"
	"net/http"
	"os"
	"slices"
	"time"
)

type Config struct {
	ModHandlerWorkdirPath string
	DepHandlerWorkdirPath string
	ManagerIDPath         string
	CoreID                string
	CewBaseUrl            string
	HttpTimeout           time.Duration
}

type Service struct {
	managerId          string
	managerIDPath      string
	modulesHandler     *mod_hdl.Handler
	deploymentsHandler *dep_hdl.Handler
	storageHandler     *storage_hdl.Handler
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
		config.ModHandlerWorkdirPath,
	)
	err = modHandler.Init(0770)
	if err != nil {
		return nil, err
	}

	cewClient := cew_client.New(newHttpClient(config.HttpTimeout), config.CewBaseUrl)

	depHandler := dep_hdl.New(
		storageHandler,
		cewClient,
		config.DepHandlerWorkdirPath,
		managerId,
	)
	err = depHandler.InitWorkspace(0770)
	if err != nil {
		return nil, err
	}

	return &Service{
		managerId:          managerId,
		managerIDPath:      config.ManagerIDPath,
		modulesHandler:     modHandler,
		deploymentsHandler: depHandler,
		storageHandler:     storageHandler,
	}, nil
}

func (s *Service) GetManagerId() string {
	return s.managerId
}

func (s *Service) ReadManagerIdFile() ([]byte, error) {
	file, err := os.Open(s.managerIDPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func (s *Service) GetModules(ctx context.Context) (map[string]ModuleAndDeployment, error) {
	modules, err := s.modulesHandler.List(ctx)
	if err != nil {
		return nil, err
	}
	deployments, err := s.deploymentsHandler.List(ctx)
	if err != nil {
		return nil, err
	}
	advertisements, err := s.storageHandler.ListDepAdv(ctx, pkg_model.DepAdvFilter{})
	if err != nil {
		return nil, err
	}
	depMap, err := getDepMap(deployments)
	if err != nil {
		return nil, err
	}
	advMap, err := getAdvMap(advertisements)
	if err != nil {
		return nil, err
	}
	result := make(map[string]ModuleAndDeployment)
	for id, module := range modules {
		result[id] = ModuleAndDeployment{
			Module:            module.Module,
			Dir:               module.Dir,
			Deployment:        depMap[id],
			DepAdvertisements: advMap[id],
		}
	}
	return result, nil
}

func getDepMap(deployments map[string]model.Deployment) (map[string]model.Deployment, error) {
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
	return depMap, nil
}

func getAdvMap(advertisements map[string]pkg_model.DepAdvertisement) (map[string][]pkg_model.DepAdvertisement, error) {
	mulDep := make(map[string][]string)
	for _, advertisement := range advertisements {
		depIds := mulDep[advertisement.ModuleID]
		if !slices.Contains(depIds, advertisement.DepID) {
			depIds = append(depIds, advertisement.DepID)
		}
		mulDep[advertisement.ModuleID] = depIds
	}
	depCount := 0
	for _, depIds := range mulDep {
		if len(depIds) > 1 {
			depCount++
		}
	}
	if depCount > 0 {
		return nil, errors.New(fmt.Sprintf("multiple deployments: %v", mulDep))
	}
	advMap := make(map[string][]pkg_model.DepAdvertisement)
	for _, advertisement := range advertisements {
		advMap[advertisement.ModuleID] = append(advMap[advertisement.ModuleID], advertisement)
	}
	return advMap, nil
}

type ModuleAndDeployment struct {
	model.Module
	Dir               string                       `json:"dir"`
	Deployment        model.Deployment             `json:"deployment"`
	DepAdvertisements []pkg_model.DepAdvertisement `json:"dep_advertisements"`
}

func newHttpClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}
