/*
 * Copyright 2023 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dep_hdl

import (
	"context"
	"errors"
	"fmt"
	"mgw-module-manager-migration/pkg/old_impl/handlers/context_hdl"
	"mgw-module-manager-migration/pkg/old_impl/libs/hm_lib"
	"mgw-module-manager-migration/pkg/old_impl/libs/module_lib"
	"mgw-module-manager-migration/pkg/old_impl/model"
	"mgw-module-manager-migration/pkg/old_impl/util/parser"
)

func (h *Handler) getDepAssets(ctx context.Context, mod *module_lib.Module, dID string, depInput model.DepInput) (map[string]hm_lib.HostResource, map[string]secret, map[string]model.DepConfig, error) {
	hostResources, err := h.getHostResources(ctx, mod.HostResources, depInput.HostResources)
	if err != nil {
		return nil, nil, nil, err
	}
	userConfigs, err := h.getUserConfigs(mod.Configs, depInput.Configs)
	if err != nil {
		return nil, nil, nil, err
	}
	secrets, err := h.getSecrets(ctx, mod, dID, depInput.Secrets)
	if err != nil {
		return nil, nil, nil, err
	}
	return hostResources, secrets, userConfigs, nil
}

func (h *Handler) newDepAssets(hostResources map[string]hm_lib.HostResource, secrets map[string]secret, userConfigs map[string]model.DepConfig) model.DepAssets {
	depAssets := model.DepAssets{
		HostResources: make(map[string]string),
		Secrets:       make(map[string]model.DepSecret),
		Configs:       userConfigs,
	}
	for ref, hostResource := range hostResources {
		depAssets.HostResources[ref] = hostResource.ID
	}
	for ref, sec := range secrets {
		var variants []model.DepSecretVariant
		for _, variant := range sec.Variants {
			variants = append(variants, model.DepSecretVariant{
				Item:    variant.Item,
				AsMount: variant.Path != "",
				AsEnv:   variant.AsEnv,
			})
		}
		depAssets.Secrets[ref] = model.DepSecret{
			ID:       sec.ID,
			Variants: variants,
		}
	}
	return depAssets
}

func (h *Handler) getUserConfigs(mConfigs module_lib.Configs, userInput map[string]any) (map[string]model.DepConfig, error) {
	userConfigs := make(map[string]model.DepConfig)
	for ref, mConfig := range mConfigs {
		val, ok := userInput[ref]
		if !ok || val == nil {
			if mConfig.Default == nil && mConfig.Required {
				return nil, model.NewInvalidInputError(fmt.Errorf("config '%s' requried", ref))
			}
		} else {
			var v any
			var err error
			if mConfig.IsSlice {
				v, err = parser.AnyToDataTypeSlice(val, mConfig.DataType)
			} else {
				v, err = parser.AnyToDataType(val, mConfig.DataType)
			}
			if err != nil {
				return nil, model.NewInvalidInputError(fmt.Errorf("parsing user input '%s' failed: %s", ref, err))
			}
			if err = h.cfgVltHandler.ValidateValue(mConfig.Type, mConfig.TypeOpt, v, mConfig.IsSlice, mConfig.DataType); err != nil {
				return nil, model.NewInvalidInputError(err)
			}
			if mConfig.Options != nil && !mConfig.OptExt {
				if err = h.cfgVltHandler.ValidateValInOpt(mConfig.Options, v, mConfig.IsSlice, mConfig.DataType); err != nil {
					return nil, model.NewInvalidInputError(err)
				}
			}
			userConfigs[ref] = model.DepConfig{
				Value:    v,
				DataType: mConfig.DataType,
				IsSlice:  mConfig.IsSlice,
			}
		}
	}
	return userConfigs, nil
}

func (h *Handler) getHostResources(ctx context.Context, mHostRes map[string]module_lib.HostResource, userInput map[string]string) (map[string]hm_lib.HostResource, error) {
	usrHostRes, missing, err := getUserHostRes(userInput, mHostRes)
	if err != nil {
		return nil, model.NewInvalidInputError(err)
	}
	if len(missing) > 0 {
		return nil, model.NewInternalError(errors.New("host resource discovery not implemented"))
	}
	ch := context_hdl.New()
	defer ch.CancelAll()
	hostRes := make(map[string]hm_lib.HostResource)
	for ref, id := range usrHostRes {
		res, err := h.hmClient.GetHostResource(ch.Add(context.WithTimeout(ctx, h.httpTimeout)), id)
		if err != nil {
			return nil, model.NewInternalError(err)
		}
		hostRes[ref] = res
	}
	return hostRes, nil
}

func (h *Handler) getSecrets(ctx context.Context, mod *module_lib.Module, dID string, userInput map[string]string) (map[string]secret, error) {
	usrSecrets, missing, err := getUserSecrets(userInput, mod.Secrets)
	if err != nil {
		return nil, model.NewInvalidInputError(err)
	}
	if len(missing) > 0 {
		return nil, model.NewInternalError(errors.New("secret discovery not implemented"))
	}
	ch := context_hdl.New()
	defer ch.CancelAll()
	secrets := make(map[string]secret)
	variants := make(map[string]secretVariant)
	for ref, sID := range usrSecrets {
		sec, ok := secrets[ref]
		if !ok {
			sec.ID = sID
			sec.Variants = make(map[string]secretVariant)
		}
		for _, service := range mod.Services {
			for _, target := range service.SecretMounts {
				if target.Ref == ref {
					vID := newSecretVariantID(sID, target.Item)
					variant, ok := variants[vID]
					if variant.Path == "" {
						if !ok {
							variant.Item = target.Item
						}
						variants[vID] = variant
					}
					sec.Variants[vID] = variant
				}
			}
			for _, target := range service.SecretVars {
				if target.Ref == ref {
					vID := newSecretVariantID(sID, target.Item)
					variant, ok := variants[vID]
					if !variant.AsEnv {
						if !ok {
							variant.Item = target.Item
						}
						variant.AsEnv = true
						variants[vID] = variant
					}
					sec.Variants[vID] = variant
				}
			}
		}
		secrets[ref] = sec
	}
	return secrets, nil
}

func newSecretVariantID(id string, item *string) string {
	if item != nil {
		return id + *item
	}
	return id
}

func getUserHostRes(userInput map[string]string, mHostRes map[string]module_lib.HostResource) (map[string]string, []string, error) {
	usrHostRes := make(map[string]string)
	var missing []string
	for ref, mHR := range mHostRes {
		id, ok := userInput[ref]
		if ok {
			usrHostRes[ref] = id
		} else {
			if mHR.Required {
				if len(mHR.Tags) > 0 {
					missing = append(missing, ref)
				} else {
					return nil, nil, fmt.Errorf("host resource '%s' required", ref)
				}
			}
		}
	}
	return usrHostRes, missing, nil
}

func getUserSecrets(userInput map[string]string, mSecrets map[string]module_lib.Secret) (map[string]string, []string, error) {
	usrSecrets := make(map[string]string)
	var missing []string
	for ref, mS := range mSecrets {
		id, ok := userInput[ref]
		if ok {
			usrSecrets[ref] = id
		} else {
			if mS.Required {
				if len(mS.Tags) > 0 {
					missing = append(missing, ref)
				} else {
					return nil, nil, fmt.Errorf("secret '%s' required", ref)
				}
			}
		}
	}
	return usrSecrets, missing, nil
}
