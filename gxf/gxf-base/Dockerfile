FROM ubuntu:16.04
  
RUN apt-get update && apt-get install -y sudo wget git psmisc &&\
    groupadd -r -g 166 dev &&\
    useradd -r -u 166 -g 166 -d /home/dev -s /sbin/nologin -c "dev user" dev

RUN mkdir -p /home/manifests
RUN mkdir -p /home/dev
RUN touch /home/dev/.profile

COPY manifests /home/manifests
#COPY eclipse-2020-03-R.tar.gz /home/dev/Downloads/osgp/
COPY moveFile.sh /home/

WORKDIR /home
RUN git clone https://github.com/OSGP/Config.git
RUN cd /home/Config && git checkout release-5.20.0
RUN rm -r /home/Config/puppet/manifests
RUN mv /home/manifests /home/Config/puppet/manifests
WORKDIR /home/Config/puppet

RUN ./install-puppet.sh
RUN ./run-puppet.sh
RUN puppet apply manifests/install-java.pp
RUN /home/moveFile.sh
RUN cd /home/dev/Sources/OSGP/open-smart-grid-platform && git checkout release-5.20.0
RUN cd /home/dev/Sources/OSGP/Config && git checkout release-5.20.0
RUN cd /home/dev/Sources/OSGP/Documentation && git checkout release-5.20.0

WORKDIR /home/dev/Sources/OSGP/open-smart-grid-platform
COPY web-demo-app.properties /home/dev/Sources/OSGP/open-smart-grid-platform/public-lighting-demo-app/web-demo-app/src/main/resources/
RUN mvn clean install

ENV JAVA_OPTS='-Xms512m -Xmx2048m -Xss512k -XX:MaxMetaspaceSize=1024m -XX:+CMSClassUnloadingEnabled -XX:+UseConcMarkSweepGC -Dcom.sun.management.jmxremote=true'
WORKDIR /home/dev/Tools/activemq/bin/linux-x86-64
