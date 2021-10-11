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
	"profile-manager/nodes"
)

func main() {
	c := NewConfig()
	err := healthCheckProfile(c.Cfg)
	if err != nil {
		return
	}
	err = applyProfile(c)
	if err != nil {
		return
	}
}

func healthCheckProfile(cfg *jsone.O) error {
	configs, err := cfg.GetObject("config")
	if err != nil {
		return err
	}

	for key, _ := range configs {
		nodeObj, _ := configs.GetObject(key)
		if nodeObj.Has("readiness") {
			return nodes.GetProfileNode(key).ReadinessCheck(&nodeObj)
		}
	}
	return nil
}

func applyProfile(configs *Config) error {
	profile := configs.Profile
	profileList, err := profile.GetArray("profile")
	if err != nil {
		return err
	}
	cfg := configs.Cfg
	config, err := cfg.GetObject("config")
	if err != nil {
		return err
	}

	for _, node := range *profileList {
		for _, key := range (node.(jsone.O)).Keys() {
			profileNode := nodes.GetProfileNode(key)

			if profileNode == nil {
				continue
			}

			nodeObj, err := (node.(jsone.O)).GetObject(key)
			if err != nil {
				return err
			}

			for _, cfgType := range nodeObj.Keys() {
				prfCfgList, err := nodeObj.GetArray(cfgType)
				if err != nil {
					return err
				}
				cfgObj, err := config.GetObject(key)
				if err != nil {
					return err
				}
				result := profileNode.ApplyConfig(cfgType, prfCfgList, &cfgObj, cfg)
				if result != nil {
					return result
				}
			}
		}
	}
	return nil
}
