FROM swr.cn-north-4.myhuaweicloud.com/eg-common/grafana/grafana:8.1.5-ubuntu
#RUN apk update && apk add libc6-compat
USER root
#COPY grafanaplugin/ grafanaplugin/
RUN apt-get -y upgrade && \
    apt-get -y update && \
    apt-get install -y git && git clone https://github.com/taosdata/grafanaplugin.git && \
    cp -r grafanaplugin/dist /var/lib/grafana/plugins/tdengine && \
    rm -rf grafanaplugin && \
    sed -i 's/;allow_loading_unsigned_plugins =/;allow_loading_unsigned_plugins = taosdata-tdengine-datasource/g' /etc/grafana/grafana.ini
USER grafana
ENTRYPOINT [ "/run.sh" ]