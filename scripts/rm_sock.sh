#!/bin/bash
set -ex
echo -e "\033[0;32m[rm sock]\033[0m"
# rm /tmp/kraken-agent-registry.sock
rm /tmp/kraken-proxy-registry.sock
rm /tmp/kraken-proxy-registry-override.sock
rm /tmp/kraken-origin.sock
rm /tmp/kraken-tracker.sock
rm /tmp/kraken-build-index.sock
