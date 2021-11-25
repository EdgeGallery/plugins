FROM swr.cn-north-4.myhuaweicloud.com/eg-common/alpine:latest
RUN apk update && apk add libc6-compat
WORKDIR /
COPY mqtt-simulator /mqtt-simulator
CMD ["sh", "-c", "/mqtt-simulator"]