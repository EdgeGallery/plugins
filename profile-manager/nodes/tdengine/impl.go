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

package tdengine

import (
	"fmt"
	"github.com/libujacob/jsone"
	log "github.com/sirupsen/logrus"
	_const "profile-manager/common/const"
	"profile-manager/common/util"
	"profile-manager/nodes/nodeif"
	"strings"
	"time"
)

type TdEngine struct {
	nodeif.Node
	dbName          string
	topic           string
	sTable          string
	tableNameJqPath string
	dataMapping     *jsone.A
}

func (t *TdEngine) ReadinessCheck(node *jsone.O) error {
	return nil
}

func (t *TdEngine) buildAdapterConfigString(key string, val interface{}, config string) string {
	switch key {
	case "tableNameJqPath":
		tableNameJqPath, _ := (val.(jsone.O)).GetString(key)
		config = strings.Replace(config, "_"+key+"_", tableNameJqPath, 1)
	case "dbName":
		dbName, _ := (val.(jsone.O)).GetString(key)
		config = strings.Replace(config, "_"+key+"_", dbName, 1)
	case "sTable":
		sTable, _ := (val.(jsone.O)).GetString(key)
		config = strings.Replace(config, "_"+key+"_", sTable, 1)
	case "topic":
		topic, _ := (val.(jsone.O)).GetString(key)
		config = strings.Replace(config, "_"+key+"_", topic, 1)
	case "dataMapping":
		dataMapping, _ := (val.(jsone.O)).GetArray(key)
		dbMapping := "["
		moreRec := ""
		for _, dbMap := range *dataMapping {
			dbMapping += moreRec
			field, _ := (dbMap.(jsone.O)).GetString("field")
			jqPath, _ := (dbMap.(jsone.O)).GetString("jqPath")
			dataType, _ := (dbMap.(jsone.O)).GetString("dataType")
			dbMapping += "{\"field\": \"" + field + "\",\"jqPath\": \"" + jqPath + "\",\"dataType\": \"" + dataType + "\"}"
			moreRec = ","
		}
		config = strings.Replace(config, "_"+key+"_", dbMapping+"]", 1)
	}
	return config
}

func (t *TdEngine) applyAdapterConfig(configStr string, cfgObj *jsone.O, config *jsone.O) error {
	host, port, err := util.GetBrokerHostAndPort(config)
	if err != nil {
		log.Error(err, "Getting broker host or port failed.")
		return err
	}
	url := util.GetValuesByKey(cfgObj, "hostUrl") + _const.TdEngine
	data := fmt.Sprintf(configStr, host, port)
	log.Infof("data : %s", data)

	_, err = util.SendConfigToNode(url, data, "POST")
	if err != nil {
		log.Errorf(err.Error(), "Sending to tdEngine node failed.")
		return err
	}

	log.Info("TdEngine adapter configuratiin is success")
	time.Sleep(time.Second * 1) //sleep 1 seconds
	return nil
}

func (t *TdEngine) ApplyConfig(prfCfgType string, prfCfg *jsone.A, cfgObj *jsone.O, config *jsone.O) error {
	// Get the values
	configStr := _const.TdEngineFmt
	for _, val := range *prfCfg {
		for _, key := range (val.(jsone.O)).Keys() {
			configStr = t.buildAdapterConfigString(key, val, configStr)
		}
	}

	err := t.applyAdapterConfig(configStr, cfgObj, config)
	if err != nil {
		log.Error(err, "TdEngine apply config failed.")
		return err
	}
	return nil
}

func (t *TdEngine) RevertConfig(prfCfgType string, prfCfg *jsone.A, cfg *jsone.O) error {
	return nil
}
