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

import (
	"database/sql"
	"fmt"
	"github.com/savaki/jq"
	_ "github.com/taosdata/driver-go/taosSql"
	"strings"
	"time"
)

type TDEngine struct {
	url string
	db  *sql.DB
}

var taosDriverName = "taosSql"

//
//func main() {
//	store := model.DbConns{
//		Topic:           "abcd",
//		STable:          "meters1",
//		TableNameJqPath: ".asset",
//		DataMapping:     make([]model.DataMapping, 0),
//	}
//	store.DataMapping = append(store.DataMapping, model.DataMapping{
//		Field:    "ts",
//		JqPath:   ".timestamp",
//		DataType: "TIMESTAMP",
//	})
//	store.DataMapping = append(store.DataMapping, model.DataMapping{
//		Field:    "temperature",
//		JqPath:   ".readings.temp",
//		DataType: "TINYINT",
//	})
//	subscription := model.Subscription{
//		BrokerHost: "172.16.105.1",
//		BrokerPort: 6030,
//		DBName:     "iotdb",
//		DbConns:      make([]model.DbConns, 0),
//	}
//	subscription.DbConns = append(subscription.DbConns, store)
//	newStore := NewStore(&subscription)
//	fmt.Println(newStore.CheckStableExists("meters1"))
//	fmt.Println(newStore.CreateStable(&store))
//	fmt.Println(newStore.InsertStore(&store, []byte("{\"asset\":\"asset-101\","+
//		"\"timestamp\":\"2021-09-20 12:14:47.317\",\"readings\":{\"temp\":67}}")))
//}

func NewStore(sub *Subscription) *TDEngine {
	dbName := sub.DBName
	engine := &TDEngine{url: fmt.Sprintf("root:taosdata@/tcp(%s:%d)/", sub.DbAddr.Host, sub.DbAddr.Port)}
	db, err := sql.Open(taosDriverName, engine.url)
	if err != nil {
		fmt.Printf("Open database error: %s\n", err)
		return nil
	}
	engine.db = db
	if !engine.execute(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName)) {
		engine.Close()
		return nil
	}
	if !engine.execute(fmt.Sprintf("USE %s;", dbName)) {
		engine.Close()
		return nil
	}
	return engine
}

func (e *TDEngine) Close() {
	e.db.Close()
}

func (e TDEngine) execute(statement string) bool {
	fmt.Printf("- %s\n", statement)
	_, err := e.db.Exec(statement)
	if err != nil {
		return false
	}
	return true
}

func (e *TDEngine) CheckStableExists(stableName string) bool {

	sqlStr := "show stables;"
	fmt.Printf("- %s\n", sqlStr)
	res, err := e.db.Query(sqlStr)
	if !e.checkErr(err, sqlStr) {
		return false
	}
	defer res.Close()
	for res.Next() {
		var (
			name    string
			ts      time.Time
			columns int64
			tags    int64
			tables  int64
		)
		err = res.Scan(&name, &ts, &columns, &tags, &tables)
		if !e.checkErr(err, sqlStr) {
			return false
		}
		if name == stableName {
			return true
		}
	}
	return false
}

func (e *TDEngine) CreateStable(store *Store) bool {
	sqlStr := "create stable if not exists " + store.STable + " ("
	for _, field := range store.DataMapping {
		sqlStr += field.Field + " " + field.DataType + ","
	}
	sqlStr = sqlStr[:len(sqlStr)-1]
	sqlStr += ") tags(asset binary(64));"
	fmt.Printf("- %s\n", sqlStr)
	_, err := e.db.Exec(sqlStr)
	return e.checkErr(err, sqlStr)
}

func (e *TDEngine) InsertStore(store *Store, data []byte) bool {
	op, _ := jq.Parse(store.TableNameJqPath)
	value, _ := op.Apply(data)
	tblName := strings.Replace(string(value), "-", "_", -1)
	sqlStr := "INSERT INTO " + tblName
	fieldNames := make([]string, 0)
	values := make([]string, 0)
	for _, field := range store.DataMapping {
		op1, _ := jq.Parse(field.JqPath)
		value1, _ := op1.Apply(data)
		fieldNames = append(fieldNames, field.Field)
		values = append(values, string(value1))
	}
	sqlStr += " (" + strings.Join(fieldNames, ",") + ") "
	sqlStr += " USING " + store.STable + " TAGS (" + tblName + ") VALUES "
	sqlStr += " (" + strings.Join(values, ",") + ");"

	//fmt.Printf("- %s\n", sqlStr)
	_, err := e.db.Exec(sqlStr)
	return e.checkErr(err, sqlStr)
}

func (e *TDEngine) checkErr(err error, prompt string) bool {
	if err != nil {
		fmt.Printf("ERROR: %s\n", prompt)
		return false
	}
	return true
}
