FROM gxf-base:latest

ARG VERSION=5.20.0

COPY startup.sh /home/
#COPY web-device-simulator.war /home/dev/Tools/tomcat8/webapps/
COPY web-device-simulator.properties /etc/osp/
COPY move-war-file.sh /home/

RUN /home/move-war-file.sh

CMD ["/home/startup.sh"]
EXPOSE 443
