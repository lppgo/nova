#!/bin/bash

echo -e "\033[0;32m[export port]\033[0m"
# Define optional storage ports.
TESTFS_PORT=14000
# REDIS_PORT=14001
REDIS_PORT=6379

# Define herd ports.
PROXY_PORT=15000
ORIGIN_PEER_PORT=15001
ORIGIN_SERVER_PORT=15002
TRACKER_PORT=15003
BUILD_INDEX_PORT=15004
PROXY_SERVER_PORT=15005
HOSTNAME=localhost


echo -e "\033[0;32m[rm sock]\033[0m"
# rm /tmp/kraken-agent-registry.sock
rm /tmp/kraken-proxy-registry.sock
rm /tmp/kraken-proxy-registry-override.sock
rm /tmp/kraken-origin.sock
rm /tmp/kraken-tracker.sock
rm /tmp/kraken-build-index.sock

# echo -e "\033[0;32m[start redis]\033[0m"
# redis-server --port ${REDIS_PORT} &


echo -e "\033[0;32m[run testfs]\033[0m"
# ./testfs --port=14000
./testfs --port=${TESTFS_PORT} \
    &> /var/log/kraken/kraken-testfs/stdout.log &

echo -e "\033[0;32m[run origin]\033[0m"
#./origin \
#    --config=/etc/kraken/config/origin/base.yaml \
#    --blobserver-hostname=localhost \
#    --blobserver-port=15002 \
#    --peer-ip=localhost \
#    --peer-port=15001

./origin \
    --config=/etc/kraken/config/origin/base.yaml \
    --blobserver-hostname=${HOSTNAME} \
    --blobserver-port=${ORIGIN_SERVER_PORT} \
    --peer-ip=${HOSTNAME} \
    --peer-port=${ORIGIN_PEER_PORT} \
    &> /var/log/kraken/kraken-origin/stdout.log &

echo -e "\033[0;32m[run tracker]\033[0m"
#./tracker \
#    --config=/etc/kraken/config/tracker/base.yaml \
#    --port=15003
./tracker \
    --config=/etc/kraken/config/tracker/base.yaml \
    --port=${TRACKER_PORT} \
    &> /var/log/kraken/kraken-tracker/stdout.log &

echo -e "\033[0;32m[run build-index]\033[0m"
#./build-index \
#    --config=/etc/kraken/config/build-index/base.yaml \
#    --port=15004
./build-index \
    --config=/etc/kraken/config/build-index/base.yaml \
    --port=${BUILD_INDEX_PORT} \
    &> /var/log/kraken/kraken-build-index/stdout.log &

sleep 1

# Poor man's supervisor.
while : ; do
    # for c in redis-server kraken-testfs kraken-origin kraken-tracker kraken-build-index kraken-push-cli; do
    for c in kraken-testfs kraken-origin kraken-tracker kraken-build-index kraken-proxy; do
        ps aux | grep $c | grep -q -v grep
        status=$?
        if [ $status -ne 0 ]; then
            echo "$c exited unexpectedly. Logs:"
            tail -100 /var/log/kraken/$c/stdout.log
            exit 1
        fi
    done
    sleep 10
done
