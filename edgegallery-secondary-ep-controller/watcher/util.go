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

package watcher

import (
	"errors"
	"regexp"
	"strings"

	"github.com/intel/multus-cni/types"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

//GetNetworkAnnotations Read network annotation from svc
func GetNetworkAnnotations(obj interface{}) string {
	metaObject := obj.(metav1.Object)
	annotations, ok := metaObject.GetAnnotations()[selectionsKey]
	if !ok {
		return ""
	}
	return annotations
}

//ParsePodNetworkSelections Parse network selection
func ParsePodNetworkSelections(podNetworks string, defaultNamespace string) ([]*types.NetworkSelectionElement, error) {
	var networkSelections []*types.NetworkSelectionElement

	if len(podNetworks) == 0 {
		err := errors.New("empty string passed as network selection elements list")
		log.Error(err)
		return nil, err
	}

	for _, networkSelection := range strings.Split(podNetworks, ",") {
		networkSelection = strings.TrimSpace(networkSelection)
		networkSelectionElement, err := ParsePodNetworkSelectionElement(networkSelection, defaultNamespace)
		if err != nil {
			err := errors.New("error parsing network selection element")
			log.Error(err)
			return nil, err
		}
		networkSelections = append(networkSelections, networkSelectionElement)
	}

	return networkSelections, nil
}

//ResolveNetworkAnnotation parse network annotation from svc
func ResolveNetworkAnnotation(annotation string, defaultNamespace string) (string, string, error) {
	var namespace, name string
	units := strings.Split(annotation, "/")
	if len(units) == 1 {
		namespace = defaultNamespace
		name = units[0]
	} else if len(units) == 2 {
		namespace = units[0]
		name = units[1]
	} else {
		err := errors.New("invalid network selection element")
		log.Error(err)
		return "", "", err
	}
	return namespace, name, nil

}

//ParsePodNetworkSelectionElement parse network selection element
func ParsePodNetworkSelectionElement(selection string, defaultNamespace string) (*types.NetworkSelectionElement, error) {
	var namespace, name string
	var networkSelectionElement *types.NetworkSelectionElement

	namespace, name, err := ResolveNetworkAnnotation(selection, defaultNamespace)
	if err != nil {
		log.Error(err)
		return networkSelectionElement, err
	}

	validNameRegex, _ := regexp.Compile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	for _, unit := range []string{namespace, name} {
		ok := validNameRegex.MatchString(unit)
		if !ok && len(unit) > 0 {
			err := errors.New("invalid network annotation")
			log.Error(err)
			return networkSelectionElement, err
		}
	}

	networkSelectionElement = &types.NetworkSelectionElement{
		Namespace: namespace,
		Name:      name,
	}
	return networkSelectionElement, nil
}

//IsInNetworkSelectionElementsArray finds if networks belongs to list
func IsInNetworkSelectionElementsArray(name string, networks []*types.NetworkSelectionElement) bool {
	for i := range networks {
		log.Infof("checking service network %s === pod network %s ", name, networks[i].Name)
		_, netname, err := ResolveNetworkAnnotation(name, networks[i].Name)
		if err != nil {
			return false
		}
		if netname == networks[i].Name {
			return true
		}
	}
	return false
}

func createRecorder(clientSet kubernetes.Interface, comp string) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(log.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: clientSet.CoreV1().Events("")})
	return eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: comp})
}