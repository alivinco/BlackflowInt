FROM alpine
COPY ./bin/blackflowint /opt/blackflowint/
# RUN apk --no-cache add ca-certificates && update-ca-certificates
VOLUME ["/etc/blackflowint","/var/lib/blackflowint"]
ENV ZM_LOGLEVEL="info"
ENV ZM_STORAGELOCATION="/var/lib/blackflowint"
ENV ZM_ADMINRESTAPIBINDADDRES=":5015"
EXPOSE 5015
WORKDIR /opt/blackflowint/
ENTRYPOINT ["/opt/blackflowint/blackflowint"]