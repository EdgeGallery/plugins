/*
 * Copyright 2021 Huawei Technologies Co., Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"github.com/libujacob/jsone"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path/filepath"
	_const "profile-manager/common/const"
	"sigs.k8s.io/yaml"
)

type Config struct {
	Cfg     *jsone.O
	Profile *jsone.O
	Broker  struct {
		Host string
		Port int
	}
}

func NewConfig() *Config {
	c := &Config{}
	if c.readConfig() != nil {
		return nil
	}
	if c.readProfile() != nil {
		return nil
	}
	return c
}

func (c *Config) readConfig() error {

	configFilePath := filepath.FromSlash(_const.ConfigPath)
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Error("Reading profile manager file error.", nil)
		return err
	}

	json, err := yaml.YAMLToJSON(configData)
	out := jsone.ParseJsonObject(json)
	c.Cfg = &out

	broker, err := c.Cfg.GetObject("broker")
	if err != nil {
		return err
	}

	brokerHost, err := broker.GetString("host")
	if err != nil {
		return err
	}
	c.Broker.Host = brokerHost
	brokerPort, err := broker.GetInt64("port")
	if err != nil {
		return err
	}
	c.Broker.Port = int(brokerPort)
	return nil
}

func (c *Config) readProfile() error {
	profileFilePath := filepath.FromSlash(_const.ProficlePath)
	configData, err := ioutil.ReadFile(profileFilePath)
	if err != nil {
		log.Error("Reading profile file error.", nil)
		return err
	}

	json, err := yaml.YAMLToJSON(configData)
	out := jsone.ParseJsonObject(json)
	c.Profile = &out
	return nil
}
