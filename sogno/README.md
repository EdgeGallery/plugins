# Sogno 
## Architecture
![](https://sogno-platform.github.io/docs/examples/state-estimation/state-estimation.svg) 
---
## Deployment Steps
```Note: Kubernetes and helm are required. Huge pages is required if one's using device simulator. ```

```
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

# Restart k3s service to apply changes
sudo systemctl restart k3s    # If you are using k8s then restart k8s

# Ensure the KUBECONFIG env is still set correctly
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
```

1. Deploy sogno profile from installer repo. Take the yaml file from [here](https://gitee.com/edgegallery/installer/ProfileManagment/SampleSognoProfile/sogno-profile-deploy.yaml)
2. Deploy device simulator using [this](device-simulator/dpsim/README.md) README file. Huge pages is required for device simulator. Once you enable Huge page then restart k8s. Not deploy it in edge gallery environment.
3. Deploy dp-sim adapter using [this](device-simulator/dpsim-adapter/dpsim-plugin.yaml) yaml file.
4. Deploy application using [README](application/README.md) file. Helm is required to install these application.
5. Check the grafana using https://localhost:31230 . Make sure that every thing is working.
NOTE: grafana username is demo and password is demo.