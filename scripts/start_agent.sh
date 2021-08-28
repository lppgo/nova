#!/bin/bash

# echo -e "\033[0;32m[export port]\033[0m"
# PeerIP=localhost
# PeerPort=16001
# AgentServerPort=16002
# AgentRegistryPort=16000

# sleep 1

echo -e "\033[0;32m[run kraken-agent]\033[0m"
/usr/local/bin/kraken-agent --config=/etc/kraken/config/agent/base.yaml \
    --peer-ip=localhost \
    --peer-port=16001 \
    --agent-server-port=16002 \
    --agent-registry-port=16000
# &>/var/log/kraken/kraken-agent/stdout.log &
