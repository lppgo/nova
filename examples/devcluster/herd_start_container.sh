#!/bin/bash

set -ex
# --------------------------------------------
# source examples/devcluster/herd_param.sh
# Define optional storage ports.
TESTFS_PORT=14000
REDIS_PORT=14001

# Define herd ports.
PROXY_PORT=15000
ORIGIN_PEER_PORT=15001
ORIGIN_SERVER_PORT=15002
TRACKER_PORT=15003
BUILD_INDEX_PORT=15004
PROXY_SERVER_PORT=15005



# docker0
HOSTNAME=172.17.0.1
# HOSTNAME=host.docker.internal

# --------------------------------------------

# Start kraken herd.
docker run -d \
    -p ${TESTFS_PORT}:${TESTFS_PORT} \
    -p ${ORIGIN_SERVER_PORT}:${ORIGIN_SERVER_PORT} \
    -p ${ORIGIN_PEER_PORT}:${ORIGIN_PEER_PORT} \
    -p ${TRACKER_PORT}:${TRACKER_PORT} \
    -p ${BUILD_INDEX_PORT}:${BUILD_INDEX_PORT} \
    -p ${PROXY_PORT}:${PROXY_PORT} \
    -p ${PROXY_SERVER_PORT}:${PROXY_SERVER_PORT} \
    -v $(pwd)/examples/devcluster/config/origin/development.yaml:/etc/kraken/config/origin/development.yaml \
    -v $(pwd)/examples/devcluster/config/tracker/development.yaml:/etc/kraken/config/tracker/development.yaml \
    -v $(pwd)/examples/devcluster/config/build-index/development.yaml:/etc/kraken/config/build-index/development.yaml \
    -v $(pwd)/examples/devcluster/config/proxy/development.yaml:/etc/kraken/config/proxy/development.yaml \
    -v $(pwd)/examples/devcluster/herd_param.sh:/etc/kraken/herd_param.sh \
    -v $(pwd)/examples/devcluster/herd_start_processes.sh:/etc/kraken/herd_start_processes.sh \
    --name kraken-herd \
    kraken-herd:dev ./herd_start_processes.sh
