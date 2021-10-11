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

package kuiper

import (
	"fmt"
	"github.com/libujacob/jsone"
	log "github.com/sirupsen/logrus"
	_const "profile-manager/common/const"
	"profile-manager/common/util"
	"profile-manager/nodes/nodeif"
	"time"
)

type Kuiper struct {
	nodeif.Node
	streamSql string
	rulesSql  string
}

func (f *Kuiper) ReadinessCheck(cfg *jsone.O) error {
	return nil
}

func addStreamsWithZmq(url string, sql string) error {
	data := "{\"sql\":\"" + sql + "\"}"
	log.Infof("data : %s", data)
	_, err := util.SendConfigToNode(url+_const.Streams, data, "POST")
	if err != nil {
		log.Errorf(err.Error(), "sending streams with zmq failed.")
		return err
	}

	log.Info("Add streams with zmq is success")
	time.Sleep(time.Second * 1) //sleep 1 seconds
	return nil
}

//addRulesWithZmq - API to add streams with zmq
func addRulesWithZmq(url string, sql string, config *jsone.O) error {
	host, port, err := util.GetBrokerHostAndPort(config)
	if err != nil {
		log.Error(err, "Getting broker host or port failed.")
		return err
	}

	server := fmt.Sprintf(_const.ActionServer, host, port)
	data := fmt.Sprintf(_const.RulesFmt, sql, server)
	log.Infof("data : %s", data)
	_, err = util.SendConfigToNode(url+_const.Rules, data, "POST")
	if err != nil {
		log.Errorf(err.Error(), "Sending add rules with zmq is failed.")
		return err
	}

	log.Info("adding rules with zmq is success")
	return nil
}

func (f *Kuiper) ApplyConfig(prfCfgType string, prfCfg *jsone.A, cfgObj *jsone.O, config *jsone.O) error {
	url := util.GetValuesByKey(cfgObj, "hostUrl")

	for _, sqls := range *prfCfg {
		sql, _ := (sqls.(jsone.O)).GetString("sql")
		switch prfCfgType {
		case "stream":
			f.streamSql = sql
		case "rules":
			f.rulesSql = sql
		}
	}

	if len(f.rulesSql) == 0 || len(f.streamSql) == 0 {
		return nil
	}
	err := addStreamsWithZmq(url, f.streamSql)
	if err != nil {
		log.Errorf(err.Error(), "Add stream is failed.")
		return err
	}

	err = addRulesWithZmq(url, f.rulesSql, config)
	if err != nil {
		log.Errorf(err.Error(), "Add rules are failed.")
		return err
	}

	return nil
}

func (f *Kuiper) RevertConfig(prfCfgType string, prfCfg *jsone.A, cfg *jsone.O) error {
	return nil
}
