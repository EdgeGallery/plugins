# DP sim
## Deployment Steps
```
Note: huge page is required
```

1. git clone https://github.com/sogno-platform/helm-charts
2. cd helm-charts/charts/
3. Go to dpsim-demo/files/Shmem_CIGRE_MV.conf and change following
```
host = "mqtt" # service name of emqx broker
port = 1883 # internal port of emqx broker
username = "admin"  # username of emqx broker
password = "****"   # password of emqx broker
out = {
       publish = "/fledge"
}
```
4. helm install dpsim-demo ./dpsim-demo/ 