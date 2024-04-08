/*******************************************************************************
 * Copyright 2024 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package config

import (
	"encoding/json"
	"fmt"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
)

const (
	HeaderValueJson string = "application/json"
)

type MongoConfig struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Collection string `json:"collection,omitempty"`
	DbName     string `json:"dbName,omitempty"`
}

type EndpointInfo struct {
	Certificate string      `json:"certificate,omitempty"`
	Key         string      `json:"key,omitempty"`
	Service     ServiceInfo `json:"service,omitempty"`
}

// ServiceInfo describes a service endpoint that the deployed service is a client of. HTTP or TCP for example.
type ServiceInfo struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

// Uri constructs a string from the populated elements of the ServiceInfo
func (s ServiceInfo) Uri() string {
	return fmt.Sprintf("%s://%s:%v", s.Protocol, s.Host, s.Port)
}

type ApplicationConfig struct {
	Endpoint EndpointInfo       `json:"endpoint,omitempty"`
	Mongo    MongoConfig        `json:"mongo,omitempty"`
	NextHop  ServiceInfo        `json:"nextHop,omitempty"`
	Sdk      config.SdkInfo     `json:"sdk,omitempty"`
	Logging  config.LoggingInfo `json:"logging,omitempty"`
}

func (a ApplicationConfig) AsString() string {
	b, _ := json.Marshal(a)
	return string(b)
}
