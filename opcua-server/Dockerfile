FROM swr.cn-north-4.myhuaweicloud.com/eg-common/python:3.8-slim
RUN mkdir -p /usr/app
ENV APP_HOME=/usr/app

# Set the home directory to our app user's home.
ENV APP_HOME=/usr/app
ENV UID=166
ENV GID=166
ENV USER_NAME=eguser
ENV GROUP_NAME=eggroup
ENV ENV="/etc/profile"
ENV PYTHONUNBUFFERED=0

# Create an app user so our program doesn't run as root.
RUN apt-get -y update &&\
    groupadd -r -g $GID $GROUP_NAME &&\
    useradd -r -u $UID -g $GID -d $APP_HOME -s /sbin/nologin -c "Docker image user" $USER_NAME

RUN apt install -y git

RUN apt-get install -y --fix-missing \
    build-essential \
    && apt-get clean && rm -rf /tmp/* /var/tmp/*

COPY requirements.txt $APP_HOME/

# Set the working directory.
WORKDIR $APP_HOME

RUN pip install opcua
RUN pip install -r requirements.txt

RUN git clone https://github.com/FreeOpcUa/python-opcua.git

WORKDIR $APP_HOME/python-opcua/examples/
CMD python3 ./server-minimal.py
