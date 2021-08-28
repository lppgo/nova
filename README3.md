[toc]

# P2P 文件分发系统(nova) 说明文档

## 1: 项目介绍

### 1.1: 项目背景

本系统是实现一个通过 BitTorrent 协议(P2P 方式)，进行分布式，快速进行文件分发的系统。

### 1.2: 项目参考

本项目是根据 Uber 的[kraken](https://github.com/uber/kraken) 进行二次开发实现的。

## 2: 项目部署

### 2.1 源项目 kraken 代码分析

- 查看 kraken 源码:https://github.com/uber/kraken
- 查看 kraken 源码分析文档`./docs/需规和设计文档/004_参考数据中心-文件分发系统设计说明书 - 阶段1-现有P2P和Docker逻辑梳理.md`

### 2.2 Deploy

#### **Nginx:**

- 目前该项目启动需要依赖 nginx，对内部服务做反向代理
- 多个进程(agent,proxy,origin,tracker,build-index)启动的时候都会根据 template 生成各自的 nginx 配置文件
- proxy 组件已经修改成 push 命令行，对应的 nginx 依赖也取消了
- [TODO] 后期可以考虑去掉或者使用 RPC 通信代替 HTTP

#### **Redis:**

- 暂时还不得知

### 2.3 修改对应的配置文件

- 将该项目下`config`目录 copy 至`/etc/kraken/`下。
- 也可自定义配置文件位置，启动的时候指定即可

### 2.3 build 并运行各组件(nova 后台服务)

#### 2.3.1 分别在对应目录下编译为可执行文件

- **testfs:** 在目录`nova/tools/bin/testfs`,执行`go build`
- **origin:** 在目录`nova/origin`,执行`go build`
- **tracker:** 在目录`nova/tracker`,执行`go build`
- **build-index:** 在目录`nova/build-index`,执行`go build`

#### 2.3.2 start.sh 启动脚本

start.sh

```bash
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
rm /tmp/kraken-push-cli-registry.sock
rm /tmp/kraken-push-cli-registry-override.sock
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

## 3: Push/Pull 使用介绍

### 3.1 PushCLI 上传文件客户端

#### 3.1.1 build

- 是根据 kraken 原项目的 proxy 组件修改的
- **proxy:** 在目录`nova/proxy`,执行`go build`

#### 3.1.2 run

具体参数可看源码 flags 说明，并按实际需要修改。

```go
// Flags defines push-cli CLI flags.
type Flags struct {
	ConfigFile    string // 指定配置文件
	KrakenCluster string
	SecretsFile   string
	UploadFile    string // 待上传的文件(绝对路径/filename)
	Tag           string // 待上传Tag,默认是"latest"
}
```

```bash
./push-cli \
--config=/etc/kraken/config/push-cli/base.yaml \
--file=/mnt/e/sse/sources4 \
--tag=latest
```

### 3.2 PullCLi 下载文件客户端

#### 3.2.1 build

- 新加的 CLI 工具，和 agent 对应
- **proxy:** 在目录`nova/agent`,执行`go build`
- **proxy:** 在目录`nova/pull-cli`,执行`go build`

#### 3.2.2 run
##### pull-agent
- **Agent**
- 注意 agent 暂时去掉了docker相关逻辑，还有nginx 配置逻辑。
```bash
// NOTE agent server 启动命令
// linux
./agent --config=/etc/kraken/config/agent/base.yaml \
--peer-ip=localhost \
--peer-port=16001 \
--agent-server-port=16002

// windows
./agent.exe --config=./base.yaml \
--peer-ip=localhost \
--peer-port=18001 \
--agent-server-port=18002- 
```
##### pull-cli

- **pull-cli**
  具体参数可看源码 flags 说明，并按实际需要修改。

```bash
// linux
./pull-cli \
--agent-server-port=16002 \
--tag=sse1.zip:latest \
--download-dir=/mnt/e/sse/

// windows
./agent.exe --config=./base.yaml \
--peer-ip=localhost \
--peer-port=18001 \
--agent-server-port=18002
```

