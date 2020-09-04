FROM golang:alpine as builder
ENV GOPROXY https://proxy.golang.org
ENV GO111MODULE on
ENV HOME=/usr/src

RUN apk update &&\
    apk add --no-cache libc-dev gcc shadow
RUN mkdir -p $HOME
COPY . $HOME
WORKDIR $HOME
RUN go mod vendor
RUN go mod download
RUN GOOS=linux go build -buildmode=pie -ldflags '-linkmode "external" -extldflags "-static"' -o edgegallery-secondary-ep-controller


FROM alpine:latest
RUN apk update &&\
    apk add shadow

COPY --from=builder /usr/src/edgegallery-secondary-ep-controller/ /usr/bin/

WORKDIR /
RUN chmod 750 /usr/bin/edgegallery-secondary-ep-controller
LABEL io.k8s.display-name="Network Attachment Definitions Controller"

ENTRYPOINT ["/usr/bin/edgegallery-secondary-ep-controller"]