FROM golang:1.19.3-alpine as builder

ENV HADOOP_EXPORTER_WORKDIR=/hadoop_exporter
ENV HADOOP_EXPORTER_METRICS_DIR=${HADOOP_EXPORTER_WORKDIR}/rules
ENV HADOOP_EXPORTER_PORT=9123

WORKDIR ${HADOOP_EXPORTER_WORKDIR}

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
