# Sogno 
## Architecture
![](https://sogno-platform.github.io/docs/examples/state-estimation/state-estimation.svg) 
---
## Deployment Steps
```Note: Kubernetes and helm are required. Huge pages is required if one's using device simulator. ```

```
1. Deploy sogno profile from installer repo. 
   Take the yaml file from [here](https://gitee.com/edgegallery/installer/ProfileManagment/SampleSognoProfile/sogno-profile-deploy.yaml)

2. start Device simulator:
Prerequisite: it need 1024 huge pages enabled in k8s cluster
# Verify HugePages
cat /proc/meminfo | grep Huge

AnonHugePages:    104448 kB
ShmemHugePages:        0 kB
FileHugePages:         0 kB
HugePages_Total:       0		# we require a minimum of 1024
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
Hugetlb:               0 kB

# Increase No of HPgs
echo 1024 | sudo tee /proc/sys/vm/nr_hugepages

# Check it worked
cat /proc/meminfo | grep Huge

AnonHugePages:    104448 kB
ShmemHugePages:        0 kB
FileHugePages:         0 kB
HugePages_Total:    1024
HugePages_Free:     1024
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
Hugetlb:         2097152 kB

If you don't see 1024 next to HugePages_Total, you may need to restart
your system and try again with a fresh boot.

# Restart k8s service to apply changes
sudo systemctl restart kubelet
```


Now Deploy device simulator using [this](device-simulator/dpsim/README.md) README file. 
 Huge pages is required for device simulator.  
Deploy dp-sim adapter using [this](device-simulator/dpsim-adapter/dpsim-plugin.yaml) yaml file.

Next Deploy Application
Deploy application using [README](application/README.md) file. Helm is required to install these application.
 Check the grafana using https://localhost:31230 . Make sure that every thing is working.
 NOTE: grafana username is demo and password is demo.

Some common troubleshoot for k8s enviorment:
For application deployment:
- if rabbitMq pod failed with below error:
/opt/bitnami/scripts/librabbitmq.sh: line 750: ulimit: open files: cannot modify limit: Operation not permitted

Then make change in helm chart for below enviorment value:
make this env as empty string
RABBITMQ_ULIMIT_NOFILES ""

- If influx DB failed with PVC not able to bound 
  make exiting storage class as defoult and delete pvc.
  then restart helm

- 
