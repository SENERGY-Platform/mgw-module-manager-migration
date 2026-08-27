package old_impl

import (
	"database/sql"
	"mgw-module-manager-migration/old_impl/clients/cew_client"
	"mgw-module-manager-migration/old_impl/clients/hm_client"
	"mgw-module-manager-migration/old_impl/handlers/cfg_valid_hdl"
	"mgw-module-manager-migration/old_impl/handlers/cfg_valid_hdl/validators"
	"mgw-module-manager-migration/old_impl/handlers/dep_hdl"
	"mgw-module-manager-migration/old_impl/handlers/mod_hdl"
	"mgw-module-manager-migration/old_impl/handlers/modfile_hdl"
	"mgw-module-manager-migration/old_impl/handlers/storage_hdl"
	"mgw-module-manager-migration/old_impl/libs/modfile_lib/modfile"
	"mgw-module-manager-migration/old_impl/libs/modfile_lib/v1/v1dec"
	"mgw-module-manager-migration/old_impl/libs/modfile_lib/v1/v1gen"
	"mgw-module-manager-migration/old_impl/util"
	"mgw-module-manager-migration/old_impl/util/naming_hdl"
	"net/http"
	"time"
)

type HttpClientConfig struct {
	CewBaseUrl string        `json:"cew_base_url" env_var:"CEW_BASE_URL"`
	CmBaseUrl  string        `json:"cm_base_url" env_var:"CM_BASE_URL"`
	HmBaseUrl  string        `json:"hm_base_url" env_var:"HM_BASE_URL"`
	SmBaseUrl  string        `json:"sm_base_url" env_var:"SM_BASE_URL"`
	Timeout    time.Duration `json:"timeout" env_var:"HTTP_TIMEOUT"`
}

type ModHandlerConfig struct {
	WorkdirPath string `json:"workdir_path" env_var:"MH_WORKDIR_PATH"`
}

type DepHandlerConfig struct {
	WorkdirPath string `json:"workdir_path" env_var:"DH_WORKDIR_PATH"`
	HostDepPath string `json:"host_dep_path" env_var:"DH_HOST_DEP_PATH"`
	HostSecPath string `json:"host_sec_path" env_var:"DH_HOST_SEC_PATH"`
	ModuleNet   string `json:"module_net" env_var:"DH_MODULE_NET"`
}

type Config struct {
	ModHandler      ModHandlerConfig `json:"module_handler"`
	DepHandler      DepHandlerConfig `json:"deployment_handler"`
	ConfigDefsPath  string           `json:"config_defs_path"`
	HttpClient      HttpClientConfig `json:"http_client"`
	ManagerIDPath   string           `json:"manager_id_path"`
	CoreID          string           `json:"core_id"`
	DatabaseTimeout time.Duration    `json:"database_timeout"`
}

type Service struct {
	ManagerId          string
	ModulesHandler     *mod_hdl.Handler
	DeploymentsHandler *dep_hdl.Handler
}

var inputValidators = map[string]cfg_valid_hdl.Validator{
	"regex":            validators.Regex,
	"number_compare":   validators.NumberCompare,
	"text_len_compare": validators.TextLenCompare,
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
		config.HttpClient.Timeout,
		config.ModHandler.WorkdirPath,
	)
	err = modHandler.Init(0770)
	if err != nil {
		return nil, err
	}

	cfgDefs, err := cfg_valid_hdl.LoadDefs(config.ConfigDefsPath)
	if err != nil {
		return nil, err
	}

	cfgValidHandler, err := cfg_valid_hdl.New(cfgDefs, inputValidators)
	if err != nil {
		return nil, err
	}

	hmClient := hm_client.New(http.DefaultClient, config.HttpClient.HmBaseUrl)

	cewClient := cew_client.New(http.DefaultClient, config.HttpClient.CewBaseUrl)

	depHandler := dep_hdl.New(
		storageHandler,
		cfgValidHandler,
		cewClient,
		hmClient,
		config.DatabaseTimeout,
		config.HttpClient.Timeout,
		config.DepHandler.WorkdirPath,
		config.DepHandler.HostDepPath,
		config.DepHandler.HostSecPath,
		managerId,
		config.DepHandler.ModuleNet,
		config.CoreID,
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
