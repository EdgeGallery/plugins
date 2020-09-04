/*
 *  Copyright 2020 Huawei Technologies Co., Ltd.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package main

import (
	"flag"
	clientset "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"edgegallery-secondary-ep-controller/watcher"
	"os"
	"os/signal"
)

var (
	masterURL  string
	kubeconfigPath string
)

func main() {

	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Required if out-of-cluster.")
	flag.StringVar(&kubeconfigPath, "kubeconfig", "", "Path to a kubeconfig. Required if out-of-cluster.")
	flag.Parse()

	kubecfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
		return
	}

	kubeClientSet, err := kubernetes.NewForConfig(kubecfg)
	if err != nil {
		log.Fatalf("error building kubernetes clientset: %s", err.Error())
		return
	}

	netAttachDefClientSet, err := clientset.NewForConfig(kubecfg)
	if err != nil {
		log.Fatalf("error creating net-attach-def clientset: %s", err.Error())
		return
	}

	networkController := watcher.NewNetworkController(
		kubeClientSet,
		netAttachDefClientSet)

	stopChan := make(chan struct{})
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		close(stopChan)
		<-signals
		os.Exit(1)
	}()

	networkController.Run(stopChan)
}

