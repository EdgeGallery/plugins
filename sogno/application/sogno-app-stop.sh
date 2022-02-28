#! /bin/bash

echo -e "\n************* Sogno application is stoping ************\n"

sudo kubectl delete -f example-deployments/pyvolt-dpsim-demo/visualization/dashboard-configmap.yaml
sudo rm -rf example-deployments

sudo helm delete grafana
sudo helm delete telegraf
sudo helm delete influxdb

sudo helm repo remove sogno
sudo helm repo remove bitnami
sudo helm repo remove influxdata
sudo helm repo remove grafana

echo -e "\n************* Sogno application stop successfull ************\n"
