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

package services

import (
	"errors"
	"fmt"
	"io/fs"
	"mgw-module-manager-migration/old_impl/libs/modfile_lib/v1/model"
	"mgw-module-manager-migration/old_impl/libs/modfile_lib/v1/v1gen/generic"
	"mgw-module-manager-migration/old_impl/libs/module_lib"
	"time"
)

func GenServices(mfSs map[string]model.Service) (map[string]*module_lib.Service, error) {
	mSs := make(map[string]*module_lib.Service)
	for ref, mfS := range mfSs {
		mBMs, err := GenBindMounts(mfS.Include)
		if err != nil {
			return nil, fmt.Errorf("service '%s' invalid bind mount: %s", ref, err)
		}
		mTMs, err := GenTmpfsMounts(mfS.Tmpfs)
		if err != nil {
			return nil, fmt.Errorf("service '%s' invalid tmpfsMount: %s", ref, err)
		}
		mHEs, err := GenHttpEndpoints(mfS.HttpEndpoints)
		if err != nil {
			return nil, fmt.Errorf("service '%s' invalid http endpoint: %s", ref, err)
		}
		mPs, err := GenPorts(mfS.Ports)
		if err != nil {
			return nil, fmt.Errorf("service '%s' invalid port mapping: %s", ref, err)
		}
		mSs[ref] = &module_lib.Service{
			Name:              mfS.Name,
			Image:             mfS.Image,
			RunConfig:         GenRunConfig(mfS.RunConfig),
			BindMounts:        mBMs,
			Tmpfs:             mTMs,
			HttpEndpoints:     mHEs,
			RequiredSrv:       generic.GenStringSet(mfS.RequiredServices),
			Ports:             mPs,
			DeviceCGroupRules: mfS.DeviceCGroupRules,
		}
	}
	return mSs, nil
}

func GenAuxServices(mfSs map[string]model.AuxService) (map[string]*module_lib.AuxService, error) {
	mAs := make(map[string]*module_lib.AuxService)
	for ref, mfS := range mfSs {
		mBMs, err := GenBindMounts(mfS.Include)
		if err != nil {
			return nil, fmt.Errorf("aux service '%s' invalid bind mount: %s", ref, err)
		}
		mTMs, err := GenTmpfsMounts(mfS.Tmpfs)
		if err != nil {
			return nil, fmt.Errorf("aux service '%s' invalid tmpfsMount: %s", ref, err)
		}
		mAs[ref] = &module_lib.AuxService{
			Name:       mfS.Name,
			RunConfig:  GenRunConfig(mfS.RunConfig),
			BindMounts: mBMs,
			Tmpfs:      mTMs,
		}
	}
	return mAs, nil
}

func GenRunConfig(mfRC model.RunConfig) module_lib.RunConfig {
	mRC := module_lib.RunConfig{
		MaxRetries:  5,
		RunOnce:     mfRC.RunOnce,
		StopTimeout: 5 * time.Second,
		StopSignal:  mfRC.StopSignal,
		PseudoTTY:   mfRC.PseudoTTY,
	}
	if len(mfRC.Command) > 0 {
		mRC.Command = mfRC.Command
	}
	if mfRC.MaxRetries != nil {
		mRC.MaxRetries = uint(*mfRC.MaxRetries)
	}
	if mfRC.StopTimeout != nil {
		mRC.StopTimeout = time.Duration(*mfRC.StopTimeout)
	}
	return mRC
}

func GenBindMounts(mfBMs []model.BindMount) (map[string]module_lib.BindMount, error) {
	mBMs := make(map[string]module_lib.BindMount)
	for _, mfBM := range mfBMs {
		if v, ok := mBMs[mfBM.MountPoint]; ok {
			if v.Source == mfBM.Source && v.ReadOnly == mfBM.ReadOnly {
				continue
			}
			return nil, fmt.Errorf("duplicate '%s'", mfBM.MountPoint)
		}
		mBMs[mfBM.MountPoint] = module_lib.BindMount{
			Source:   mfBM.Source,
			ReadOnly: mfBM.ReadOnly,
		}
	}
	return mBMs, nil
}

func GenTmpfsMounts(mfTMs []model.TmpfsMount) (map[string]module_lib.TmpfsMount, error) {
	mTMs := make(map[string]module_lib.TmpfsMount)
	for _, mfTM := range mfTMs {
		if v, ok := mTMs[mfTM.MountPoint]; ok {
			if v.Size == uint64(mfTM.Size) && (mfTM.Mode == nil || v.Mode == fs.FileMode(*mfTM.Mode)) {
				continue
			}
			return nil, fmt.Errorf("duplicate '%s'", mfTM.MountPoint)
		}
		mTM := module_lib.TmpfsMount{
			Size: uint64(mfTM.Size),
			Mode: fs.FileMode(504),
		}
		if mfTM.Mode != nil {
			mTM.Mode = fs.FileMode(*mfTM.Mode)
		}
		mTMs[mfTM.MountPoint] = mTM
	}
	return mTMs, nil
}

func GenHttpEndpoints(mfHEs []model.HttpEndpoint) (map[string]module_lib.HttpEndpoint, error) {
	mHEs := make(map[string]module_lib.HttpEndpoint)
	for _, mfHE := range mfHEs {
		if _, ok := mHEs[mfHE.ExtPath]; ok {
			return nil, fmt.Errorf("duplicate '%s'", mfHE.ExtPath)
		}
		mHE := module_lib.HttpEndpoint{
			Name: mfHE.Name,
			Port: mfHE.Port,
			Path: mfHE.Path,
			ProxyConf: module_lib.HttpEndpointProxyConf{
				Headers:   mfHE.ProxyConf.Headers,
				WebSocket: mfHE.ProxyConf.WebSocket,
			},
			StringSub: module_lib.HttpEndpointStrSub{
				ReplaceOnce: mfHE.StringSub.ReplaceOnce,
				MimeTypes:   mfHE.StringSub.MimeTypes,
				Filters:     mfHE.StringSub.Filters,
			},
		}
		if mfHE.ProxyConf.ReadTimeout != nil {
			mHE.ProxyConf.ReadTimeout = time.Duration(*mfHE.ProxyConf.ReadTimeout)
		}
		mHEs[mfHE.ExtPath] = mHE
	}
	return mHEs, nil
}

func GenPorts(mfSPs []model.SrvPort) ([]module_lib.Port, error) {
	var mPs []module_lib.Port
	for _, mfSP := range mfSPs {
		proto := module_lib.TcpPort
		if mfSP.Protocol != nil {
			proto = *mfSP.Protocol
		}
		ep, err := mfSP.Port.Parse()
		if err != nil {
			return nil, err
		}
		var hp []uint
		if mfSP.HostPort != nil {
			hp, err = mfSP.HostPort.Parse()
			if err != nil {
				return nil, err
			}
		}
		lep := len(ep)
		lhp := len(hp)
		if lhp > 0 {
			if lep > lhp {
				return nil, errors.New("range mismatch: ports > host ports")
			}
			if lep > 1 && lep < lhp {
				return nil, errors.New("range mismatch: ports < host ports")
			}
		}
		if lep == 1 {
			mPs = append(mPs, module_lib.Port{
				Name:     mfSP.Name,
				Number:   ep[0],
				Protocol: proto,
				Bindings: hp,
			})
		} else {
			for i, n := range ep {
				mP := module_lib.Port{
					Name:     mfSP.Name,
					Number:   n,
					Protocol: proto,
				}
				if lhp > 0 {
					mP.Bindings = []uint{hp[i]}
				}
				mPs = append(mPs, mP)
			}
		}
	}
	return mPs, nil
}
