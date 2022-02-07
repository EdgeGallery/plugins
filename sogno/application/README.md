# Application
## Deployment Steps
```
Note: Kubernetes and helm are required 
```
1. git clone https://github.com/sogno-platform/example-deployments.git
2. cd example-deployments/pyvolt-dpsim-demo
3. helm repo add sogno https://sogno-platform.github.io/helm-charts
3. helm repo add bitnami https://charts.bitnami.com/bitnami
4. helm repo add influxdata https://influxdata.github.io/helm-charts
5. helm repo add grafana https://grafana.github.io/helm-charts
6. helm repo update
7. helm install influxdb influxdata/influxdb -f database/influxdb-helm-values.yaml
8. helm install telegraf influxdata/telegraf -f ts-adapter/telegraf-values.yaml
9. helm install grafana grafana/grafana -f visualization/grafana_values.yaml
10. kubectl apply -f visualization/dashboard-configmap.yaml