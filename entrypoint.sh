#!/bin/bash

if [[ "x$1" != "x" && "$1" != -* ]]; then
  exec "$@"
else
  /hadoop_exporter "$@"
fi
