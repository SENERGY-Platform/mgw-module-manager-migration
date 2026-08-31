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
	"mgw-module-manager-migration/pkg/components/old_impl/clients/base_client"
	"mgw-module-manager-migration/pkg/components/old_impl/libs/cew_lib"
	"net/http"
)

type Client struct {
	baseClient *base_client.Client
	baseUrl    string
}

func New(httpClient base_client.HTTPClient, baseUrl string) *Client {
	return &Client{
		baseClient: base_client.New(httpClient, customError, cew_lib.HeaderRequestID),
		baseUrl:    baseUrl,
	}
}

func customError(code int, err error) error {
	switch code {
	case http.StatusInternalServerError:
		err = cew_lib.NewInternalError(err)
	case http.StatusNotFound:
		err = cew_lib.NewNotFoundError(err)
	case http.StatusBadRequest:
		err = cew_lib.NewInvalidInputError(err)
	}
	return err
}
