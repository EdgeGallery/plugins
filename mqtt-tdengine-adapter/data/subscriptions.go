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

package data

type Subscription struct {
	BrokerAddr Address `json:"brokerAddr"`
	DbAddr     Address `json:"dbAddr"`
	DBName     string  `json:"dbName"`
	Store      []Store `json:"store"`
}

type Address struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Store struct {
	Topic           string        `json:"topic"`
	STable          string        `json:"sTable"`
	TableNameJqPath string        `json:"tableNameJqPath"`
	DataMapping     []DataMapping `json:"dataMapping"`
}

type DataMapping struct {
	Field    string `json:"field"`
	JqPath   string `json:"jqPath"`
	DataType string `json:"dataType"`
}
