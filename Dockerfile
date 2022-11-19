FROM golang:1.19.3-alpine as builder

WORKDIR /hadoop_exporter
COPY ./hadoop_exporter /hadoop_exporter

RUN set -ex \
    && go mod download \
    && env GOOS=linux GOARCH=amd64 go build

FROM debian:10-slim

ENV HADOOP_EXPORTER_METRICS_DIR=/etc/hadoop_exporter/rules
ENV HADOOP_EXPORTER_PORT=9123

COPY ./rules ${HADOOP_EXPORTER_METRICS_DIR}
COPY ./entrypoint.sh /entrypoint.sh
COPY --from=builder /hadoop_exporter/hadoop_exporter /hadoop_exporter

RUN set -ex \
    && chmod +x /entrypoint.sh

EXPOSE ${HADOOP_EXPORTER_PORT}

ENTRYPOINT ["/entrypoint.sh"]
