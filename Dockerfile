FROM alpine
COPY ./bin/blackflowint /opt/blackflowint/
# RUN apk --no-cache add ca-certificates && update-ca-certificates
VOLUME ["/etc/blackflowint","/var/lib/blackflowint"]
ENV ZM_INT_LOGLEVEL="info"
ENV ZM_INT_STORAGELOCATION="/var/lib/blackflowint"
ENV ZM_INT_ADMINRESTAPIBINDADDRES=":5015"
ENV ZM_INT_ACTIVE_INTEGRATIONS = "restmqttproxy"
EXPOSE 5015
WORKDIR /opt/blackflowint/
ENTRYPOINT ["/opt/blackflowint/blackflowint"]