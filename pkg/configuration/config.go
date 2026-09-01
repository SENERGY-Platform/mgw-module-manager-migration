package configuration

import (
	"time"

	sb_config_hdl "github.com/SENERGY-Platform/go-service-base/config-hdl"
	sb_config_types "github.com/SENERGY-Platform/go-service-base/config-hdl/types"
)

type DatabaseConfig struct {
	Address               string                   `json:"address" env_var:"DATABASE_ADDRESS"`
	Database              string                   `json:"database" env_var:"DATABASE_NAME"`
	User                  string                   `json:"user" env_var:"DATABASE_USER"`
	Password              sb_config_types.Secret   `json:"password" env_var:"DATABASE_PASSWORD"`
	Timeout               sb_config_types.Duration `json:"timeout" env_var:"DATABASE_TIMEOUT"`
	MaxOpenConnections    int                      `json:"max_open_connections" env_var:"DATABASE_MAX_OPEN_CONNECTIONS"`
	MaxIdleConnections    int                      `json:"max_idle_connections" env_var:"DATABASE_MAX_IDLE_CONNECTIONS"`
	ConnectionMaxLifetime sb_config_types.Duration `json:"connection_max_lifetime" env_var:"DATABASE_CONNECTION_MAX_LIFETIME"`
}

type OldImplConfig struct {
	ModHandlerWorkdirPath string                   `json:"mod_handler_workdir_path" env_var:"OLD_MOD_WORKDIR_PATH"`
	DepHandlerWorkdirPath string                   `json:"dep_handler_workdir_path" env_var:"OLD_DEP_WORKDIR_PATH"`
	ManagerIDPath         string                   `json:"manager_id_path" env_var:"OLD_MANAGER_ID_PATH"`
	CewBaseUrl            string                   `json:"cew_base_url"  env_var:"CEW_BASE_URL"`
	HttpTimeout           sb_config_types.Duration `json:"http_timeout" env_var:"HTTP_TIMEOUT"`
}

type NewImplConfig struct {
	ManagerIDPath     string `json:"manager_id_path" env_var:"NEW_MANAGER_ID_PATH"`
	RepositorySource  string `json:"repository_source" env_var:"REPOSITORY_SOURCE"`
	RepositoryChannel string `json:"repository_channel" env_var:"REPOSITORY_CHANNEL"`
}

type Config struct {
	CoreId    string         `json:"core_id" env_var:"CORE_ID"`
	ManagerId string         `json:"manager_id" env_var:"MANAGER_ID"`
	Database  DatabaseConfig `json:"database"`
	OldImpl   OldImplConfig  `json:"old_impl"`
	NewImpl   NewImplConfig  `json:"new_impl"`
}

var defaultConfig = Config{
	Database: DatabaseConfig{
		Database:              "module_manager",
		Timeout:               sb_config_types.Duration(time.Second * 30),
		MaxOpenConnections:    25,
		MaxIdleConnections:    25,
		ConnectionMaxLifetime: sb_config_types.Duration(time.Minute * 5),
	},
	OldImpl: OldImplConfig{
		ModHandlerWorkdirPath: "/opt/module-manager/modules",
		DepHandlerWorkdirPath: "/opt/module-manager/deployments",
		ManagerIDPath:         "/opt/module-manager/data/mid",
		CewBaseUrl:            "http://core-api/ce-wrapper",
		HttpTimeout:           sb_config_types.Duration(time.Second * 30),
	},
	NewImpl: NewImplConfig{
		ManagerIDPath:     "/opt/module-manager/service/mid",
		RepositorySource:  "github.com/SENERGY-Platform/mgw-module-repository",
		RepositoryChannel: "main",
	},
}

func New(path string) (Config, error) {
	cfg := defaultConfig
	err := sb_config_hdl.Load(&cfg, nil, envTypeParser, nil, path)
	return cfg, err
}

var envTypeParser = []sb_config_hdl.EnvTypeParser{
	sb_config_types.SecretEnvTypeParser,
	sb_config_types.DurationEnvTypeParser,
}
