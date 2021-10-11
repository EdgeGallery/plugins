FROM alpine
RUN apk update && apk add libc6-compat
WORKDIR /
COPY mqtt-simulator /mqtt-simulator
CMD ["sh", "-c", "/mqtt-simulator"]