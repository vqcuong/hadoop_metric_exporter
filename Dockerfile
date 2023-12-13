FROM golang:1.19.3-alpine as builder

WORKDIR /hadoop_metric_exporter
COPY ./collector /hadoop_metric_exporter/
COPY ./server /hadoop_metric_exporter/
COPY ./utils /hadoop_metric_exporter/
COPY ./go.mod /hadoop_metric_exporter/
COPY ./go.sum /hadoop_metric_exporter/
COPY ./main.go /hadoop_metric_exporter/

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN set -ex \
  && go mod download \
  && go build

FROM debian:10-slim

ENV HADOOP_EXPORTER_METRICS_DIR=/etc/hadoop_metric_exporter/rules
ENV HADOOP_EXPORTER_PORT=9123

COPY ./rules ${HADOOP_EXPORTER_METRICS_DIR}
COPY ./entrypoint.sh /entrypoint.sh
COPY --from=builder /hadoop_metric_exporter/hadoop_metric_exporter /hadoop_metric_exporter

RUN set -ex \
  && chmod +x /entrypoint.sh

EXPOSE ${HADOOP_EXPORTER_PORT}

ENTRYPOINT ["/entrypoint.sh"]
