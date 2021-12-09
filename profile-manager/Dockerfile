# Copyright 2021 Huawei Technologies Co., Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM swr.cn-north-4.myhuaweicloud.com/eg-common/golang:1.14.2 as builder
RUN apt update &&\
    apt install libc-dev gcc
ENV HOME=/
WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . $HOME
WORKDIR $HOME
# RUN GOOS=linux go build -buildmode=pie -ldflags '-linkmode "external" -extldflags "-static"' -o app
RUN GOOS=linux go build -o app
RUN chmod -R 550 $HOME/app


FROM swr.cn-north-4.myhuaweicloud.com/eg-common/alpine:latest
RUN apk update && apk add libc6-compat
WORKDIR /
RUN mkdir conf
COPY --from=builder $HOME/app $HOME
COPY --from=builder $HOME/conf/config.yaml $HOME/conf/
COPY --from=builder $HOME/conf/profile.yaml $HOME/conf/
RUN ls /app
CMD ["sh", "-c", "/app"]