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
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"mqtt-engine-adapter/api"
	"mqtt-engine-adapter/data"
	"os"
)

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	adaptData := data.GetDataInstance()
	adaptData.Subscriptions = make(map[string]*data.Store, 0)
	adaptData.MqttClients = make(map[string]*data.MqttConn, 0)
	adaptData.DbConns = make(map[string]*data.TDEngine, 0)
	api.Start()
}
