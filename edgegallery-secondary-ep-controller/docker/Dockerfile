# Copyright 2020 Huawei Technologies Co., Ltd.
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

FROM swr.cn-north-4.myhuaweicloud.com/eg-common/golang:1.14.2-alpine3.11 as builder
ENV HOME=/usr/src

RUN mkdir -p $HOME
COPY . $HOME
WORKDIR $HOME

RUN apk update &&\
    apk add --no-cache libc-dev gcc shadow &&\
    go mod download &&\
    GOOS=linux go build -buildmode=pie -ldflags '-linkmode "external" -extldflags "-static"' -o edgegallery-secondary-ep-controller

RUN chmod 550 $HOME/edgegallery-secondary-ep-controller

FROM swr.cn-north-4.myhuaweicloud.com/eg-common/alpine:latest
RUN sed -i "s|umask 022|umask 027|g" /etc/profile
# Create the home directory for the new mep user.
RUN mkdir -p /usr/mep
RUN mkdir -p /usr/mep/bin
# Set the home directory to our app user's home.
ENV HOME=/usr/mep
ENV UID=166
ENV GID=166
ENV USER_NAME=eguser
ENV GROUP_NAME=eggroup
ENV ENV="/etc/profile"
# Create an app user so our program doesn't run as root.
RUN apk update &&\
    apk add shadow &&\
    groupadd -r -g $GID $GROUP_NAME &&\
    useradd -r -u $UID -g $GID -d $HOME -s /sbin/nologin -c "Docker image user" $USER_NAME

WORKDIR $HOME
#RUN chmod 750 /usr/bin/edgegallery-secondary-ep-controller
RUN chmod 750 $HOME &&\
    chmod 550 -R $HOME/bin/ &&\
    mkdir -p -m 750 $HOME/log &&\
    chown -hR $USER_NAME:$GROUP_NAME $HOME

# Copy in the application exe.
COPY --from=builder --chown=$USER_NAME:$GROUP_NAME /usr/src/edgegallery-secondary-ep-controller/ $HOME/bin/

USER $USER_NAME

LABEL io.k8s.display-name="edgegallery-secondary-ep-controller"
RUN ls $HOME/bin
CMD ["sh", "-c", "$HOME/bin/edgegallery-secondary-ep-controller"]