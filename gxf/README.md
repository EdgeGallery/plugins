# Gxf


## Build
```shell
$ ./create-docker.sh
```



## Deploy
``` Note: add PGPASSWORD=1234 as env variable in gxf-profile.yaml and in gxf-demo.yaml```
1. Take gxf publiclight domain profile from installer repo. Deploy [gxf-publiclight-domain](https://gitee.com/edgegallery/installer/ProfileManagment/SampleGxfProfile/gxf-profile.yaml)
2. kubectl apply -f gxf-demo.yaml

```
NOTE: Go to https://grid-exchange-fabric.gitbook.io/gxf/userguide/installationguide/request/testosgpdemoapp page for demo instruction. 
Follow the instruction except Opening Device Simulator to Add a Device. 
Change only one field and fill IP Address: gxf-device-simulator
```

TO DO: Generate pem file inside an image.
