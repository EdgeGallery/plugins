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
	"fmt"
	"github.com/savaki/jq"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type MqttConn struct {
	client mqtt.Client
}

func NewMQTTConnection(host string, port int) *MqttConn {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", host, port)).SetClientID("emqx_test_client")

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		//fmt.Printf("TOPIC: %s\n", msg.Topic())
		//fmt.Printf("MSG: %s\n", msg.Payload())
		instance := GetDataInstance()
		key := fmt.Sprintf("%s:%d", host, port)
		subKey := fmt.Sprintf("%s:%s", key, msg.Topic())
		for i := 0; i < 10000; i++ {
			op, _ := jq.Parse(fmt.Sprintf(".[%d]", i))
			value, err := op.Apply(msg.Payload())
			if err != nil || value == nil || len(value) == 0 {
				return
			}
			instance.DbConns[key].InsertStore(
				instance.Subscriptions[subKey], value)
		}
	})
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return nil
	}

	return &MqttConn{c}
	// Subscribe to a topic
	//if token := c.Subscribe("Room1/#", 0, nil); token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	os.Exit(1)
	//}
	//
	//// Publish a message
	//token := c.Publish("testtopic/1", 0, false, "Hello World")
	//token.Wait()

	//time.Sleep(6 * time.Second)
	//
	//// Unscribe
	//if token := c.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	os.Exit(1)
	//}
	//
	//// Disconnect
	//c.Disconnect(250)
	//time.Sleep(1 * time.Second)
}

func (m *MqttConn) Subscribe(topic string) bool {
	//Subscribe to a topic
	if token := m.client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return false
	}
	return true
}
