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

package sogno

import (
	"fmt"
	"github.com/libujacob/jsone"
	log "github.com/sirupsen/logrus"
	_const "profile-manager/common/const"
	"profile-manager/common/util"
	"profile-manager/nodes/nodeif"
	"time"
)

type SognoAdpt struct {
	nodeif.Node
}

func (f *SognoAdpt) ReadinessCheck(node *jsone.O) error {

	hostUrls, err := node.GetString("hostUrl")
	readiness, err := node.GetObject("readiness")
	httpGet, err := readiness.GetString("httpGet")
	if hostUrls == "" || httpGet == "" {
		log.Error(err, "url/path is invalid")
		return err
	}
	count := 0
	for {
		count++
		_, err := util.SendConfigToNode(hostUrls+httpGet, "", "GET")
		if err == nil {
			log.Info("Peer is up, readiness check is completed.")
			break
		}
		log.Infof("Waiting for peer status up(%v seconds)", count)
		if count == _const.MaxNumberRetry {
			log.Info("Server is not responding for %v seconds, readiness check stopped.", count)
			err = fmt.Errorf("server is not reachable")
			break
		}
		time.Sleep(time.Second)
	}

	return err
}

func (f *SognoAdpt) ApplyConfig(prfCfgType string, prfCfg *jsone.A,  cfgObj *jsone.O, config *jsone.O) error {
	// Get the values
	topic := ""
	host := ""
	port := ""

	for _, val := range *prfCfg {
		host, _ = val.(jsone.O).GetString("brokerHost")
		port, _ = val.(jsone.O).GetString("brokerPort")
		topic, _ = val.(jsone.O).GetString("topic")
	}

	err := f.addSouthMQTTService(prfCfgType, topic, host, port, cfgObj, config)
	if err != nil {
		log.Error("South service add failed.", err)
		return err
	}

	return nil
}

func (f *SognoAdpt) addSouthMQTTService(srvType string, topic string, host string, port string, cfg *jsone.O, config *jsone.O) error {
    url := util.GetValuesByKey(cfg, "hostUrl") + _const.SognoPath
	data := fmt.Sprintf(_const.SognoService, host, port, topic)
	log.Infof("data : %s", data)
	_, err := util.SendConfigToNode(url, data, "POST")
	if err != nil {
		log.Error(err, "Adding south service failed.")
		return err
	}

	log.Info("Subscribe MQTT south service success")
	return nil
}

func (f *SognoAdpt) RevertConfig(prfCfgType string, prfCfg *jsone.A, cfg *jsone.O) error {
	return nil
}
