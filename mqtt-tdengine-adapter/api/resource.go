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

package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mqtt-engine-adapter/data"
	"net/http"
)

func HandleReq(c echo.Context) error {
	requestData := data.Subscription{}
	if err := c.Bind(&requestData); err != nil {
		fmt.Println("Bind error")
		return err
	}
	adaptData := data.GetDataInstance()
	brokerKey := fmt.Sprintf("%s:%d", requestData.BrokerAddr.Host, requestData.BrokerAddr.Port)
	mqttConn, done := createOrGetMqtt(adaptData, brokerKey, &requestData)
	if done {
		return c.String(http.StatusNotFound, "create mqtt error")
	}
	//tdKey := fmt.Sprintf("%s:%d", requestData.DbAddr.Host, requestData.DbAddr.Port)
	tdStore, done := createOrGetTdEngine(adaptData, brokerKey, &requestData)
	if done {
		return c.String(http.StatusNotFound, "create store error")
	}
	for _, s := range requestData.Store {
		storeKey := fmt.Sprintf("%s:%s", brokerKey, s.Topic)
		mqttConn.Subscribe(s.Topic)            // add subscriptions
		tdStore.CreateStable(&s)               // create stable if not exists
		adaptData.Subscriptions[storeKey] = &s // add subscription mapping
	}
	return c.JSON(http.StatusCreated, requestData)
}

func createOrGetMqtt(adaptData *data.AdapterData, key string,
	requestData *data.Subscription) (*data.MqttConn, bool) {
	if adaptData.MqttClients[key] != nil {
		return adaptData.MqttClients[key], false
	}
	mqttConn := data.NewMQTTConnection(requestData.BrokerAddr.Host, requestData.BrokerAddr.Port)
	if mqttConn == nil {
		return nil, true
	}
	adaptData.MqttClients[key] = mqttConn
	return mqttConn, false
}

func createOrGetTdEngine(adaptData *data.AdapterData, key string,
	requestData *data.Subscription) (*data.TDEngine, bool) {
	if adaptData.DbConns[key] != nil {
		return adaptData.DbConns[key], false
	}
	tdStore := data.NewStore(requestData)
	if tdStore == nil {
		return nil, true
	}
	adaptData.DbConns[key] = tdStore
	return tdStore, false
}
