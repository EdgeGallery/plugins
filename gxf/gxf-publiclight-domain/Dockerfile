FROM gxf-base:latest

ARG CERT=/home/dev/Sources/OSGP/Config/certificates/osgp-ca/private
ARG PVT=/home/dev/Sources/OSGP/Config/certificates/osgp-ca/certs
ARG VERSION=5.20.0

RUN mkdir -p /home/dev/Tools/tomcat8/tmp/

COPY startup.sh /home/
COPY gxf-publiclight-domain.key.pem $PVT/
COPY gxf-publiclight-domain.cert.pem  $CERT/
COPY vhost.conf /home/dev/Sources/OSGP/Config/apache-httpd/
COPY move-war-file.sh /home/

RUN /home/move-war-file.sh
RUN ln -sf $PVT/gxf-publiclight-domain.key.pem /etc/ssl/private
RUN ln -sf $CERT/gxf-publiclight-domain.cert.pem /etc/ssl/certs
RUN mv /home/dev/Tools/tomcat8/tmp/osgp-core.war /home/dev/Tools/tomcat8/webapps/osgp-core.war

CMD ["/home/startup.sh"]
