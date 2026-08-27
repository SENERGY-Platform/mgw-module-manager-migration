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

package hm_lib

import (
	"context"
	"time"
)

type Api interface {
	GetHostInfo(ctx context.Context) (HostInfo, error)
	GetHostNet(ctx context.Context) (HostNet, error)
	ListHostResources(ctx context.Context, filter HostResourceFilter) ([]HostResource, error)
	GetHostResource(ctx context.Context, rID string) (HostResource, error)
	ListHostApplications(ctx context.Context) ([]HostApplication, error)
	AddHostApplication(ctx context.Context, appResBase HostApplicationBase) (string, error)
	RemoveHostApplication(ctx context.Context, aID string) error
	GetNetItfBlacklist(ctx context.Context) ([]string, error)
	NetItfBlacklistAdd(ctx context.Context, v string) error
	NetItfBlacklistRemove(ctx context.Context, v string) error
	GetNetRngBlacklist(ctx context.Context) ([]string, error)
	NetRngBlacklistAdd(ctx context.Context, v string) error
	NetRngBlacklistRemove(ctx context.Context, v string) error
	MDNSQueryService(ctx context.Context, service, domain string, window time.Duration) ([]MDNSEntry, error)
}
