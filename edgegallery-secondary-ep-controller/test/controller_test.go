package test

import (
	"edgegallery-secondary-ep-controller/watcher"
	"errors"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/intel/multus-cni/types"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	CONTROLLER_INIT string = "Error in controller init"
	ENDPOINTS       string = "Endpoints"
	DEFAULT         string = "default"
	WRK_QUEUE_EMPTY string = "Workeue should be empty"
	SERVICE         string = "Service"
	V1              string = "v1"
	IP1             string = "100.1.1.1"
	HTTP_PORT       string = "httpPort"
)

func TestInitNetworkController(t *testing.T) {
	convey.Convey("Testing Init network controller", t, func() {
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)
	})
}

func TestAddOrDelEndpointEvent(t *testing.T) {
	convey.Convey("Testing Successfully added services to workqueue after Endpoint event received", t, func() {
		fakeEp := &corev1.Endpoints{
			TypeMeta: metav1.TypeMeta{
				Kind:       ENDPOINTS,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: DEFAULT,
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							IP: IP1,
						},
					},
					Ports: []corev1.EndpointPort{
						{
							Name:     HTTP_PORT,
							Port:     30001,
							Protocol: corev1.ProtocolTCP,
						},
					},
				},
			},
		}
		fakeservice := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: DEFAULT,
				Name:      "foo",
			},
		}
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)
		var n *watcher.Controller
		patch1 := gomonkey.ApplyMethod(reflect.TypeOf(n), "GetServices", func(*watcher.Controller, string, string) (*corev1.Service, error) {
			return fakeservice, nil
		})
		networkController.AddOrDelEndpointEvent(fakeEp)
		obj, shouldQuit := networkController.Workqueue.Get()
		assert.NotEqual(t, nil, obj, CONTROLLER_INIT)
		assert.Equal(t, false, shouldQuit, CONTROLLER_INIT)

		defer patch1.Reset()
	})
}

func TestAddOrDelEndpointEventServiceNotfound(t *testing.T) {
	convey.Convey("Testing Successfully added services to workqueue after Endpoint event received", t, func() {
		fakeEp := &corev1.Endpoints{
			TypeMeta: metav1.TypeMeta{
				Kind:       ENDPOINTS,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: DEFAULT,
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							IP: IP1,
						},
					},
					Ports: []corev1.EndpointPort{
						{
							Name:     HTTP_PORT,
							Port:     30001,
							Protocol: corev1.ProtocolTCP,
						},
					},
				},
			},
		}
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)
		var n *watcher.Controller
		patch1 := gomonkey.ApplyMethod(reflect.TypeOf(n), "GetServices", func(*watcher.Controller, string, string) (*corev1.Service, error) {
			return nil, errors.New("service not found")
		})
		networkController.AddOrDelEndpointEvent(fakeEp)
		assert.Equal(t, 0, networkController.Workqueue.Len(), WRK_QUEUE_EMPTY)
		defer patch1.Reset()
	})
}

func TestUpdateEndPointEvent(t *testing.T) {
	convey.Convey("Testing Successfully added services to workqueue after Update Endpoint event received", t, func() {
		fakeEpOld := &corev1.Endpoints{
			TypeMeta: metav1.TypeMeta{
				Kind:       ENDPOINTS,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            "foo",
				Namespace:       DEFAULT,
				ResourceVersion: "1",
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							IP: IP1,
						},
					},
					Ports: []corev1.EndpointPort{
						{
							Name:     HTTP_PORT,
							Port:     30001,
							Protocol: corev1.ProtocolTCP,
						},
					},
				},
			},
		}
		fakeEpNew := &corev1.Endpoints{
			TypeMeta: metav1.TypeMeta{
				Kind:       ENDPOINTS,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            "foo",
				Namespace:       DEFAULT,
				ResourceVersion: "2",
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							IP: IP1,
						},
					},
					Ports: []corev1.EndpointPort{
						{
							Name:     HTTP_PORT,
							Port:     30001,
							Protocol: corev1.ProtocolTCP,
						},
					},
				},
			},
		}
		fakeservice := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: DEFAULT,
				Name:      "foo",
			},
		}
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)
		var n *watcher.Controller
		patch1 := gomonkey.ApplyMethod(reflect.TypeOf(n), "GetServices", func(*watcher.Controller, string, string) (*corev1.Service, error) {
			return fakeservice, nil
		})
		networkController.UpdateEndPoint(fakeEpOld, fakeEpNew)
		obj, shouldQuit := networkController.Workqueue.Get()
		assert.NotEqual(t, nil, obj, CONTROLLER_INIT)
		assert.Equal(t, false, shouldQuit, CONTROLLER_INIT)

		defer patch1.Reset()
	})
}

func TestAddServiceEventSuccess(t *testing.T) {
	convey.Convey("Testing Successfully added services to workqueue after service event received", t, func() {
		fakeService := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: DEFAULT,
			},
		}
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)
		networkController.AddServiceEvent(fakeService)
		obj, shouldQuit := networkController.Workqueue.Get()
		assert.NotEqual(t, nil, obj, CONTROLLER_INIT)
		assert.Equal(t, false, shouldQuit, CONTROLLER_INIT)
	})
}

func TestAddServiceEventFailure(t *testing.T) {
	convey.Convey("Testing Failure to add services to workqueue after service event received", t, func() {
		fakeService := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{},
		}
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)

		patch1 := gomonkey.ApplyFunc(cache.MetaNamespaceKeyFunc, func(interface{}) (string, error) {
			return "", errors.New("Returning error")
		})
		defer patch1.Reset()
		networkController.AddServiceEvent(fakeService)
		assert.Equal(t, 0, networkController.Workqueue.Len(), WRK_QUEUE_EMPTY)

	})
}

func TestUpdateSvcEvent(t *testing.T) {
	convey.Convey("Testing Update services to workqueue after service event received, but no change in service", t, func() {
		fakeServiceOld := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				ResourceVersion: "1",
			},
		}
		fakeServiceNew := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				ResourceVersion: "1",
			},
		}
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)

		patch1 := gomonkey.ApplyFunc(cache.MetaNamespaceKeyFunc, func(interface{}) (string, error) {
			return "", errors.New("Returning error")
		})
		defer patch1.Reset()
		networkController.UpdateSvc(fakeServiceOld, fakeServiceNew)
		assert.Equal(t, 0, networkController.Workqueue.Len(), WRK_QUEUE_EMPTY)

	})
}

func TestUpdatePodEvent(t *testing.T) {
	convey.Convey("Testing Update services to workqueue after pod event received", t, func() {
		fakePodOld := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "fakePod1",
				Namespace:       "fakeNamespace1",
				ResourceVersion: "0",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "fakeContainer",
						Image: "fakeImage",
					},
				},
			},
		}
		fakePodNew := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:            "fakePod1",
				Namespace:       "fakeNamespace1",
				ResourceVersion: "1",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "fakeContainer",
						Image: "fakeImage",
					},
				},
			},
		}
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)

		patch1 := gomonkey.ApplyFunc(cache.MetaNamespaceKeyFunc, func(interface{}) (string, error) {
			return "", errors.New("Returning error")
		})
		defer patch1.Reset()
		//Case1: No annotation exist
		networkController.UpdatePod(fakePodOld, fakePodNew)
		assert.Equal(t, 0, networkController.Workqueue.Len(), WRK_QUEUE_EMPTY)

		//Case2: No service exist
		fakePodNew.Annotations = make(map[string]string)
		fakePodNew.Annotations["k8s.v1.cni.cncf.io/networks"] = "default/macvlan-conf1"
		networkController.UpdatePod(fakePodOld, fakePodNew)
		assert.Equal(t, 0, networkController.Workqueue.Len(), WRK_QUEUE_EMPTY)

		//Case2: service exist and added to workqueue
		fakeService := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				ResourceVersion: "1",
			},
		}
		var services []*corev1.Service
		patch2 := gomonkey.ApplyFunc(watcher.GetPodServices, func(corelisters.ServiceLister, *corev1.Pod) ([]*corev1.Service, error) {
			services = append(services, fakeService)
			return services, nil
		})
		defer patch2.Reset()
		patch3 := gomonkey.ApplyFunc(cache.MetaNamespaceKeyFunc, func(interface{}) (string, error) {
			return DEFAULT, nil
		})
		defer patch3.Reset()
		networkController.UpdatePod(fakePodOld, fakePodNew)
		obj, _ := networkController.Workqueue.Get()
		assert.NotEqual(t, nil, obj, "Workeue should not be empty")
	})
}

func TestHandleItemError(t *testing.T) {
	convey.Convey("Testing Handle workeue item", t, func() {
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)

		patch1 := gomonkey.ApplyFunc(cache.MetaNamespaceKeyFunc, func(interface{}) (string, error) {
			return "", errors.New("Returning error")
		})
		defer patch1.Reset()

		//Service not exist
		err := networkController.HandleItem(DEFAULT)
		assert.Equal(t, "service not found", err.Error(), "Service shouldn't exist")

		//Service exist with out network annotation
		fakeService1 := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				ResourceVersion: "1",
			},
		}

		var n1 *watcher.Controller
		patch2 := gomonkey.ApplyMethod(reflect.TypeOf(n1), "GetServices", func(*watcher.Controller, string, string) (*corev1.Service, error) {
			return fakeService1, nil
		})
		err = networkController.HandleItem(DEFAULT)
		assert.Equal(t, "no network annotations", err.Error(), "Service shouldn't exist")
		patch2.Reset()

		//Service exist with many network
		fakeService2 := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            "mepserver",
				Namespace:       DEFAULT,
				ResourceVersion: "1",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "default/macvlan-conf1, default/macvlan-conf2",
				},
			},
		}

		var n2 *watcher.Controller
		patch3 := gomonkey.ApplyMethod(reflect.TypeOf(n2), "GetServices", func(*watcher.Controller, string, string) (*corev1.Service, error) {
			return fakeService2, nil
		})
		defer patch3.Reset()

		err = networkController.HandleItem("mepserver")
		assert.Equal(t, "multiple network in service spec", err.Error(), "Service shouldn't have so many network")

	})
}

func TestHandleItemError2(t *testing.T) {
	convey.Convey("Testing Handle workeue item", t, func() {
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)

		patch1 := gomonkey.ApplyFunc(cache.MetaNamespaceKeyFunc, func(interface{}) (string, error) {
			return "", errors.New("Returning error")
		})
		defer patch1.Reset()

		//Service exist with out pod
		fakeService2 := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            "mepserver",
				Namespace:       DEFAULT,
				ResourceVersion: "1",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "default/macvlan-conf1",
				},
			},
		}

		var n2 *watcher.Controller
		patch2 := gomonkey.ApplyMethod(reflect.TypeOf(n2), "GetServices", func(*watcher.Controller, string, string) (*corev1.Service, error) {
			return fakeService2, nil
		})
		patch3 := gomonkey.ApplyMethod(reflect.TypeOf(n2), "GetPodsOfService", func(*watcher.Controller, *corev1.Service) ([]*corev1.Pod, error) {
			return nil, errors.New("No pod exist with matching selector")
		})
		err := networkController.HandleItem("mepserver")
		assert.Equal(t, "No pod exist with matching selector", err.Error(), "")
		defer patch2.Reset()
		patch3.Reset() //Reset patch3 here

		fakePod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fakePod1",
				Namespace: DEFAULT,
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks":        "macvlan-conf1",
					"k8s.v1.cni.cncf.io/networks-status": "[{\n        \"name\": \"openshift-sdn\",\n        \"ips\": [\n            \"10.131.0.10\"\n        ],\n        \"default\": true,\n        \"dns\": {}\n}]",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "fakeContainer",
						Image: "fakeImage",
					},
				},
			},
		}
		var n3 *watcher.Controller
		patch4 := gomonkey.ApplyMethod(reflect.TypeOf(n3), "GetPodsOfService", func(*watcher.Controller, *corev1.Service) ([]*corev1.Pod, error) {
			var podList []*corev1.Pod
			podList = append(podList, fakePod)
			return podList, nil
		})
		defer patch4.Reset()

		err = networkController.HandleItem("mepserver")
		assert.Equal(t, "no service endpoints found", err.Error(), "")
		defer patch2.Reset()
	})
}

func TestHandleItemWithServiceExist(t *testing.T) {
	convey.Convey("Testing Handle service items", t, func() {
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)

		//Service exist with a network
		fakePod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fakePod1",
				Namespace: DEFAULT,
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks":        "macvlan-conf1",
					"k8s.v1.cni.cncf.io/networks-status": "[{\n        \"name\": \"openshift-sdn\",\n        \"ips\": [\n            \"10.131.0.10\"\n        ],\n        \"default\": true,\n        \"dns\": {}\n}]",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "fakeContainer",
						Image: "fakeImage",
					},
				},
			},
		}
		fakeService2 := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            "mepserver",
				Namespace:       DEFAULT,
				ResourceVersion: "1",
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/networks": "default/macvlan-conf1",
				},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:     "myport",
						Port:     30000,
						Protocol: corev1.ProtocolTCP,
					},
				},
			},
		}
		fakeEp := &corev1.Endpoints{
			TypeMeta: metav1.TypeMeta{
				Kind:       ENDPOINTS,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "mepserver",
				Namespace: DEFAULT,
			},
			Subsets: []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							IP: IP1,
						},
					},
					Ports: []corev1.EndpointPort{
						{
							Name:     HTTP_PORT,
							Port:     30001,
							Protocol: corev1.ProtocolTCP,
						},
					},
				},
			},
		}
		var n1 *watcher.Controller
		patch1 := gomonkey.ApplyMethod(reflect.TypeOf(n1), "GetServices", func(*watcher.Controller, string, string) (*corev1.Service, error) {
			return fakeService2, nil
		})
		defer patch1.Reset()

		patch2 := gomonkey.ApplyMethod(reflect.TypeOf(n1), "GetPodsOfService", func(*watcher.Controller, *corev1.Service) ([]*corev1.Pod, error) {
			var podList []*corev1.Pod
			podList = append(podList, fakePod)
			return podList, nil
		})
		defer patch2.Reset()

		patch3 := gomonkey.ApplyMethod(reflect.TypeOf(n1), "GetEndpoints", func(*watcher.Controller, string, string) (*corev1.Endpoints, error) {
			return fakeEp, nil
		})
		defer patch3.Reset()

		patch4 := gomonkey.ApplyMethod(reflect.TypeOf(n1), "UpdateEndpoints", func(*watcher.Controller, *corev1.Endpoints) error {
			return nil
		})
		defer patch4.Reset()

		patch5 := gomonkey.ApplyFunc(watcher.IsInNetworkSelectionElementsArray, func(string, []*types.NetworkSelectionElement) bool {
			return true
		})
		defer patch5.Reset()

		err := networkController.HandleItem("mepserver")
		assert.Equal(t, nil, err, "Service shouldn't have so many network")

	})
}

func TestProcessNextWorkItemError(t *testing.T) {
	convey.Convey("Testing workeue item", t, func() {
		fakeService := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: DEFAULT,
			},
		}
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)
		networkController.AddServiceEvent(fakeService)
		shouldContinue := networkController.ProcessNextWorkItem()
		assert.Equal(t, true, shouldContinue, "Error processing workqueue")

		networkController.Workqueue.Add(1)
		shouldContinue = networkController.ProcessNextWorkItem()
		assert.Equal(t, true, shouldContinue, "Error processing workqueue")
	})
}

func TestProcessNextWorkItemSuccess(t *testing.T) {
	convey.Convey("Testing workeue item", t, func() {
		fakeService := &corev1.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       SERVICE,
				APIVersion: V1,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: DEFAULT,
			},
		}
		k8sclientSet := k8sfake.NewSimpleClientset()
		networkController := watcher.NewNetworkController(
			k8sclientSet)
		assert.NotEqual(t, nil, networkController.ServiceLister, CONTROLLER_INIT)
		networkController.AddServiceEvent(fakeService)
		var n1 *watcher.Controller
		patch1 := gomonkey.ApplyMethod(reflect.TypeOf(n1), "HandleItem", func(*watcher.Controller, string) error {
			return nil
		})
		defer patch1.Reset()
		shouldContinue := networkController.ProcessNextWorkItem()
		assert.Equal(t, true, shouldContinue, "Error processing workqueue")
	})
}
