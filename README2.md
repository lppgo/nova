[toc]

# 1: 使用 docker 搭建 devcluster

## 1.1 修改配置文件

### 1.1.1 Agent One docker 启动配置文件

```bash
# Define agent ports.
AGENT_REGISTRY_PORT=16000
AGENT_PEER_PORT=16001
AGENT_SERVER_PORT=16002

HOSTNAME=172.17.0.1

# Container config.
AGENT_CONTAINER_NAME=kraken-agent-one

# docker run
docker run -d \
-p 16000:16000 \
-p 16001:16001 \
-p 16002:16002 \
-v /mnt/e/sse/kraken/examples/devcluster/config/agent/development.yaml:/etc/kraken/config/agent/development.yaml \
--name kraken-agent-one kraken-agent:dev \
/usr/bin/kraken-agent \
--config=/etc/kraken/config/agent/development.yaml \
--peer-ip=172.17.0.1 \
--peer-port=16001 \
--agent-server-port=16002 \
--agent-registry-port=16000
```

### 1.1.2 Agent Two docker 启动配置文件

```bash
# Define agent ports.
AGENT_REGISTRY_PORT=17000
AGENT_PEER_PORT=17001
AGENT_SERVER_PORT=17002

HOSTNAME=172.17.0.1

# Container config.
AGENT_CONTAINER_NAME=kraken-agent-two

# docker run
docker run -d \
-p 17000:17000 \
-p 17001:17001 \
-p 17002:17002 \
-v /mnt/e/sse/kraken/examples/devcluster/config/agent/development.yaml:/etc/kraken/config/agent/development.yaml \
--name kraken-agent-two kraken-agent:dev \
/usr/bin/kraken-agent \
--config=/etc/kraken/config/agent/development.yaml \
--peer-ip=172.17.0.1 \
--agent-registry-port=17000 \
--peer-port=17001 \
--agent-server-port=17002
```

### 1.1.3 herd docker 启动配置文件

```bash
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

HOSTNAME=172.17.0.1

# docker run
docker run -d \
-p 14000:14000 \
-p 15000:15000 \
-p 15001:15001 \
-p 15002:15002 \
-p 15003:15003 \
-p 15004:15004 \
-p 15005:15005 \
-v /mnt/e/sse/kraken/examples/devcluster/herd_param.sh:/etc/kraken/herd_param.sh \
-v /mnt/e/sse/kraken/examples/devcluster/config/push-cli/development.yaml:/etc/kraken/config/push-cli/development.yaml \
-v /mnt/e/sse/kraken/examples/devcluster/config/origin/development.yaml:/etc/kraken/config/origin/development.yaml \
-v /mnt/e/sse/kraken/examples/devcluster/config/tracker/development.yaml:/etc/kraken/config/tracker/development.yaml \
-v /mnt/e/sse/kraken/examples/devcluster/config/build-index/development.yaml:/etc/kraken/config/build-index/development.yaml \
-v /mnt/e/sse/kraken/examples/devcluster/herd_start_processes.sh:/etc/kraken/herd_start_processes.sh \
--name kraken-herd kraken-herd:dev \
./herd_start_processes.sh
```

# 2: 使用 binary

## 2.1 依赖组件

### 2.1.1 tls 加密

```bash
# 1. Create a self-signed certificate for server:
# 1.1  create private key (生成密钥xxx.key)
openssl genrsa -aes256 -out server.key 4096
# 1.2 create a sign request (根据密钥生成证书请求文件xxx.csr)
openssl req -new -key server.key -out server.csr
# 1.3 Generate cert （根据密钥xxx.key和证书请求文件xxx.scr 生成crt证书）
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

# 2. Create an intermediate certificate for client:
openssl genrsa -aes256 -out client.key 4096
openssl req -new -key client.key -out client.csr

# 3. Verify
openssl verify -verbose -CAfile server.crt client.crt

# 4. Decrypt client key (because curl does not support encrypted key)
openssl rsa -in client.key -out client_decrypted.key

# 5. Both client and server should enforce verification
InsecureSkipVerify should be false in client and ClientAuth should be equal to tls.RequireAndVerifyClientCert in tls.Config.
In nginx config, ssl_verify_client should be on.

# 2: nginx
# /mnt/e/sse/kraken/docker
./setup_nginx.sh

# 3: redis

```

### 2.1.2 redis

- 注意 port

### 2.1.3 nginx

```bash
/usr/sbin/nginx
/usr/lib/nginx
/etc/nginx
/usr/share/nginx
```

## 2.2 config 各组件配置文件

- 配置文件

```bash
# mkdir -p /etc/kraken/config

base.yaml
development.yaml

# 注意各参数
```

```bash

# 1: config
mkdir -p /etc/kraken/config/agent

# 2: run
./kraken-agent  --config=/etc/kraken/config/agent/base.yaml \
--peer-ip=localhost \
--peer-port=16001 \
--agent-server-port=16002 \
--agent-registry-port=16000


```

## 2.3 kraken 组件启动

### 2.3.1 remove 对应的 unix \*\*\*.sock

```bash
#!/bin/bash

set -x
# echo `rm kraken-push-cli-registry.sock`
rm /tmp/kraken-agent-registry.sock
rm /tmp/kraken-push-cli-registry.sock
rm /tmp/kraken-push-cli-registry-override.sock
rm /tmp/kraken-origin.sock
rm /tmp/kraken-tracker.sock
rm /tmp/kraken-build-index.sock
```

### 2.3.2 查看占用 port 并 kill

```bash
sudo netstat -nltp | grep 1500
```

### 2.3.3 kraken-各组件

#### 2.3.3.1 kraken-agent

- kraken-agent run

```bash
./kraken-agent --config=/etc/kraken/config/agent/base.yaml \
--peer-ip=localhost \
--peer-port=16001 \
--agent-server-port=16002 \
--agent-registry-port=16000
```

- kraken-agent flags

```go
// Flags defines agent CLI flags.
type Flags struct {
	PeerIP            string // ("peer-ip", "", "ip which peer will announce itself as")
	PeerPort          int    // ("peer-port", 0, "port which peer will announce itself as")
	AgentServerPort   int    // ("agent-server-port", 0, "port which agent server listens on")
	AgentRegistryPort int    // ("agent-registry-port", 0, "port which agent registry listens on")
	ConfigFile        string // ("config", "", "configuration file path")
	Zone              string // ("zone", "", "zone/datacenter name")
	KrakenCluster     string // ("cluster", "", "cluster name (e.g. prod01-zone1)")
	SecretsFile       string // ("secrets", "", "path to a secrets YAML file to load into configuration")
}
```

#### 2.3.3.2 kraken-testfs

```bash
./testfs --port=14000
```

#### 2.3.3.3 kraken-proxy

```bash
./kraken-push-cli --config=/etc/kraken/config/push-cli/base.yaml --port=15000 --server-port=15005
```

#### 2.3.3.4 kraken-origin

```bash
./kraken-origin \
    --config=/etc/kraken/config/origin/base.yaml \
    --blobserver-hostname=localhost \
    --blobserver-port=15002 \
    --peer-ip=localhost \
    --peer-port=15001
```

#### 2.3.3.5 kraken-tracker

```bash
./kraken-tracker --config=/etc/kraken/config/tracker/base.yaml --port=15003
```

#### 2.3.3.6 kraken-build-index

```bash
./kraken-build-index --config=/etc/kraken/config/build-index/base.yaml  --port=15004
```

### 2.3.4 shell 脚本

### 2.3.4.1 kraken-agent

```bash
./kraken-agent --config=/etc/kraken/config/agent/base.yaml \
--peer-ip=localhost \
--peer-port=16001 \
--agent-server-port=16002 \
--agent-registry-port=16000
```

### 2.3.4.1 kraken-herd

```bash
#!/bin/bash

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
# HOSTNAME=172.17.0.1


# redis-server --port ${REDIS_PORT} &

sleep 1

/usr/local/bin/kraken-testfs \
    --port=${TESTFS_PORT} \
    &> /var/log/kraken/kraken-testfs/stdout.log &

/usr/local/bin/kraken-push-cli \
    --config=/etc/kraken/config/push-cli/base.yaml \
    --port=${PROXY_PORT} \
    --server-port=${PROXY_SERVER_PORT} \
    &> /var/log/kraken/kraken-push-cli/stdout.log &

/usr/local/bin/kraken-origin \
    --config=/etc/kraken/config/origin/base.yaml \
    --blobserver-hostname=${HOSTNAME} \
    --blobserver-port=${ORIGIN_SERVER_PORT} \
    --peer-ip=${HOSTNAME} \
    --peer-port=${ORIGIN_PEER_PORT} \
    &> /var/log/kraken/kraken-origin/stdout.log &

/usr/local/bin/kraken-tracker \
    --config=/etc/kraken/config/tracker/base.yaml \
    --port=${TRACKER_PORT} \
    &> /var/log/kraken/kraken-tracker/stdout.log &

/usr/local/bin/kraken-build-index \
    --config=/etc/kraken/config/build-index/base.yaml \
    --port=${BUILD_INDEX_PORT} \
    &> /var/log/kraken/kraken-build-index/stdout.log &

sleep 1

# Poor man's supervisor.
while : ; do
    # for c in redis-server kraken-testfs kraken-origin kraken-tracker kraken-build-index kraken-push-cli; do
    for c in kraken-testfs kraken-origin kraken-tracker kraken-build-index kraken-push-cli; do
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
```

# 4: push/pull

以 `portainer/portainer:latest` images 为例：

```go
// tag
docker tag portainer/portainer:latest localhost:15000/portainer:latest

// push
docker push localhost:15000/portainer:latest

```

# 5：2021-0601 会议

1. 本周摸透所有配置
2. ftp 存储接口
3. pull agent

# 6：2021-0607 会议


1. proxy upload

BPC client
Proxy 新加一个 server，暴露 origin upload 相关 api

2. origin
   s3 backend ---<ftp>

3. agent
   cli download

# 7: 实现核心上传下载功能，文档参照README3.md
## 7.1 修改源代码
## 7.2 项目部署
## 7.3 上传下载实现
## 7.4 根据业务拓展，实现业务需求
## 7.5 项目优化






