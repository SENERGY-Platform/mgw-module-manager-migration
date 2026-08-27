/*
 * Copyright 2024 InfAI (CC SES)
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

import "time"

type MDNSEntry struct {
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Subtypes   []string  `json:"subtypes"`
	Domain     string    `json:"domain"`
	Hostname   string    `json:"hostname"`
	Port       int       `json:"port"`
	IPv4Addr   string    `json:"ipv4_addr"`
	TxtRecords []string  `json:"txt_records"`
	Expiry     time.Time `json:"expiry"`
}
