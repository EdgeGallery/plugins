### Build

```shell
$ docker build -t swr.ap-southeast-1.myhuaweicloud.com/edgegallery/mqtt-go-device:latest .
```

### Run

```shell
$ docker run -it \
    -e BROKER_HOST=172.16.104.33 \
    -e BROKER_PORT=31883 \
    -e DEVICE_COUNT=1000 \
    -e TOPIC=Room1/Condition \
    swr.ap-southeast-1.myhuaweicloud.com/edgegallery/mqtt-go-device:latest
```

### K8s command

```shell
$ kubectl run iot-device -n test --image=swr.ap-southeast-1.myhuaweicloud.com/edgegallery/mqtt-go-device:latest 
--env="BROKER_HOST=mqtt" --env="BROKER_PORT=1883" --env="DEVICE_COUNT=10" --env="TOPIC=Room1/Condition"

$ kubectl delete -n test po iot-device
```