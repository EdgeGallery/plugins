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

type ModBusInfo struct {
	devIP string
	devPort string
	reading string
	slaveAdd int64
	assetName string
	register int64
	scale float64
	offset int64
}
type DNP3Info struct {
	asset                    string
	master_id                string
	outstation_tcp_address   string
	outstation_tcp_port      string
	outstation_id            string
	outstation_scan_interval string
	data_fetch_timeout       string
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
		if GetTypeByKey(val.(jsone.O)) == "modbus" {
			var modbus ModBusInfo
			fmt.Println(val.(jsone.O))
			for _, key := range (val.(jsone.O)).Keys() {
				switch key {
				case "installPlugin":
					install, _ = (val.(jsone.O)).GetBoolean(key)
				case "devIP":
					modbus.devIP, _ = (val.(jsone.O)).GetString(key)
				case "devPort":
					modbus.devPort, _ = (val.(jsone.O)).GetString(key)
				case "reading":
					modbus.reading, _ = (val.(jsone.O)).GetString(key)
				case "slaveAdd":
					modbus.slaveAdd, _ = (val.(jsone.O)).GetInt64(key)
				case "assetName":
					modbus.assetName, _ = (val.(jsone.O)).GetString(key)
				case "register":
					modbus.register, _ = (val.(jsone.O)).GetInt64(key)
				case "scale":
					modbus.scale, _ = (val.(jsone.O)).GetFloat64(key)
				case "offset":
					modbus.offset, _ = (val.(jsone.O)).GetInt64(key)
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
				err := f.addSouthModbusService(prfCfgType, modbus, cfgObj)
				if err != nil {
					log.Error("South service add failed.", err)
					return err
				}
			}
		}
		if GetTypeByKey(val.(jsone.O)) == "dnp3" {
			f.configdnp3Device(val.(jsone.O), prfCfgType, cfgObj, config)
		}
		if GetTypeByKey(val.(jsone.O)) == "csv" {
			f.configCsvDevice(val.(jsone.O), prfCfgType, cfgObj, config)
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

func (f *FlEdge) configdnp3Device(val jsone.O, prfCfgType string, cfgObj *jsone.O, config *jsone.O) error {
	install := false
	var dnp3Info DNP3Info
	_type := ""
	for _, key := range val.Keys() {
		switch key {
		case "installPlugin":
			install, _ = val.GetBoolean(key)
		case "outstation_tcp_port":
			dnp3Info.outstation_tcp_port, _ = val.GetString(key)
		case "outstation_tcp_address":
			dnp3Info.outstation_tcp_address, _ = val.GetString(key)
		case "asset":
			dnp3Info.asset, _ = val.GetString(key)
		case "master_id":
			dnp3Info.master_id, _ = val.GetString(key)
		case "outstation_id":
			dnp3Info.outstation_id, _ = val.GetString(key)
		case "outstation_scan_interval":
			dnp3Info.outstation_scan_interval, _ = val.GetString(key)
		case "data_fetch_timeout":
			dnp3Info.data_fetch_timeout, _ = val.GetString(key)
		case "type":
			_type, _ = val.GetString(key)
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
		err := f.addSouthDnp3Service(prfCfgType, dnp3Info, _type, cfgObj, config)
		if err != nil {
			log.Error("South service add failed.", err)
			return err
		}
	}

	log.Info("DNP3 plugin and south service config success")
	return nil
}

func (f *FlEdge) configCsvDevice(val jsone.O, prfCfgType string, cfgObj *jsone.O, config *jsone.O) error {
	install := false
	asset := ""
	file := ""
	datapoint := ""
	_type := ""
	for _, key := range val.Keys() {
		switch key {
		case "installPlugin":
			install, _ = val.GetBoolean(key)
		case "asset":
			asset, _ = val.GetString(key)
		case "file":
			file, _ = val.GetString(key)
		case "datapoint":
			datapoint, _ = val.GetString(key)
		case "type":
			_type, _ = val.GetString(key)
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
		err := f.addSouthCsvService(prfCfgType, asset, file, datapoint, _type, cfgObj, config)
		if err != nil {
			log.Error("South service add failed.", err)
			return err
		}
	}

	log.Info("Csv plugin and south service config success")
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

func (f *FlEdge) addSouthModbusService(srvType string, modbus ModBusInfo, cfg *jsone.O) error {

	url := util.GetValuesByKey(cfg, "hostUrl") + _const.SubscribePath
	log.Infof("data : %v", modbus)

	data := fmt.Sprintf(_const.SubscribeModbusService, "modbus-srvc", srvType, "ModbusC", modbus.devIP, modbus.devPort, modbus.reading, modbus.slaveAdd, modbus.assetName, modbus.register, modbus.scale, modbus.offset)
	log.Infof("data : %s", data)
	_, err := util.SendConfigToNode(url, data, "POST")
	if err != nil {
		log.Error(err, "Adding south service failed.")
		return err
	}

	log.Info("Subscribe modbus south service success")
	return nil
}

func (f *FlEdge) addSouthDnp3Service(srvType string, dnp3Info DNP3Info, _type string, cfg *jsone.O, config *jsone.O) error {

	url := util.GetValuesByKey(cfg, "hostUrl") + _const.SubscribePath
	data := fmt.Sprintf(_const.SubscribeDNP3Service, "dnp3-srvc", srvType, _type, dnp3Info.asset, dnp3Info.master_id, dnp3Info.outstation_tcp_address,
		dnp3Info.outstation_tcp_port, dnp3Info.outstation_id, dnp3Info.outstation_scan_interval, dnp3Info.data_fetch_timeout)
	_, err := util.SendConfigToNode(url, data, "POST")
	if err != nil {
		log.Error(err, "Adding south service failed.")
		return err
	}

	log.Info("Subscribe dnp3 south service success")
	return nil
}

func (f *FlEdge) addSouthCsvService(srvType string, asset string, file string, datapoint string, _type string, cfg *jsone.O, config *jsone.O) error {

	url := util.GetValuesByKey(cfg, "hostUrl") + _const.SubscribePath
	data := fmt.Sprintf(_const.SubscribeCsvService, "Csv-srvc", srvType, "Csv", asset, datapoint, file)
	_, err := util.SendConfigToNode(url, data, "POST")
	if err != nil {
		log.Error(err, "Adding south service failed.")
		return err
	}

	log.Info("Subscribe CSV south service success")
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
