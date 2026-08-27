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

package hm_client

import (
	"context"
	"mgw-module-manager-migration/old_impl/libs/hm_lib"
	"net/http"
	"net/url"
)

func (c *Client) GetHostResource(ctx context.Context, id string) (hm_lib.HostResource, error) {
	u, err := url.JoinPath(c.baseUrl, hm_lib.HostResourcesPath, id)
	if err != nil {
		return hm_lib.HostResource{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return hm_lib.HostResource{}, err
	}
	var resource hm_lib.HostResource
	err = c.baseClient.ExecRequestJSON(req, &resource)
	if err != nil {
		return hm_lib.HostResource{}, err
	}
	return resource, nil
}
