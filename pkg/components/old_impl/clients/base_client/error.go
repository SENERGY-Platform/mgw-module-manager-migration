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

type cError struct {
	err error
}

type ResponseError struct {
	cError
	Code      int
	RequestID string
}

type ClientError struct {
	cError
}

type ServerError struct {
	cError
}

func newClientError(err error) *ClientError {
	return &ClientError{
		cError: cError{err: err},
	}
}

func newServerError(err error) *ServerError {
	return &ServerError{
		cError: cError{err: err},
	}
}

func newResponseError(c int, rID string, err error) *ResponseError {
	return &ResponseError{
		cError:    cError{err: err},
		Code:      c,
		RequestID: rID,
	}
}

func (e *cError) Error() string {
	return e.err.Error()
}

func (e *cError) Unwrap() error {
	return e.err
}
