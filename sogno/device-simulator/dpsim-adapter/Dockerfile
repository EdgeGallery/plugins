FROM python:3.7

RUN mkdir -p /usr/app

ENV APP_HOME=/usr/app
ENV UID=166
ENV GID=166
ENV USER_NAME=eguser
ENV GROUP_NAME=eggroup

# Create an app user so our program doesn't run as root.
RUN apt-get -y update &&\
    groupadd -r -g $GID $GROUP_NAME &&\
    useradd -r -u $UID -g $GID -d $APP_HOME -s /sbin/nologin -c "Docker image user" $USER_NAME

COPY dpsim-plugin.py /home

RUN pip3 install paho-mqtt

CMD ["python3", "/home/dpsim-plugin.py"]
