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

FROM swr.cn-north-4.myhuaweicloud.com/eg-common/tdengine/tdengine:2.2.0.2 as builder
COPY --from=swr.cn-north-4.myhuaweicloud.com/eg-common/golang:1.14.2 /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
RUN apt-get update && apt-get install -y ca-certificates openssl gcc cmake build-essential
ENV HOME=/
ENV GOPROXY=https://goproxy.cn
WORKDIR /go/cache
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . $HOME
WORKDIR $HOME
#RUN GOOS=linux go build -buildmode=pie -ldflags '-linkmode "external" -extldflags "-static"' -o app
RUN GOOS=linux go build -o app
RUN chmod -R 550 $HOME/app


FROM swr.cn-north-4.myhuaweicloud.com/eg-common/tdengine/tdengine:2.2.0.2
ENV HOME=/
COPY --from=builder $HOME/app $HOME
#RUN apk update && apk add libc6-compat
WORKDIR $HOME
CMD ["sh", "-c", "/app"]
