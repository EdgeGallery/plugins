#! /bin/bash

echo -e "\n************* Sogno application is starting ************\n"
sudo git clone https://github.com/sogno-platform/example-deployments.git
cd example-deployments/pyvolt-dpsim-demo
sudo helm repo add sogno https://sogno-platform.github.io/helm-charts
sudo helm repo add bitnami https://charts.bitnami.com/bitnami
sudo helm repo add influxdata https://influxdata.github.io/helm-charts
sudo helm repo add grafana https://grafana.github.io/helm-charts
sudo helm repo update
sudo helm install influxdb influxdata/influxdb -f database/influxdb-helm-values.yaml
sudo helm install telegraf influxdata/telegraf -f ts-adapter/telegraf-values.yaml
sudo helm install grafana grafana/grafana -f visualization/grafana_values.yaml
sudo kubectl apply -f visualization/dashboard-configmap.yaml
sleep 10s
cd ../..


echo -e "\n************* Sogno application started ************\n"
