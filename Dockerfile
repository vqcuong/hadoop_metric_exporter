FROM golang:1.19.3-alpine as builder

WORKDIR /go/app/hadoop_exporter

ENV HADOOP_EXPORTER_METRICS_DIR=/hadoop_exporter/rules
ENV HADOOP_EXPORTER_PORT=9123

COPY ./rules /hadoop_exporter/
COPY ./hadoop_exporter/go.mod ./
COPY ./hadoop_exporter/go.sum ./
RUN go mod download

COPY ./hadoop_exporter/* ./
COPY ./entrypoint.sh /entrypoint.sh

RUN set -ex \
    && go build -o /usr/local/bin/hadoop_exporter ./ \
    && chmod +x /entrypoint.sh

EXPOSE ${HADOOP_EXPORTER_PORT}

ENTRYPOINT ["/entrypoint.sh"]
