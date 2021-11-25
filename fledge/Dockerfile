FROM swr.cn-north-4.myhuaweicloud.com/eg-common/jrei/systemd-ubuntu:18.04
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get -y upgrade && \
    apt-get -y update && \
    apt-get -y install avahi-daemon curl && \
    apt-get -y install cmake g++ make build-essential autoconf automake uuid-dev && \
    apt-get -y install libtool libboost-dev libboost-system-dev libboost-thread-dev libpq-dev libssl-dev libz-dev && \
    apt-get -y install python-dev python3-dev python3-pip postgresql sqlite3 libsqlite3-dev git pkg-config && \
    apt-get -y install libcurl4-openssl-dev wget apt-utils gnupg2 software-properties-common sudo && \
    rm -rf /var/lib/apt/lists/*
RUN git clone https://github.com/fledge-iot/fledge.git && \
    cd fledge && \
    make install && \
    rm -rf /fledge
# RUN git clone https://github.com/fledge-iot/fledge-south-mqtt.git && \
#     cd fledge-south-mqtt && \
#     python3 -m pip install -r python/requirements-mqtt-readings.txt && \
#     sed -i 's/localhost/mqtt-broker/g' /fledge-south-mqtt/mqtt-pub/__main__.py
# RUN apt-get -y install wget apt-utils gnupg2 software-properties-common sudo 
RUN apt-get autoclean && apt-get autoremove -y && apt-get update
RUN     wget -q -O - http://archives.fledge-iot.org/KEY.gpg | apt-key add - && \
    add-apt-repository "deb http://archives.fledge-iot.org/latest/ubuntu1804/x86_64/ / "
RUN apt-get -y update && \
    apt-get -yqq install krb5-user libpam-krb5 && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /
COPY fledge.sh /fledge.sh
COPY mqtt_north/ /usr/local/fledge/python/fledge/plugins/north/mqtt_north/
RUN python3 -m pip install -r /usr/local/fledge/python/fledge/plugins/north/mqtt_north/requirements.txt
# RUN apt-get -y install rsyslog
CMD /fledge.sh