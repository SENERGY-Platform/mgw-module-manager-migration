package new_impl

import (
	"context"
	"database/sql"
	"fmt"
	"mgw-module-manager-migration/pkg/components/new_impl/components/handler/database"
	"mgw-module-manager-migration/pkg/components/new_impl/components/handler/database/migrations/db_init"
	"mgw-module-manager-migration/pkg/components/new_impl/components/helper/naming"
	"mgw-module-manager-migration/pkg/components/new_impl/models"
)

type Config struct {
	ManagerIdPath string
	ManagerId     string
}

type Service struct {
	config          Config
	databaseHandler *database.Handler
}

func New(config Config, db *sql.DB) *Service {
	return &Service{
		config:          config,
		databaseHandler: database.New(db),
	}
}

func (s *Service) InitDatabaseTables(ctx context.Context) error {
	return s.databaseHandler.Migrate(ctx, db_init.Migration)
}

func (s *Service) WriteManagerId(id string) error {
	if s.config.ManagerId != "" {
		return nil
	}
	return naming.SetManagerId(s.config.ManagerIdPath, id)
}

func (s *Service) WriteModules(ctx context.Context, modules []Module) error {
	tx, err := s.databaseHandler.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin database transaction: %w", err)
	}
	defer tx.Rollback()
	for _, module := range modules {
		err = s.databaseHandler.CreateModule(ctx, tx, module.DatabaseModule)
		if err != nil {
			return fmt.Errorf("write module '%s' to database: %w", module.Id, err)
		}
		if !module.IsDeployed {
			continue
		}
		err = s.databaseHandler.CreateDeployment(
			ctx,
			tx,
			module.Deployment.DeploymentBase,
			module.Deployment.HostResources,
			module.Deployment.Secrets,
			module.Deployment.UserConfigs,
			module.Deployment.GlobalConfigs,
			module.Deployment.Files,
			module.Deployment.FileGroups,
			module.Deployment.Volumes,
			module.Deployment.Containers,
		)
		if err != nil {
			return fmt.Errorf("write deployment '%s' '%s' to database: %w", module.Id, module.Deployment.Id, err)
		}
		if len(module.DepAdvertisements) > 0 {
			err = s.databaseHandler.WriteDeploymentAdvertisements(
				ctx,
				tx,
				module.Deployment.Id,
				module.DepAdvertisements,
				false,
			)
			if err != nil {
				return fmt.Errorf("write deployment advertisements '%s' '%s' to database: %w", module.Id, module.Deployment.Id, err)
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit database transaction: %w", err)
	}
	return nil
}

type Module struct {
	models.DatabaseModule
	Deployment        Deployment
	DepAdvertisements []models.DeploymentAdvertisement
	IsDeployed        bool
}

type Deployment struct {
	models.DeploymentBase
	HostResources []models.DeploymentHostResource
	Secrets       []models.DeploymentSecret
	UserConfigs   []models.DeploymentUserConfig
	GlobalConfigs []models.DeploymentGlobalConfig
	Files         []models.DeploymentFile
	FileGroups    []models.DeploymentFileGroup
	Volumes       []models.DeploymentVolume
	Containers    []models.DeploymentContainerBase
}
