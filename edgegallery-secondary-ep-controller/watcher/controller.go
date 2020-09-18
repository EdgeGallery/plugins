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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/intel/multus-cni/types"
	netattachdef "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	clientset "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/clientset/versioned"
	sharedInformers "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/client/informers/externalversions"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/kubernetes/pkg/api/v1/endpoints"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
)

// Controller main structure.
type Controller struct {
	ClientSet           kubernetes.Interface //Client to get information of k8s resources, svcs, pods etc.
	NetWatcherClientSet clientset.Interface  //Client to get information of custom net-attachment

	KubeInformar informers.SharedInformerFactory

	NetAttachDefsSynced cache.InformerSynced
	NetAttachInformer   sharedInformers.SharedInformerFactory

	PodsLister corelisters.PodLister
	PodsSynced cache.InformerSynced

	ServiceLister  corelisters.ServiceLister
	ServicesSynced cache.InformerSynced

	EndpointsLister corelisters.EndpointsLister
	EndpointsSynced cache.InformerSynced

	Workqueue workqueue.RateLimitingInterface
	Recorder  record.EventRecorder
}

const (
	selectionsKey       = "k8s.v1.cni.cncf.io/networks"
	statusesKey         = "k8s.v1.cni.cncf.io/networks-status"
	controllerAgentName = "edgegallery-secondary-ep-controller"
)

// NewNetworkController returns a new controller structure.
func NewNetworkController(
	clientSet kubernetes.Interface,
	netWatcherClientSet clientset.Interface) *Controller {
	// Records events
	recorder := createRecorder(clientSet, controllerAgentName)

	kubeInformerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*5)
	netAttachDefFactory := sharedInformers.NewSharedInformerFactory(netWatcherClientSet, time.Second*5)

	netWatcherInformer := netAttachDefFactory.K8sCniCncfIo().V1().NetworkAttachmentDefinitions()
	svcInformer := kubeInformerFactory.Core().V1().Services()
	podInformer := kubeInformerFactory.Core().V1().Pods()
	epInformer := kubeInformerFactory.Core().V1().Endpoints()

	NetworkController := &Controller{
		ClientSet:           clientSet,
		KubeInformar:        kubeInformerFactory,
		NetAttachInformer:   netAttachDefFactory,
		NetWatcherClientSet: netWatcherClientSet,
		NetAttachDefsSynced: netWatcherInformer.Informer().HasSynced,
		PodsLister:          podInformer.Lister(),
		PodsSynced:          podInformer.Informer().HasSynced,
		ServiceLister:       svcInformer.Lister(),
		ServicesSynced:      svcInformer.Informer().HasSynced,
		EndpointsLister:     epInformer.Lister(),
		EndpointsSynced:     epInformer.Informer().HasSynced,
		Workqueue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Endpoints"),
		Recorder:            recorder,
	}

	netWatcherInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: NetworkController.HandleNetAttachDefDeleteEvent,
	})

	/* setup handlers for endpoints events */
	epInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    NetworkController.AddOrDelEndpointEvent,
		UpdateFunc: NetworkController.UpdateEndPoint,
		DeleteFunc: NetworkController.AddOrDelEndpointEvent,
	})

	/* setup handlers for services events */
	svcInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    NetworkController.AddServiceEvent,
		UpdateFunc: NetworkController.UpdateSvc,
	})

	/* setup handlers for pod events */
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: NetworkController.UpdatePod,
	})

	return NetworkController
}

//AddOrDelEndpointEvent Handles Add or Delete endpoint event.
func (c *Controller) AddOrDelEndpointEvent(obj interface{}) {
	ep := obj.(*corev1.Endpoints)

	// get service associated with endpoints instance
	svc, err := c.GetServices(ep.GetNamespace(), ep.GetName())
	if err != nil {
		return
	}
	c.AddServiceEvent(svc)
}

//UpdatePod Handles update event.
func (c *Controller) UpdatePod(old, new interface{}) {
	oldPod := old.(*corev1.Pod)
	newPod := new.(*corev1.Pod)

	if oldPod.ResourceVersion == newPod.ResourceVersion || newPod.ObjectMeta.DeletionTimestamp != nil {
		return
	}
	c.addOrDelPodEvent(new)
}

func (c *Controller) addOrDelPodEvent(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}

	_, ok = pod.GetAnnotations()[selectionsKey]
	if !ok {
		log.Info("skipping pod event: network annotations missing")
		return
	}

	// if not behind any service discard
	services, err := GetPodServices(c.ServiceLister, pod)
	if err != nil {
		log.Infof("skipping pod event: %s", err)
		return
	}
	for _, svc := range services {
		c.AddServiceEvent(svc)
	}
}

// AddServiceEvent handles adding service event
func (c *Controller) AddServiceEvent(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.Workqueue.AddRateLimited(key)
}

//UpdateEndPoint handles endpoint update event
func (c *Controller) UpdateEndPoint(old, new interface{}) {
	oldEp := old.(*corev1.Endpoints)
	newEp := new.(*corev1.Endpoints)

	if oldEp.ResourceVersion == newEp.ResourceVersion {
		return
	}
	c.AddOrDelEndpointEvent(new)
}

//UpdateSvc handles sevice change event
func (c *Controller) UpdateSvc(old, new interface{}) {
	oldSvc := old.(*corev1.Service)
	newSvc := new.(*corev1.Service)

	if oldSvc.ResourceVersion == newSvc.ResourceVersion {
		return
	}
	c.AddServiceEvent(new)
}

//Run process evts from workqueue
func (c *Controller) Run(stopChan <-chan struct{}) {
	log.Infof("starting network controller at %s", time.Now().Local())
	defer runtime.HandleCrash()
	defer c.Workqueue.ShutDown()

	c.KubeInformar.Start(stopChan)
	c.NetAttachInformer.Start(stopChan)

	if ok := cache.WaitForCacheSync(stopChan, c.NetAttachDefsSynced, c.EndpointsSynced, c.ServicesSynced, c.PodsSynced); !ok {
		log.Fatalf("failed waiting for caches to sync")
	}

	go wait.Until(c.runWorker, time.Second, stopChan)

	<-stopChan

	log.Infof("shutting down network controller")
	return
}

func (c *Controller) runWorker() {
	for c.ProcessNextWorkItem() {
	}
}

//ProcessNextWorkItem handles svc
func (c *Controller) ProcessNextWorkItem() bool {
	obj, shouldQuit := c.Workqueue.Get()
	if shouldQuit {
		return false
	}
	err := func(obj interface{}) error {
		defer c.Workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.Workqueue.Forget(obj)
			log.Errorf("expected string in workqueue but got %#v", obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.HandleItem(key); err != nil {
			log.Errorf("error syncing '%s': %s", key, err.Error())
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		c.Workqueue.Forget(obj)
		log.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		log.Errorf("error processNextWorkItem %s", err.Error())
		runtime.HandleError(err)
		return true
	}
	return true
}

//HandleItem handles individual items from work queue
func (c *Controller) HandleItem(key string) error {
	pods, svc, ep, networks, err := c.GetResources(key)
	if err != nil {
		return err
	}
	subsets := make([]corev1.EndpointSubset, 0)

	for _, pod := range pods {
		addresses, err := c.GetCurrentAddressList(pod, networks)
		if err != nil {
			continue
		}
		ports, err := c.GetCurrentPortList(pod, svc)
		if err != nil {
			continue
		}

		subset := corev1.EndpointSubset{
			Addresses: addresses,
			Ports:     ports,
		}
		subsets = append(subsets, subset)
	}

	ep.SetOwnerReferences(
		[]metav1.OwnerReference{
			*metav1.NewControllerRef(svc, schema.GroupVersionKind{
				Group:   corev1.SchemeGroupVersion.Group,
				Version: corev1.SchemeGroupVersion.Version,
				Kind:    "Service",
			}),
		},
	)
	log.Infof("subsets %s", subsets)
	ep.Subsets = endpoints.RepackSubsets(subsets)

	// update endpoints resource
	err = c.UpdateEndpoints(ep)
	if err != nil {
		log.Errorf("error updating endpoint: %s", err)
		return err
	}

	log.Info("endpoint updated successfully")
	return nil
}

//GetCurrentAddressList gets addresslist from pods
func (c *Controller) GetCurrentAddressList(pod *corev1.Pod, networks []*types.NetworkSelectionElement) ([]corev1.EndpointAddress, error) {
	networksStatus := make([]types.NetworkStatus, 0)
	err := json.Unmarshal([]byte(pod.Annotations[statusesKey]), &networksStatus)
	if err != nil {
		log.Warningf("error reading pod networks status: %s", err)
		return nil, err
	}
	addresses := make([]corev1.EndpointAddress, 0)
	// find networks used by pod and match network annotation of the service
	for _, status := range networksStatus {
		log.Infof("found pod %s/%s: found network interface %s with IP addresses %s",
			pod.Namespace, pod.Name, status.Interface, status.IPs)
		if IsInNetworkSelectionElementsArray(status.Name, networks) {
			log.Infof("processing pod %s/%s: found network interface %s with IP addresses %s",
				pod.Namespace, pod.Name, status.Interface, status.IPs)
			// all IPs of matching network are added as endpoints
			for _, ip := range status.IPs {
				epAddress := corev1.EndpointAddress{
					IP:       ip,
					NodeName: &pod.Spec.NodeName,
					TargetRef: &corev1.ObjectReference{
						Kind:            "Pod",
						Name:            pod.GetName(),
						Namespace:       pod.GetNamespace(),
						ResourceVersion: pod.GetResourceVersion(),
						UID:             pod.GetUID(),
					},
				}
				addresses = append(addresses, epAddress)
			}
		}
	}
	return addresses, nil
}

//GetCurrentPortList gets port list from svcs
func (c *Controller) GetCurrentPortList(pod *corev1.Pod, svc *corev1.Service) ([]corev1.EndpointPort, error) {
	ports := make([]corev1.EndpointPort, 0)
	for i := range svc.Spec.Ports {
		// check whether pod has the ports needed by service and add them to endpoints if so
		portNumber, err := podutil.FindPort(pod, &svc.Spec.Ports[i])
		if err != nil {
			log.Infof("Could not find pod port for service %s/%s: %s, skipping...", svc.Namespace, svc.Name, err)
			continue
		}

		port := corev1.EndpointPort{
			Port:     int32(portNumber),
			Protocol: svc.Spec.Ports[i].Protocol,
			Name:     svc.Spec.Ports[i].Name,
		}
		ports = append(ports, port)
	}
	return ports, nil
}

//GetResources gets svc, pods, eps and netwroak attacment definition
func (c *Controller) GetResources(key string) ([]*corev1.Pod, *corev1.Service, *corev1.Endpoints, []*types.NetworkSelectionElement, error) {
	log.Infof("key: %s\n", key)
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return nil, nil, nil, nil, errors.New("No valid namespace")
	}
	log.Infof("ns, name: %s %s\n", ns, name)

	svc, err := c.GetServices(ns, name)
	if err != nil {
		return nil, nil, nil, nil, errors.New("service not found")
	}
	annotations := GetNetworkAnnotations(svc)
	if len(annotations) == 0 {
		return nil, nil, nil, nil, errors.New("no network annotations")
	}
	log.Infof("service network annotation found: %v", annotations)
	networks, err := ParsePodNetworkSelections(annotations, ns)
	if err != nil {
		return nil, nil, nil, nil, errors.New("no service networks")
	}
	if len(networks) > 1 {
		msg := fmt.Sprintf("multiple network in service spec")
		log.Warningf(msg)
		return nil, nil, nil, nil, errors.New(msg)
	}

	// get pods matching service selector
	pods, err := c.GetPodsOfService(svc)
	if err != nil {
		return nil, nil, nil, nil, errors.New("No pod exist with matching selector")
	}

	// get endpoints of the service
	ep, err := c.GetEndpoints(ns, name)
	if err != nil {
		return nil, nil, nil, networks, errors.New("no service endpoints found")
	}

	return pods, svc, ep, networks, nil
}

// HandleNetAttachDefDeleteEvent handles network atachment deletion
func (c *Controller) HandleNetAttachDefDeleteEvent(obj interface{}) {
	netAttachDef, ok := obj.(metav1.Object)
	if ok {
		name := netAttachDef.GetName()
		namespace := netAttachDef.GetNamespace()
		log.Infof("handling deletion of %s/%s", namespace, name)
		pods, _ := c.GetAllPods()
		/* check whether net-attach-def requested to be removed is still in use by any of the pods */
		for _, pod := range pods {
			podNetworks, err := c.GetMatchingPodNetworks(pod)
			if err != nil {
				continue
			}
			for _, net := range podNetworks {
				if net.Namespace == namespace && net.Name == name {
					log.Infof("pod %s uses net-attach-def %s/%s which needs to be recreated\n", pod.ObjectMeta.Name, namespace, name)
					/* check whether the object somehow still exists */
					c.RecoverNetworkAttachmentDefinitions(obj, netAttachDef)
				}
			}
		}
	}
}

// GetMatchingPodNetworks Get matching pod networks
func (c *Controller) GetMatchingPodNetworks(pod *corev1.Pod) ([]*types.NetworkSelectionElement, error) {
	netAnnotations, ok := pod.ObjectMeta.Annotations[selectionsKey]
	if !ok {
		return nil, errors.New("No network annotation found")
	}
	podNetworks, err := ParsePodNetworkSelections(netAnnotations, pod.ObjectMeta.Namespace)
	if err != nil {
		return nil, errors.New("No network annotation found")
	}
	return podNetworks, nil
}

// RecoverNetworkAttachmentDefinitions recreating network attacment if svc are pending
func (c *Controller) RecoverNetworkAttachmentDefinitions(obj interface{}, netAttachDef metav1.Object) error {
	_, err := c.NetWatcherClientSet.K8sCniCncfIoV1().
		NetworkAttachmentDefinitions(netAttachDef.GetNamespace()).
		Get(context.TODO(), netAttachDef.GetName(), metav1.GetOptions{})
	if err != nil {
		/* recover deleted object */
		recovered := obj.(*netattachdef.NetworkAttachmentDefinition).DeepCopy()
		recovered.ObjectMeta.ResourceVersion = "" // ResourceVersion field needs to be cleared before recreating the object
		_, err = c.NetWatcherClientSet.
			K8sCniCncfIoV1().
			NetworkAttachmentDefinitions(netAttachDef.GetNamespace()).
			Create(context.TODO(), recovered, metav1.CreateOptions{})
		if err != nil {
			log.Errorf("error recreating recovered object: %s", err.Error())
			return errors.New("error recreating recovered object")
		}
		log.Infof("net-attach-def recovered: %v", recovered)
	}
	return nil
}

//GetPodServices get corresponding svces based on pod selector
func GetPodServices(sl corelisters.ServiceLister, pod *corev1.Pod) ([]*corev1.Service, error) {
	allServices, err := sl.Services(pod.Namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	var services []*corev1.Service
	for i := range allServices {
		service := allServices[i]
		if service.Spec.Selector == nil {
			// services with nil selectors match nothing, not everything.
			continue
		}
		selector := labels.Set(service.Spec.Selector).AsSelectorPreValidated()
		if selector.Matches(labels.Set(pod.Labels)) {
			services = append(services, service)
		}
	}

	return services, nil
}

//GetServices Get all services belongs to namespaces
func (c *Controller) GetServices(namespaces string, name string) (*corev1.Service, error) {
	svc, err := c.ServiceLister.Services(namespaces).Get(name)
	if err != nil {
		return nil, errors.New("service not found")
	}
	return svc, nil
}

//GetPodsOfService filter services related to pods
func (c *Controller) GetPodsOfService(svc *corev1.Service) ([]*corev1.Pod, error) {
	selector := labels.Set(svc.Spec.Selector).AsSelector()
	pods, err := c.PodsLister.List(selector)
	if err != nil {
		// no selector or no pods running
		log.Warn("error listing pods matching service selector: %s", selector.String())
		return nil, errors.New("No pod exist with matching selector")
	}
	return pods, err
}

//GetEndpoints find endpoints by name
func (c *Controller) GetEndpoints(namespaces string, name string) (*corev1.Endpoints, error) {
	ep, err := c.EndpointsLister.Endpoints(namespaces).Get(name)
	if err != nil {
		log.Infof("error getting service endpoints: %s", name)
		return nil, errors.New("no service endpoints found")
	}
	return ep, nil
}

//UpdateEndpoints call update endpoint api
func (c *Controller) UpdateEndpoints(ep *corev1.Endpoints) error {
	_, err := c.ClientSet.CoreV1().Endpoints(ep.Namespace).Update(context.TODO(), ep, metav1.UpdateOptions{})
	if err != nil {
		log.Errorf("error getting service endpoints: %s", err.Error())
		return errors.New("endpoint update failed")
	}
	return nil
}

//GetAllPods get All pods
func (c *Controller) GetAllPods() ([]*corev1.Pod, error) {
	pods, err := c.PodsLister.Pods("").List(labels.Everything())
	if err != nil {
		return nil, errors.New("No pod found")
	}
	return pods, nil
}