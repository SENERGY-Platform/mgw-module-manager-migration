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

package cew_client

import (
	"context"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/cew_lib"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) GetContainers(ctx context.Context, filter cew_lib.ContainerFilter) ([]cew_lib.Container, error) {
	u, err := url.JoinPath(c.baseUrl, cew_lib.ContainersPath)
	if err != nil {
		return nil, err
	}
	u += genGetContainersQuery(filter)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	var containers []cew_lib.Container
	err = c.baseClient.ExecRequestJSON(req, &containers)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func genGetContainersQuery(filter cew_lib.ContainerFilter) string {
	var q []string
	if filter.Name != "" {
		q = append(q, "name="+filter.Name)
	}
	if filter.State != "" {
		q = append(q, "state="+filter.State)
	}
	if len(filter.Labels) > 0 {
		q = append(q, "labels="+genLabels(filter.Labels, "=", ","))
	}
	if len(q) > 0 {
		return "?" + strings.Join(q, "&")
	}
	return ""
}
