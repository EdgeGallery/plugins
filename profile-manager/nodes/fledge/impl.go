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

package fledge

import (
	"fmt"
	"github.com/libujacob/jsone"
	log "github.com/sirupsen/logrus"
	_const "profile-manager/common/const"
	"profile-manager/common/util"
	"profile-manager/nodes/nodeif"
	"time"
)

type FlEdge struct {
	nodeif.Node
}

func (f *FlEdge) ReadinessCheck(node *jsone.O) error {

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

func GetTypeByKey(cfg jsone.O) string {
	var result string = ""
	for _, k := range cfg.Keys() {
		if k == "type" {
			result, _ = cfg.GetString(k)
		}
	}
	return result
}

func (f *FlEdge) ApplyConfig(prfCfgType string, prfCfg *jsone.A,  cfgObj *jsone.O, config *jsone.O) error {
	install := false
	topic := ""
	srvUrl := ""
	nodeId := ""
	_type := ""
	// Get the values
	for _, val := range *prfCfg {
		if GetTypeByKey(val.(jsone.O)) == "mqtt-readings" {
			for _, key := range (val.(jsone.O)).Keys() {
				switch key {
				case "installPlugin":
					install, _ = (val.(jsone.O)).GetBoolean(key)
				case "topic":
					topic, _ = (val.(jsone.O)).GetString(key)
				case "type":
					_type, _ = (val.(jsone.O)).GetString(key)
				}
			}

			//1.install plugin
			if install {
				err := f.installPlugin(_type, cfgObj)
				if err != nil {
					log.Error("Install plugin failed.", err)
					return err
				}
			}

			//2.add south service
			if prfCfgType == "south" {
				err := f.addSouthMQTTService(prfCfgType, topic, _type, cfgObj, config)
				if err != nil {
					log.Error("South service add failed.", err)
					return err
				}
			}
		}
		if GetTypeByKey(val.(jsone.O)) == "opcua" {
			for _, key := range (val.(jsone.O)).Keys() {
				switch key {
				case "installPlugin":
					install, _ = (val.(jsone.O)).GetBoolean(key)
				case "srvUrl":
					srvUrl, _ = (val.(jsone.O)).GetString(key)
				case "type":
					_type, _ = (val.(jsone.O)).GetString(key)
				case "nodeId":
					nodeId, _ = (val.(jsone.O)).GetString(key)
				}
			}

			//1.install plugin
			if install {
				err := f.installPlugin(_type, cfgObj)
				if err != nil {
					log.Error("Install plugin failed.", err)
					return err
				}
			}

			//2.add south service
			if prfCfgType == "south" {
				err := f.addSouthOPCUAService(prfCfgType, srvUrl, nodeId, _type, cfgObj, config)
				if err != nil {
					log.Error("South service add failed.", err)
					return err
				}
			}
		}
	}

	//3.add north service : defaut
	err := f.addNorthService("north", "", "", cfgObj, config)
	if err != nil {
		log.Error("North service add failed.", err)
		return err
	}

	return nil
}

func (f *FlEdge) installPlugin(name string, cfg *jsone.O) error {
	repository := "repository"
	version := ""

	url := util.GetValuesByKey(cfg, "hostUrl") + _const.PluginPath
	data := fmt.Sprintf(_const.PluginFormat, repository, "fledge-south-"+name, version)
	log.Infof("data : %s", data)

	_, err := util.SendConfigToNode(url, data, "POST")
	if err != nil {
		log.Errorf(err.Error(), "plugin failed while sending.")
		return err
	}

	log.Info("Node plugin install success")
	time.Sleep(time.Second * _const.Delay) //sleep 10 seconds
	return nil
}

func (f *FlEdge) addSouthMQTTService(srvType string, topic string, _type string, cfg *jsone.O, config *jsone.O) error {
    host, port, err := util.GetBrokerHostAndPort(config)
	if err != nil {
		log.Error(err, "Getting broker host or port failed.")
		return err
	}
	url := util.GetValuesByKey(cfg, "hostUrl") + _const.SubscribePath
	data := fmt.Sprintf(_const.SubscribeService, "in", srvType, _type, host, port, topic, "mqtt123")
	log.Infof("data : %s", data)
	_, err = util.SendConfigToNode(url, data, "POST")
	if err != nil {
		log.Error(err, "Adding south service failed.")
		return err
	}

	log.Info("Subscribe MQTT south service success")
	return nil
}

func (f *FlEdge) addSouthOPCUAService(srvType string, srvUrl string, nodeId string, _type string, cfg *jsone.O, config *jsone.O) error {

	url := util.GetValuesByKey(cfg, "hostUrl") + _const.SubscribePath
	data := fmt.Sprintf(_const.SubscribeOPCUAService, "opcua-srv", srvType, _type, srvUrl, nodeId)
	log.Infof("data : %s", data)
	_, err := util.SendConfigToNode(url, data, "POST")
	if err != nil {
		log.Error(err, "Adding south service failed.")
		return err
	}

	log.Info("Subscribe opcua south service success")
	return nil
}

func (f *FlEdge) addNorthService(srvType string, topic string, _type string, cfg *jsone.O, config *jsone.O) error {
	url := util.GetValuesByKey(cfg, "hostUrl") + _const.SubscribePath
	host, port, err := util.GetBrokerHostAndPort(config)
	if err != nil {
		log.Error(err, "Getting broker host or port failed.")
		return err
	}
	north := util.GetDefaultByKey(cfg, "north")
	if north != nil {
		if topic == "" {
			topic = util.GetValuesByKey(north, "topic")
		}
		if _type == "" {
			_type = util.GetValuesByKey(north, "type")
		}
	}

	data := fmt.Sprintf(_const.NorthService, "out", srvType, host, port, topic)
	log.Infof("data : %s", data)
	_, err = util.SendConfigToNode(url, data, "POST")
	if err != nil {
		log.Errorf(err.Error(), "sending north service failed.")
		return err
	}

	log.Info("Subscribe north service success")
	time.Sleep(time.Second * _const.Delay) //sleep 10 seconds
	return nil
}

func (f *FlEdge) RevertConfig(prfCfgType string, prfCfg *jsone.A, cfg *jsone.O) error {
	return nil
}
