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

package _const

const (
	FledgeStr   = "fledge"
	KuiperStr   = "kuiper"
	TdEngineStr = "tdengine"
	SognoStr = "sogno"
)

const (
	ConfigPath   = "conf/config.yaml"
	ProficlePath = "conf/profile.yaml"
)

const MaxNumberRetry = 100

const PluginFormat string = "{\"format\":\"%s\",\"name\":\"%s\",\"version\":\"%s\"}"
const PluginPath = "/fledge/plugins"

//Fledge constants
const SubscribeService = "{\"name\": \"%s\", \"type\": \"%s\", \"plugin\" :\"%s\",\"enabled\":true,\"config\": {\"brokerHost\": {\"value\": \"%s\"},\"brokerPort\": {\"value\": \"%d\"},\"topic\": {\"value\": \"%s\"},\"assetName\":{\"value\":\"%s\"}}}"
const SubscribeOPCUAService = "{\"name\": \"%s\", \"type\": \"%s\", \"plugin\" :\"%s\",\"enabled\":true,\"config\": {\"asset\":{\"value\":\"opcua00\"},\"url\":{\"value\":\"%s\"},\"subscription\": {\"value\":{\"subscriptions\": [\"%s\"]}}}}"
const SubscribeModbusService = "{\"name\": \"%s\", \"type\": \"%s\", \"plugin\" :\"%s\",\"enabled\":true,\"config\": " +
	"{\"protocol\":{\"value\":\"TCP\"},\"address\":{\"value\": \"%s\"}, \"port\":{\"value\": \"%s\"},\"map\":{\"value\":{\"values\":[{\"name\":\"%s\",\"slave\":%d,\"assetName\":\"%s\",\"register\":%d,\"scale\":%f,\"offset\":%d}]}}}}"

const SubscribeDNP3Service = "{\"name\": \"%s\", \"type\": \"%s\", \"plugin\" :\"%s\",\"enabled\":true,\"config\": " +
	"{\"asset\":{\"value\":\"%s\"},\"master_id\":{\"value\": \"%s\"},\"outstation_tcp_address\":{\"value\": \"%s\"},\"outstation_tcp_port\":{\"value\": \"%s\"},\"outstation_id\":{\"value\": \"%s\"},\"outstation_scan_interval\":{\"value\": \"%s\"},\"data_fetch_timeout\":{\"value\": \"%s\"}}}"

const SubscribeCsvService = "{\"name\": \"%s\", \"type\": \"%s\", \"plugin\" :\"%s\",\"enabled\":true,\"config\": {\"asset\": {\"value\": \"%s\"}, \"datapoint\": {\"value\": \"%s\"},\"file\": {\"value\": \"%s\"}}}"

const SubscribePath = "/fledge/service"
const NorthService = "{\"name\":\"%s\", \"plugin\":\"mqtt_north\", \"type\":\"%s\", \"enabled\":true, \"config\": {\"brokerHost\": {\"value\": \"%s\"},\"brokerPort\": {\"value\": \"%d\"},\"topic\": {\"value\": \"%s\"}}}"

//Kuiper constants
const Streams = "/streams"
const Rules = "/rules"
const ActionServer = "tcp://%s:%d"
const RulesFmt = "{\"id\": \"rule2\",\"sql\": \"%s\",\"actions\": [{\"mqtt\": {\"server\": \"%s\",\"topic\": \"kuiper\"}}]}"

// TdEngine
const TdEngineFmt = "{\"brokerAddr\": {\"host\": \"%s\",\"port\": %d},\"dbAddr\": {\"host\": \"0.0.0.0\"," +
	"\"port\": 6030},\"dbName\": \"_dbName_\",\"store\": [{\"topic\": \"_topic_\",\"sTable\": \"_sTable_\"," +
	"\"tableNameJqPath\": \"_tableNameJqPath_\",\"dataMapping\": _dataMapping_}]}"
const TdEngine = "/api/v1/resource"
const Delay = 30

//sogno
const SognoService = "{\"brokerHost\": \"%s\", \"brokerPort\": \"%s\", \"topic\": \"%s\"}"
const SognoPath = "/sogno/fledgeadapter/topic"