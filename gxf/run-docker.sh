sudo docker network create gxf-net
sudo docker run -d -p 8083:443 --net gxf-net --name gxf-device-simulator gxf-device-simulator
sudo docker run -d --net gxf-net --name gxf-publiclight-domain gxf-publiclight-domain
sudo docker run -d -p 8081:443 --net gxf-net --name gxf-publiclight gxf-publiclight
