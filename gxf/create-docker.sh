cd gxf-base
sudo docker build -t gxf-base:latest .
cd ../gxf-publiclight
sudo docker build -t gxf-publiclight:latest .
cd ../gxf-publiclight-domain
sudo docker build -t gxf-publiclight-domain:latest .
cd ../gxf-device-simulator
sudo docker build -t gxf-device-simulator:latest .
