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

package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	brokerHost := os.Getenv("BROKER_HOST")
	brokerPort := os.Getenv("BROKER_PORT")
	devCountSt := os.Getenv("DEVICE_COUNT")
	topic := os.Getenv("TOPIC")
	if len(topic) == 0 {
		topic = "Room1/conditions"
	}
	deviceCount := 1000

	if len(devCountSt) != 0 {
		i, err := strconv.Atoi(devCountSt)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		deviceCount = i
	}

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", brokerHost,
		brokerPort)).SetClientID("device_simulate_client")

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for i := 0; i < deviceCount; i++ {
		go func(counter int) {
			r := rand.New(rand.NewSource(99))
			for true {
				_ = c.Publish(topic, 0, false,
					fmt.Sprintf("{\"device\": \"asset-%04d\", \"humidity\": %.2f, \"temp\": %.2f}",
						counter,
						float32(50+r.Intn(50))+r.Float32(),
						float32(10+r.Intn(30))+r.Float32()))
				//token.Wait()
				time.Sleep(1 * time.Second)
			}
		}(i)
	}

	for true {
		time.Sleep(5 * time.Second)
	}
}
