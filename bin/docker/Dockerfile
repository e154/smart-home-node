FROM debian
MAINTAINER Filippov Alex <af@suretly.com>

ENV APP_DIR /opt/smart-home-node
RUN sed 's/main$/main universe/' -i /etc/apt/sources.list \
    && apt-get update \
    && apt-get upgrade -y \
    && apt-get install ca-certificates -y \
    && mkdir -p $APP_DIR

WORKDIR $APP_DIR

ADD . $APP_DIR

ENTRYPOINT (cd $APP_DIR && ./node)

EXPOSE 3000


ENV NAME="node1"
ENV TOPIC="node1"
ENV MQTT_KEEP_ALIVE="300"
ENV MQTT_CONNECT_TIMEOUT="2"
ENV MQTT_SESSIONS_PROVIDER="mem"
ENV MQTT_TOPICS_PROVIDER="mem"
ENV MQTT_USERNAME="node1"
ENV MQTT_PASSWORD="node1"
ENV MQTT_IP="127.0.0.1"
ENV MQTT_PORT="1883"
ENV SERIAL="tty.usbserial,cu.wchusbserial,tty.SLAB_USB,ttyS,ttyUSB"

VOLUME $APP_DIR/conf
VOLUME $APP_DIR/data