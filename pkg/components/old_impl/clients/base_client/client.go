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

package base_client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient  HTTPClient
	cErrorFunc  func(code int, err error) error
	headerReqID string
}

func New(httpClient HTTPClient, cErrorFunc func(code int, err error) error, headerReqID string) *Client {
	return &Client{
		httpClient:  httpClient,
		cErrorFunc:  cErrorFunc,
		headerReqID: headerReqID,
	}
}

func (c *Client) ExecRequestJSON(req *http.Request, v any) error {
	body, err := c.execRequest(req)
	if err != nil {
		return err
	}
	defer body.Close()
	err = readJSON(body, v)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ExecRequestString(req *http.Request) (string, error) {
	body, err := c.execRequest(req)
	if err != nil {
		return "", err
	}
	defer body.Close()
	return readString(body)
}

func (c *Client) ExecRequestVoid(req *http.Request) error {
	body, err := c.execRequest(req)
	if err != nil {
		return err
	}
	defer body.Close()
	_ = readVoid(body)
	return nil
}

func (c *Client) execRequest(req *http.Request) (io.ReadCloser, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		errMsg, err := readString(resp.Body)
		if err != nil || errMsg == "" {
			errMsg = resp.Status
		}
		return nil, c.getError(resp.StatusCode, resp.Header.Get(c.headerReqID), errMsg)
	}
	return resp.Body, nil
}

func (c *Client) getError(sc int, rID, msg string) error {
	var err error
	err = newResponseError(sc, rID, errors.New(msg))
	if sc < 500 {
		err = newClientError(err)
	}
	if sc >= 500 {
		err = newServerError(err)
	}
	return c.cErrorFunc(sc, err)
}

func readVoid(rc io.ReadCloser) error {
	_, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	return nil
}

func readString(rc io.ReadCloser) (string, error) {
	b, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func readJSON(rc io.ReadCloser, v any) error {
	jd := json.NewDecoder(rc)
	err := jd.Decode(v)
	if err != nil {
		_ = readVoid(rc)
		return err
	}
	return nil
}
