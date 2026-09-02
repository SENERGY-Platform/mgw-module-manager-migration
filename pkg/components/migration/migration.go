package migration

import (
	"context"
	"fmt"
	"mgw-module-manager-migration/pkg/components/new_impl"
	"mgw-module-manager-migration/pkg/components/new_impl/components/helper/configs"
	new_model "mgw-module-manager-migration/pkg/components/new_impl/models"
	"mgw-module-manager-migration/pkg/components/old_impl"
	old_model "mgw-module-manager-migration/pkg/components/old_impl/model"
	"os"
)

func Run(ctx context.Context, oldSrv *old_impl.Service, newSrv *new_impl.Service, repoSource, repoChannel string) error {
	err := oldSrv.RenameDatabaseTables(ctx)
	if err != nil {
		return fmt.Errorf("backup old database: %w", err)
	}
	oldModules, err := oldSrv.GetModules(ctx)
	if err != nil {
		return fmt.Errorf("read old database: %w", err)
	}
	var newModules []new_impl.Module
	for _, oldModule := range oldModules {
		newModule := new_impl.Module{
			DatabaseModule: new_model.DatabaseModule{
				Id:      oldModule.ID,
				DirName: oldModule.Dir,
				Source:  repoSource,
				Channel: repoChannel,
				Added:   oldModule.Added,
				Updated: oldModule.Updated,
			},
			IsDeployed: oldModule.IsDeployed,
		}
		if oldModule.IsDeployed {
			newModule.Deployment = getNewDeployment(oldModule.ID, repoSource, repoChannel, oldModule.Deployment)
			for _, advertisement := range oldModule.DepAdvertisements {
				newModule.DepAdvertisements = append(newModule.DepAdvertisements, new_model.DeploymentAdvertisement{
					Id:           advertisement.ID,
					DeploymentId: oldModule.Deployment.ID,
					ModuleId:     oldModule.ID,
					Reference:    advertisement.Ref,
					Timestamp:    advertisement.Timestamp,
					Items:        advertisement.Items,
				})
			}
		}
		newModules = append(newModules, newModule)
	}
	err = newSrv.WriteManagerId(oldSrv.GetManagerId())
	if err != nil {
		return err
	}
	err = newSrv.InitDatabaseTables(ctx)
	if err != nil {
		return fmt.Errorf("init new database: %w", err)
	}
	err = newSrv.WriteModules(ctx, newModules)
	if err != nil {
		return fmt.Errorf("write database: %w", err)
	}
	return nil
}

func getNewDeployment(moduleId, repoSource, repoChannel string, oldDeployment old_model.Deployment) new_impl.Deployment {
	newDeployment := new_impl.Deployment{
		DeploymentBase: new_model.DeploymentBase{
			Id:            oldDeployment.ID,
			ModuleId:      moduleId,
			ModuleSource:  repoSource,
			ModuleChannel: repoChannel,
			ModuleVersion: oldDeployment.Module.Version,
			DirName:       oldDeployment.Dir,
			FilesDirName:  oldDeployment.Dir + "_files", // TODO create folder in new impl service
			Enabled:       oldDeployment.Enabled,
			Created:       oldDeployment.Created,
			Updated:       oldDeployment.Updated,
		},
	}
	for ref, id := range oldDeployment.HostResources {
		newDeployment.HostResources = append(newDeployment.HostResources, new_model.DeploymentHostResource{
			Id:           id,
			DeploymentId: oldDeployment.ID,
			Reference:    ref,
		})
	}
	for ref, secret := range oldDeployment.Secrets {
		var items []new_model.DeploymentSecretItem
		for _, variant := range secret.Variants {
			var name string
			if variant.Item != nil {
				name = *variant.Item
			}
			items = append(items, new_model.DeploymentSecretItem{
				Name:    name,
				AsMount: variant.AsMount,
				AsEnv:   variant.AsEnv,
			})
		}
		newDeployment.Secrets = append(newDeployment.Secrets, new_model.DeploymentSecret{
			Id:           secret.ID,
			DeploymentId: oldDeployment.ID,
			Reference:    ref,
			Items:        items,
		})
	}
	for ref, config := range oldDeployment.Configs {
		value, err := configs.GetValue(config.Value, configs.GetDataType(config.DataType), config.IsSlice)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get config value '%s' '%s' '%s': %s\n", moduleId, oldDeployment.ID, ref, err)
			continue
		}
		newDeployment.UserConfigs = append(newDeployment.UserConfigs, new_model.DeploymentUserConfig{
			Id:           oldDeployment.ID + "_" + ref,
			DeploymentId: oldDeployment.ID,
			Reference:    ref,
			Value:        value,
		})
	}
	for ref, volName := range oldDeployment.Volumes {
		newDeployment.Volumes = append(newDeployment.Volumes, new_model.DeploymentVolume{
			DeploymentId: oldDeployment.ID,
			Reference:    ref,
			Name:         volName,
		})
	}
	for ref, container := range oldDeployment.Containers {
		name := container.ID
		if container.Info != nil {
			name = container.Info.Name
		}
		newDeployment.Containers = append(newDeployment.Containers, new_model.DeploymentContainerBase{
			Name:         name,
			DeploymentId: oldDeployment.ID,
			Reference:    ref,
			Alias:        container.Alias,
		})
	}
	return newDeployment
}
