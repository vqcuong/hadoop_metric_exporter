FROM golang:1.19.3-alpine

ENV HADOOP_EXPORTER_WORKDIR=/hadoop_exporter
ENV HADOOP_EXPORTER_METRICS_DIR=${HADOOP_EXPORTER_WORKDIR}/metrics
ENV HADOOP_EXPORTER_PORT=9123
ENV GOPATH=/hadoop_exporter/.go

WORKDIR ${HADOOP_EXPORTER_WORKDIR}

COPY ./hadoop_exporter/* ./
COPY ./entrypoint.sh /entrypoint.sh

RUN set -ex \
    && go mod download \
    && go build -o /usr/local/bin/hadoop_exporter \
    && chmod +x /entrypoint.sh

EXPOSE ${HADOOP_EXPORTER_PORT}

ENTRYPOINT ["/entrypoint.sh"]
