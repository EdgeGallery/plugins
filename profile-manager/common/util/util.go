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

package util

import (
	"fmt"
	"github.com/libujacob/jsone"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendConfigToNode(hostUrl string, data string, method string) ([]byte, error) {
	payload := strings.NewReader(data)

	req, err := http.NewRequest(method, hostUrl, payload)
	if err != nil {
		log.Error("New Request error.", err)
		return nil, err
	}

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("sending error.", err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("IoT config response error.", err)
		return nil, err
	}

	fmt.Println(string(body))
	return body, nil
}

func GetValuesByKey(cfg *jsone.O, key string) string {
	var result string = ""
	for _, k := range cfg.Keys() {
		if k == key {
			result, _ = cfg.GetString(key)
		}
	}
	return result
}

func GetBrokerHostAndPort(config *jsone.O) (string, int64, error) {
	broker, err := config.GetObject("broker")
	if err != nil {
		log.Errorf(err.Error(), "invalid param while adding south service.")
		return "", 0, err
	}

	brokerHost, _ := broker.GetString("host")
	brokerPort, _ := broker.GetInt64("port")
	return brokerHost, brokerPort, nil
}
