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

type Config struct {
	CoreId   string         `json:"core_id" env_var:"CORE_ID"`
	Database DatabaseConfig `json:"database"`
}

var defaultConfig = Config{
	Database: DatabaseConfig{
		Database:              "module_manager",
		Timeout:               sb_config_types.Duration(time.Second * 30),
		MaxOpenConnections:    25,
		MaxIdleConnections:    25,
		ConnectionMaxLifetime: sb_config_types.Duration(time.Minute * 5),
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
