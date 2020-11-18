module edgegallery-secondary-ep-controller

go 1.14

require (
	github.com/agiledragon/gomonkey v2.0.1+incompatible
	github.com/containernetworking/cni v0.8.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/intel/multus-cni v0.0.0-20180818113950-86af6ab69fe2
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/k8snetworkplumbingwg/network-attachment-definition-client v1.1.0
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/net v0.0.0-20200813134508-3edf25e44fcc // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	k8s.io/api v0.18.4
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v0.18.4
	k8s.io/kubernetes v1.11.0-alpha.1.0.20180420161653-9c60fd5242c4
)

replace (
	github.com/containernetworking/cni v0.8.0 => github.com/containernetworking/cni v0.7.0-alpha1
	github.com/googleapis/gnostic v0.5.1 => github.com/googleapis/gnostic v0.4.0
)
