FROM debian:12-slim
MAINTAINER Filippov Alex <support@e154.ru>

ENV APP_DIR /opt/smart-home-node
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates iputils-ping
RUN update-ca-certificates
WORKDIR $APP_DIR
ADD . $APP_DIR
ENTRYPOINT (cd $APP_DIR && ./node)
USER nobody

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
