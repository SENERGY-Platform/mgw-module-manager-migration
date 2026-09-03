package migration

import (
	"context"
	"fmt"
	"mgw-module-manager-migration/pkg/components/new_impl"
	"mgw-module-manager-migration/pkg/components/new_impl/components/helper/configs"
	new_model "mgw-module-manager-migration/pkg/components/new_impl/models"
	"mgw-module-manager-migration/pkg/components/old_impl"
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
			dep, err := getNewDeployment(oldModule, repoSource, repoChannel)
			if err != nil {
				fmt.Fprintf(os.Stderr, "transform deployment '%s' '%s' failed: %s\n", oldModule.ID, oldModule.Deployment.ID, err)
				continue
			}
			newModule.Deployment = dep
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

func getNewDeployment(oldModule old_impl.ModuleAndDeployment, repoSource, repoChannel string) (new_impl.Deployment, error) {
	newDeployment := new_impl.Deployment{
		DeploymentBase: new_model.DeploymentBase{
			Id:            oldModule.Deployment.ID,
			ModuleId:      oldModule.ID,
			ModuleSource:  repoSource,
			ModuleChannel: repoChannel,
			ModuleVersion: oldModule.Deployment.Module.Version,
			DirName:       oldModule.Deployment.Dir,
			FilesDirName:  oldModule.Deployment.Dir + "_files",
			Enabled:       oldModule.Deployment.Enabled,
			Created:       oldModule.Deployment.Created,
			Updated:       oldModule.Deployment.Updated,
		},
	}
	for ref, id := range oldModule.Deployment.HostResources {
		newDeployment.HostResources = append(newDeployment.HostResources, new_model.DeploymentHostResource{
			Id:           id,
			DeploymentId: oldModule.Deployment.ID,
			Reference:    ref,
		})
	}
	for ref, secret := range oldModule.Deployment.Secrets {
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
			DeploymentId: oldModule.Deployment.ID,
			Reference:    ref,
			Items:        items,
		})
	}
	for ref, config := range oldModule.Deployment.Configs {
		modConfig, ok := oldModule.Configs[ref]
		if !ok {
			fmt.Fprintf(os.Stderr, "config not defined in module '%s' '%s' '%s': %s\n", oldModule.ID, oldModule.Deployment.ID, ref)
			continue
		}
		value, err := configs.GetValue(config.Value, configs.GetDataType(config.DataType), config.IsSlice)
		if err != nil {
			return new_impl.Deployment{}, fmt.Errorf("get config value '%s': %w", ref, err)
		}
		if modConfig.Default != nil {
			defValue, err := configs.GetValue(modConfig.Default, configs.GetDataType(config.DataType), config.IsSlice)
			if err != nil {
				fmt.Fprintf(os.Stderr, "get config default value '%s' '%s' '%s': %s\n", oldModule.ID, oldModule.Deployment.ID, ref, err)
			} else {
				if configs.ValueIsEqual(value, defValue) {
					continue
				}
			}
		}
		newDeployment.UserConfigs = append(newDeployment.UserConfigs, new_model.DeploymentUserConfig{
			Id:           oldModule.Deployment.ID + "_" + ref,
			DeploymentId: oldModule.Deployment.ID,
			Reference:    ref,
			Value:        value,
		})
	}
	for ref, volName := range oldModule.Deployment.Volumes {
		newDeployment.Volumes = append(newDeployment.Volumes, new_model.DeploymentVolume{
			DeploymentId: oldModule.Deployment.ID,
			Reference:    ref,
			Name:         volName,
		})
	}
	for ref, container := range oldModule.Deployment.Containers {
		name := container.ID
		if container.Info != nil {
			name = container.Info.Name
		}
		newDeployment.Containers = append(newDeployment.Containers, new_model.DeploymentContainerBase{
			Name:         name,
			DeploymentId: oldModule.Deployment.ID,
			Reference:    ref,
			Alias:        container.Alias,
		})
	}
	return newDeployment, nil
}
