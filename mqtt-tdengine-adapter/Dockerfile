FROM swr.cn-north-4.myhuaweicloud.com/eg-common/tdengine/tdengine:2.2.0.2
#RUN apk update && apk add libc6-compat
WORKDIR /
COPY mqtt-engine-adapter /mqtt-engine-adapter
CMD ["sh", "-c", "/mqtt-engine-adapter"]