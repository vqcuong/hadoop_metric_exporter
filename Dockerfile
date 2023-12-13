FROM golang:1.21.4-alpine as builder

WORKDIR /hadoop_metric_exporter
COPY . /hadoop_metric_exporter/

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN set -ex \
  && go mod download \
  && go build

FROM debian:10-slim

ENV HADOOP_METRIC_EXPORTER_RULES_DIR=/exporter/rules
ENV HADOOP_METRIC_EXPORTER_PORT=9123

COPY ./rules ${HADOOP_METRIC_EXPORTER_RULES_DIR}
COPY ./entrypoint.sh /entrypoint.sh
COPY --from=builder /hadoop_metric_exporter/hadoop_metric_exporter /hadoop_metric_exporter

RUN set -ex \
  && chmod +x /entrypoint.sh

EXPOSE ${HADOOP_METRIC_EXPORTER_PORT}

ENTRYPOINT ["/entrypoint.sh"]
