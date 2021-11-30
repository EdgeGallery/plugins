FROM swr.cn-north-4.myhuaweicloud.com/eg-common/grafana/grafana:6.7.6
USER root
#ARG PLUGIN_VERSION=3.1.2
ENV GF_INSTALL_PLUGINS=https://github.com/taosdata/grafanaplugin/releases/download/v3.1.3/tdengine-datasource-3.1.3.zip;tdengine-datasource
USER grafana
ENTRYPOINT [ "/run.sh" ]