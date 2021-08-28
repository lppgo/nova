<div align='center' ><font size='50'>参考数据中心-文件分发系统</font></div>
<div align='center' ><font size='50'>设计说明书</font></div>
<div align='center' ><font size='50'>阶段1-现有逻辑梳理</font></div>
[toc]

# 1: docker Registry API

## 1.1 官方文档

https://docs.docker.com/registry/spec/api/

## 1.2 个人整理文档

https://lppgo.github.io/post/20210728-01-docker-registry-http-api-v2/

# 2: kraken-agent

- 对 Tracker 返回的 peers 列表之间进行相互连接
- 提供文件下载的接口

## 2.1 kraken-agent 启动流程

![](./statics/017_G4参考数据中心设计图-agent-start.svg)

## 2.2 kraken-agent-server

```go
	r.Get("/health", handler.Wrap(s.healthHandler))

	r.Get("/tags/{tag}", handler.Wrap(s.getTagHandler))

	r.Get("/namespace/{namespace}/blobs/{digest}", handler.Wrap(s.downloadBlobHandler))

	r.Delete("/blobs/{digest}", handler.Wrap(s.deleteBlobHandler))

	// Preheat/preload endpoints.
	r.Get("/preload/tags/{tag}", handler.Wrap(s.preloadTagHandler))
```

### 2.1 健康检查 `r.Get("/health", ...)`

- 接口说明：服务健康检查
- 接口地址：`localhost:{{agent-server-port}}/health`

- 返回格式：`text/plain`

- 请求方式：`GET`

- 请求示例：`localhost:16002/health`

- 接口备注：

- 请求参数说明：

  | 名称              | 类型 | 必填 | 说明                                                            |
  | ----------------- | ---- | ---- | --------------------------------------------------------------- |
  | agent-server-port | int  | true | 启动的命令行参数 agent-server-port ，即 agent-server 服务的端口 |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：
    ```json
    OK
    ```

### 2.2 获取 tag `r.Get("/tags/{tag}", ...)`

- 接口说明：获取 tag
- 接口地址：`localhost:{{agent-server-port}}/tags/{tag}`

- 返回格式：`text/plain`

- 请求方式：`GET`

- 请求示例：`localhost:16002/tags/debian:latest`

- 接口备注：

- 请求参数说明：

  | 名称              | 类型   | 必填 | 说明                                                            |
  | ----------------- | ------ | ---- | --------------------------------------------------------------- |
  | agent-server-port | int    | true | 启动的命令行参数 agent-server-port ，即 agent-server 服务的端口 |
  | tag               | string | true | 镜像的 tag 信息，例如: `debian:latest`                          |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

    ```json
    sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86
    ```

### 2.3 下载 blob `r.Get("/namespace/{namespace}/blobs/{digest}", ...)`

- 接口说明：下载 blob
- 接口地址：`localhost:{{agent-server-port}}/namespace/{{namespace}}/blobs/{digest}`

- 返回格式：`Json`

- 请求方式：`GET`

- 请求示例：`localhost:16002/namespace/.*/blobs/sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86`

- 接口备注：

- 请求参数说明：

  | 名称              | 类型   | 必填 | 说明                                                            |
  | ----------------- | ------ | ---- | --------------------------------------------------------------- |
  | agent-server-port | int    | true | 启动的命令行参数 agent-server-port ，即 agent-server 服务的端口 |
  | namespace         | string | true | 指 docker repository 的 namespace,repositories                  |
  | digest            | string | true | 指的是 Image layer 的摘要 digest                                |

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：
    ```json
    {
      "schemaVersion": 2,
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "config": {
        "mediaType": "application/vnd.docker.container.image.v1+json",
        "size": 1462,
        "digest": "sha256:0980b84bde890bbdd5db43522a34b4f7c3c96f4d026527f4a7266f7ee408780d"
      },
      "layers": [
        {
          "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
          "size": 50435626,
          "digest": "sha256:627b765e08d177e63c9a202ca4991b711448905b934435c70b7cbd7d4a9c7959"
        }
      ]
    }
    ```

### 2.4 RemoveTorrent `r.Delete("/blobs/{digest}", ...)`

- 接口说明：根据 digest 从 disk 删除 torrent
- 接口地址：`localhost:{{agent-server-port}}/blobs/{digest}`

- 返回格式：`text/plain`

- 请求方式：`GET`

- 请求示例：`localhost:16002/blobs/sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86`

- 接口备注：

- 请求参数说明：

  | 名称              | 类型   | 必填 | 说明                                                            |
  | ----------------- | ------ | ---- | --------------------------------------------------------------- |
  | agent-server-port | int    | true | 启动的命令行参数 agent-server-port ，即 agent-server 服务的端口 |
  | digest            | string | true | 指的是 Image layer 的摘要 digest                                |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：
    ```json
    1
    ```

### 2.5 预下载 `r.Get("/preload/tags/{tag}", ...)`

- 接口说明：通过触发 docker daemon 下载指定的 docker image
- 接口地址：`localhost:{{agent-server-port}}/preload/tags/{tag}`

- 返回格式：`text/plain`

- 请求方式：`GET`

- 请求示例：`localhost:16002/preload/tags/debian:latest`

- 接口备注：

- 请求参数说明：

  | 名称              | 类型   | 必填 | 说明                                                            |
  | ----------------- | ------ | ---- | --------------------------------------------------------------- |
  | agent-server-port | int    | true | 启动的命令行参数 agent-server-port ，即 agent-server 服务的端口 |
  | tag               | string | true | 镜像的 tag 信息，例如: `debian:latest`                          |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：
    ```json
    1
    ```

## 2.3 kraken-agent-client

目前 agent-client 实现了下面 2 个 API：

```go
// Client defines a client for accessing the agent server.
type Client interface {

    // 调用 r.Get("/tags/{tag}", handler.Wrap(s.getTagHandler))
	GetTag(tag string) (core.Digest, error)

    // 调用 r.Get("/namespace/{namespace}/blobs/{digest}", handler.Wrap(s.downloadBlobHandler))
	Download(namespace string, d core.Digest) (io.ReadCloser, error)
}
```

# 3: kraken-proxy

- 实现文件注册接口
- 上传文件到对应的的 Origin

## 3.1 kraken-proxy 启动流程

![](./statics/017_G4参考数据中心设计图-proxy.svg)

## 3.2 registryoverride

```go
r.Get("/v2/_catalog", handler.Wrap(s.catalogHandler))
```

> overrides Docker registry endpoints.

- 接口说明：处理`/v2/_catalog`的请求，获取 repositry list
- 接口地址：`[localhost:{{port}}/health](http://localhost:15000/v2/_catalog)`

- 返回格式：`json`

- 请求方式：`GET`

- 请求示例：`http://localhost:15000/v2/_catalog`

- 接口备注：

- 请求参数说明：

  | 名称 | 类型 | 必填 | 说明                                                   |
  | ---- | ---- | ---- | ------------------------------------------------------ |
  | port | int  | true | 启动的命令行参数 port ，即 registryoverride 服务的端口 |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

    ```json
    { "repositories": ["debian"] }
    ```

## 3.3 proxyserver

```go
	r.Get("/health", handler.Wrap(s.healthHandler))

	r.Post("/registry/notifications", handler.Wrap(s.preheatHandler.Handle))
```

### 3.3.1 proxyserver 服务健康检查 `r.Get("/health", ...)`

- 接口说明：服务健康检查
- 接口地址：`localhost:{{proxy-server-port}}/health`

- 返回格式：`text/plain`

- 请求方式：`GET`

- 请求示例：`localhost:15005/health`

- 接口备注：

- 请求参数说明：

  | 名称              | 类型 | 必填 | 说明                                                      |
  | ----------------- | ---- | ---- | --------------------------------------------------------- |
  | proxy-server-port | int  | true | 启动的命令行参数 server-port ，即 proxy-server 服务的端口 |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：
    ```json
    OK
    ```

### 3.3.2 通知 Origin 缓存 Image 关联的 blob `r.Post("/registry/notifications",...)`

本质上是调用` clusterClient.DownloadBlob(namespace, d, dst)` ，
也就是 request origin server 的 url:`/namespace/{namespace}/blobs/{digest}`,
通过 origin 下载 blob 到`/var/cache/kraken/kraken-origin/cache`

- 接口说明：通知 origin 下载 Image 关联的 blob 到 origin cache 目录
- 接口地址：`localhost:15005/registry/notifications`

- 返回格式：`text/plain`

- 请求方式：`POST`

- 请求示例：`localhost:15005/registry/notifications`

- 接口备注：

- 请求参数说明（`kraken/proxy/proxyserver/registry_events.go/Notification` 结构体）,也可以通过 `TestPreheat` 测试函数构造：

```json
{
  "Events": [
    {
      "Id": "1",
      "TimeStamp": "2021-07-29T14:58:18.1345785+08:00",
      "Action": "push",
      "Target": {
        "MediaType": "application/vnd.docker.distribution.manifest.v2+json",
        "Digest": "sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86",
        "Repository": "debian",
        "Url": "",
        "Tag": "latest"
      }
    }
  ]
}
```

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：无返回参数

    ```json

    ```

# 4: kraken-origin

- 专用的播种器 Seeder
- 将文件存储在磁盘上或者文件系统(FS)或者对象存储
- 形成一个自我修复的 hash 环以分配负载

## 4.1 kraken-origin 启动流程

![](./statics/017_G4参考数据中心设计图-origin.svg)

## 4.2 blobserver

```go
	// Public endpoints:

	r.Get("/health", handler.Wrap(s.healthCheckHandler))

	r.Get("/blobs/{digest}/locations", handler.Wrap(s.getLocationsHandler))
	// upload
	r.Post("/namespace/{namespace}/blobs/{digest}/uploads", handler.Wrap(s.startClusterUploadHandler))
	r.Patch("/namespace/{namespace}/blobs/{digest}/uploads/{uid}", handler.Wrap(s.patchClusterUploadHandler))
	r.Put("/namespace/{namespace}/blobs/{digest}/uploads/{uid}", handler.Wrap(s.commitClusterUploadHandler))

	// download
	r.Get("/namespace/{namespace}/blobs/{digest}", handler.Wrap(s.downloadBlobHandler))

	r.Post("/namespace/{namespace}/blobs/{digest}/remote/{remote}", handler.Wrap(s.replicateToRemoteHandler))

	r.Post("/forcecleanup", handler.Wrap(s.forceCleanupHandler))

	// Internal endpoints:

	r.Post("/internal/blobs/{digest}/uploads", handler.Wrap(s.startTransferHandler))
	r.Patch("/internal/blobs/{digest}/uploads/{uid}", handler.Wrap(s.patchTransferHandler))
	r.Put("/internal/blobs/{digest}/uploads/{uid}", handler.Wrap(s.commitTransferHandler))

	r.Delete("/internal/blobs/{digest}", handler.Wrap(s.deleteBlobHandler))

	r.Post("/internal/blobs/{digest}/metainfo", handler.Wrap(s.overwriteMetaInfoHandler))

	r.Get("/internal/peercontext", handler.Wrap(s.getPeerContextHandler))

	r.Head("/internal/namespace/{namespace}/blobs/{digest}", handler.Wrap(s.statHandler))

	r.Get("/internal/namespace/{namespace}/blobs/{digest}/metainfo", handler.Wrap(s.getMetaInfoHandler))

	r.Put("/internal/duplicate/namespace/{namespace}/blobs/{digest}/uploads/{uid}",
		handler.Wrap(s.duplicateCommitClusterUploadHandler))
```

### 4.2.1 blobserver 服务健康检查`r.Get("/health", ...)`

- 接口说明：服务健康检查
- 接口地址：`http://localhost:{{blobserver-port}}/health`

- 返回格式：`text/plain`

- 请求方式：`GET`

- 请求示例：`localhost:15002/health`

- 接口备注：

- 请求参数说明：

  | 名称            | 类型 | 必填 | 说明                                                        |
  | --------------- | ---- | ---- | ----------------------------------------------------------- |
  | blobserver-port | int  | true | 启动的命令行参数 blobserver-port ，即 blobserver 服务的端口 |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：
    ```json
    OK
    ```

### 4.2.2 获取 digest 获取可用的 origin locations`r.Get("/blobs/{digest}/locations", ...)`

- 接口说明：获取 digest 获取可用的 origin locations,将结果以 response>header 返回
- 接口地址：`http://localhost:{{origin_server_port}}/namespace/{{namespace}}/blobs/{{debian-digest}}/uploads`

- 返回格式：`text/plain`

- 请求方式：`GET`

- 请求示例：`http://localhost:15002/blobs/sha256:3e24baa60967d085b95a45129f82af4eb9d1e33aff9559173542ebb15c5d9cb5/locations`

- 接口备注：

- 请求参数说明：

  | 名称            | 类型   | 必填 | 说明                                                                                               |
  | --------------- | ------ | ---- | -------------------------------------------------------------------------------------------------- |
  | blobserver-port | int    | true | 启动的命令行参数 blobserver-port ，即 blobserver 服务的端口                                        |
  | debian-digest   | string | true | 要操作镜像的 digest ,比如`sha256:3e24baa60967d085b95a45129f82af4eb9d1e33aff9559173542ebb15c5d9cb5` |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
    | Origin-Locations |string | response.Header |
  - response Body 返回示例：
    ```json
    1
    ```

### 4.2.3 为外部上传初始化 uploader `r.Post("/namespace/{namespace}/blobs/{digest}/uploads", ...)`

- 接口说明：初始化并启动 uploader,并根据 digest 返回一个待上传 layer 的 uuid,并在 uploader 所在 dir 创建一个 uuid 对应的文件
- 接口地址：`http://localhost:{{blobserver-port}}/namespace/{{namespace}}/blobs/{{digest}}/uploads`

- 返回格式：`text/plain`

- 请求方式：`POST`

- 请求示例：`http://localhost:15002/blobs/sha256:3e24baa60967d085b95a45129f82af4eb9d1e33aff9559173542ebb15c5d9cb5/locations`

- 接口备注：

- 请求参数说明：

  | 名称            | 类型   | 必填 | 说明                                                                                                      |
  | --------------- | ------ | ---- | --------------------------------------------------------------------------------------------------------- |
  | blobserver-port | int    | true | 启动的命令行参数 blobserver-port ，即 blobserver 服务的端口                                               |
  | namespace       | string | true | namespace 命名空间                                                                                        |
  | digest          | string | true | 要操作镜像 layer 的 digest ,比如`sha256:3e24baa60967d085b95a45129f82af4eb9d1e33aff9559173542ebb15c5d9cb5` |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
    | Location |string | response.Header,根据生成的 uuid,比如:`231afdb9-c6d9-456c-94b5-dbd175c95d84`,同时在 uploader 的 dir 生成对应的文件 |
  - response Body 返回示例：
    ```json
    1
    ```

### 4.2.4 分块上传`r.Patch("/namespace/{namespace}/blobs/{digest}/uploads/{uid}", ...)`

- 接口说明：uploads a chunk of a blob for external uploads.
- 接口地址：`http://localhost:{{blobserver-port}}/namespace/{{namespace}}/blobs/{{digest}}/uploads`

- 返回格式：`text/plain`

- 请求方式：`Patch`

- 请求示例：``

- 接口备注：

- 请求参数说明：

  | 名称            | 类型   | 必填 | 说明                                                                                                      |
  | --------------- | ------ | ---- | --------------------------------------------------------------------------------------------------------- |
  | blobserver-port | int    | true | 启动的命令行参数 blobserver-port ，即 blobserver 服务的端口                                               |
  | namespace       | string | true | namespace 命名空间                                                                                        |
  | digest          | string | true | 要操作镜像 layer 的 digest ,比如`sha256:3e24baa60967d085b95a45129f82af4eb9d1e33aff9559173542ebb15c5d9cb5` |

```bash
>>> PATCH /v2/<name>/blobs/uploads/<uuid>
>>> Content-Length: <size of chunk>
>>> Content-Range: <start of range>-<end of range>
>>> Content-Type: application/octet-stream
<Layer Chunk Binary Data>
```

- 返回参数说明：(无)

### 4.2.5 异步提交 blob 上传`r.Put("/namespace/{namespace}/blobs/{digest}/uploads/{uid}", ...)`

- 接口说明：commits an external blob upload asynchronously,meaning the blob will be written back to
  remote storage in a non-blocking fashion.
- 接口地址：`http://localhost:{{blobserver-port}}/namespace/{{namespace}}/blobs/{{digest}}/uploads`

- 返回格式：`text/plain`

- 请求方式：`Put`

- 请求示例：``

- 接口备注：

- 请求参数说明：

  | 名称            | 类型   | 必填 | 说明                                                                                                      |
  | --------------- | ------ | ---- | --------------------------------------------------------------------------------------------------------- |
  | blobserver-port | int    | true | 启动的命令行参数 blobserver-port ，即 blobserver 服务的端口                                               |
  | namespace       | string | true | namespace 命名空间                                                                                        |
  | digest          | string | true | 要操作镜像 layer 的 digest ,比如`sha256:3e24baa60967d085b95a45129f82af4eb9d1e33aff9559173542ebb15c5d9cb5` |

### 4.2.6 获取 manifest `r.Get("/namespace/{namespace}/blobs/{digest}",...)`

- 接口说明：获取 manifest
- 接口地址：`http://localhost:{{blobserver-port}}/namespace/{{namespace}}/blobs/{{debian-digest}}`

- 返回格式：`json`

- 请求方式：`GET`

- 请求示例：`http://localhost:15002/namespace/.*/blobs/sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86`

- 接口备注：

- 请求参数说明：

  | 名称            | 类型   | 必填 | 说明                                                                                                      |
  | --------------- | ------ | ---- | --------------------------------------------------------------------------------------------------------- |
  | blobserver-port | int    | true | 启动的命令行参数 blobserver-port ，即 blobserver 服务的端口                                               |
  | namespace       | string | true | namespace 命名空间,`.*`                                                                                   |
  | digest          | string | true | 要操作镜像 layer 的 digest ,比如`sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86` |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
    | Content-Type |string | application/octet-stream-v1 |
  - response Body 返回示例：
    ```json
    {
      "schemaVersion": 2,
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "config": {
        "mediaType": "application/vnd.docker.container.image.v1+json",
        "size": 1462,
        "digest": "sha256:0980b84bde890bbdd5db43522a34b4f7c3c96f4d026527f4a7266f7ee408780d"
      },
      "layers": [
        {
          "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
          "size": 50435626,
          "digest": "sha256:627b765e08d177e63c9a202ca4991b711448905b934435c70b7cbd7d4a9c7959"
        }
      ]
    }
    ```

### 4.2.7 上传 blob 到远程 origin cluster`r.Post("/namespace/{namespace}/blobs/{digest}/remote/{remote}", ...)`

- 接口说明：将 digest 对应的 blob 复制到远程 origin cluster
- 接口地址：`http://localhost:{{blobserver-port}}/namespace/{{namespace}}/blobs/{{digest}}/remote/{{remote}}`

- 返回格式：`text`

- 请求方式：`POST`

- 请求示例：`http://localhost:15002/namespace/.*/blobs/sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86/remote/localhost:15002`

- 接口备注：

- 请求参数说明：

  | 名称            | 类型   | 必填 | 说明                                                        |
  | --------------- | ------ | ---- | ----------------------------------------------------------- |
  | blobserver-port | int    | true | 启动的命令行参数 blobserver-port ，即 blobserver 服务的端口 |
  | namespace       | string | true | namespace 命名空间,`.*`                                     |
  | remote          | string | true | 指的是其他远程 origin cluster                               |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：
    `1 `

### 4.2.8 清理 cache files `r.Post("/forcecleanup", ...)`

- 接口说明：通过设置 TTL 来清理 cache files
- 接口地址：`http://localhost:{{blobserver-port}}/forcecleanup?ttl_hr=10`

- 返回格式：`json`

- 请求方式：`POST`

- 请求示例：`http://localhost:15002/forcecleanup?ttl_hr=10`

- 接口备注：

- 请求参数说明：

  | 名称            | 类型   | 必填 | 说明                                                        |
  | --------------- | ------ | ---- | ----------------------------------------------------------- |
  | blobserver-port | int    | true | 启动的命令行参数 blobserver-port ，即 blobserver 服务的端口 |
  | ttl_hr          | string | true | url 参数(query arguments)，指的是文件 TTL 的小时数          |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：
    ```json
    {
      "deleted": null,
      "errors": null
    }
    ```

## 4.3 blobclient

### 4.3.1 Client

```go
// Client provides a wrapper around all Server HTTP endpoints.
type Client interface {

    // Addr返回一个给client提供的server地址
	Addr() string

    // Locations returns the origin server addresses which d is sharded on.
	Locations(d core.Digest) ([]string, error)
    // DeleteBlob deletes the blob corresponding to d.
	DeleteBlob(d core.Digest) error
    // TransferBlob uploads a blob to a single origin server.
	TransferBlob(d core.Digest, blob io.Reader) error

    // Stat return blob info.
	Stat(namespace string, d core.Digest) (*core.BlobInfo, error)
    // Stat return local blob info.
	StatLocal(namespace string, d core.Digest) (*core.BlobInfo, error)

    // GetMateInfo return metainfo for d.
	GetMetaInfo(namespace string, d core.Digest) (*core.MetaInfo, error)
    // OverwriteMetaInfo overwrites existing metainfo for d with new metainfo configured with pieceLength.
	OverwriteMetaInfo(d core.Digest, pieceLength int64) error

    // UploadBlob 上传并复制blob到origin cluster. 异步的备份Blob到remote storage configured for namespace.
	UploadBlob(namespace string, d core.Digest, blob io.Reader) error
    // DuplicateUploadBlob duplicates an blob upload request, which will attempt to
    // write-back at the given delay.
	DuplicateUploadBlob(namespace string, d core.Digest, blob io.Reader, delay time.Duration) error

    // DownloadBlob downloads blob for d.文件状态错误返回202，文件不存在返回404.
	DownloadBlob(namespace string, d core.Digest, dst io.Writer) error

    // ReplicateToRemote replicates the blob of d to a remote origin cluster.
	ReplicateToRemote(namespace string, d core.Digest, remoteDNS string) error
    // GetPeerContext 获取与服务器一起运行的 p2p 客户端的 PeerContext.
	GetPeerContext() (core.PeerContext, error)

    // ForceCleanup forces cache cleanup to run.
	ForceCleanup(ttl time.Duration) error
}
```

### 4.3.2 ClusterClient

```go
// ClusterClient defines a top-level origin cluster client which handles blob
// location resolution and retries.
type ClusterClient interface {
    // UploadBlob uploads blob to origin cluster.
	UploadBlob(namespace string, d core.Digest, blob io.Reader) error
    // DownloadBlob pulls a blob from the origin cluster.
	DownloadBlob(namespace string, d core.Digest, dst io.Writer) error
    // GetMetaInfo returns the metainfo for d.
	GetMetaInfo(namespace string, d core.Digest) (*core.MetaInfo, error)
    // Stat checks availability of a blob in the cluster.
	Stat(namespace string, d core.Digest) (*core.BlobInfo, error)

    // OverwriteMetaInfo overwrites existing metainfo for d with new metainfo configured
    // with pieceLength on every origin server.
	OverwriteMetaInfo(d core.Digest, pieceLength int64) error
    // Owners returns the origin peers which own d.
	Owners(d core.Digest) ([]core.PeerContext, error)
    // ReplicateToRemote replicates d to a remote origin cluster.
	ReplicateToRemote(namespace string, d core.Digest, remoteDNS string) error
}
```

#### 4.3.2.1 ClusterClient.UploadBlob

上传一个 blob 到 origin cluster

- 根据 digest 获取 origin cluster 的 clients;
- 每个 origin cluster 的 client 执行`client.UploadBlob(namespace, d, blob)`;
- 执行`runChunkedUpload`;
- 执行`runChunkedUploadHelper`;
- uploader.path 进行分块上传;

```go
// clusterClient UploadBlob
func (c *clusterClient) UploadBlob(namespace string, d core.Digest, blob io.Reader) (err error) {
	clients, err := c.resolver.Resolve(d)
	if err != nil {
		return fmt.Errorf("resolve clients: %s", err)
	}

	// We prefer the origin with highest hashing score so the first origin will handle
	// replication to origins with lower score. This is because we want to reduce upload
	// conflicts between local replicas.
	for _, client := range clients {
		err = client.UploadBlob(namespace, d, blob)
		// Allow retry on another origin if the current upstream is temporarily
		// unavailable or under high load.
		if httputil.IsNetworkError(err) || httputil.IsRetryable(err) {
			continue
		}
		break
	}
	return err
}

// ...
// 调用patch分块上传
func runChunkedUploadHelper(u uploader, d core.Digest, blob io.Reader, chunkSize int64) error {
	uid, err := u.start(d)
	if err != nil {
		return err
	}
	var pos int64
	buf := make([]byte, chunkSize)
	for {
		n, err := blob.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read blob: %s", err)
		}
		chunk := bytes.NewReader(buf[:n])
		stop := pos + int64(n)
		if err := u.patch(d, uid, pos, stop, chunk); err != nil {
			return err
		}
		pos = stop
	}
	return u.commit(d, uid)
}

// ...
// path 请求blobserver 对应的handler
r.Patch("/internal/blobs/{digest}/uploads/{uid}", handler.Wrap(s.patchTransferHandler))
```

#### 4.3.2.2 ClusterClient.DownloadBlob

从 origin cluster 下载一个 blob

- 调用`client.DownloadBlob(namespace, d, dst)`;
- client 调用 blobserver 对应的 http 方法`r.Get("/namespace/{namespace}/blobs/{digest}", handler.Wrap(s.downloadBlobHandler))`;

```go
// DownloadBlob pulls a blob from the origin cluster.
func (c *clusterClient) DownloadBlob(namespace string, d core.Digest, dst io.Writer) error {
	err := Poll(c.resolver, c.defaultPollBackOff(), d, func(client Client) error {
		return client.DownloadBlob(namespace, d, dst)
	})
	if httputil.IsNotFound(err) {
		err = ErrBlobNotFound
	}
	return err
}

// ...

// DownloadBlob .
func (c *HTTPClient) DownloadBlob(namespace string, d core.Digest, dst io.Writer) error {
	r, err := httputil.Get(
		fmt.Sprintf("http://%s/namespace/%s/blobs/%s", c.addr, url.PathEscape(namespace), d),
		httputil.SendTLS(c.tls))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if _, err := io.Copy(dst, r.Body); err != nil {
		return fmt.Errorf("copy body: %s", err)
	}
	return nil
}

// ...
// get请求blobserver 对应的handler
r.Get("/namespace/{namespace}/blobs/{digest}", handler.Wrap(s.downloadBlobHandler))

```

#### 4.3.2.3 ClusterClient.GetMetaInfo

获取 metainfo

```go
// GetMetaInfo returns the metainfo for d.
func (c *clusterClient) GetMetaInfo(namespace string, d core.Digest) (mi *core.MetaInfo, err error) {
	clients, err := c.resolver.Resolve(d)
	if err != nil {
		return nil, fmt.Errorf("resolve clients: %s", err)
	}
	for _, client := range clients {
		mi, err = client.GetMetaInfo(namespace, d)
		// Do not try the next replica on 202 errors.
		if err != nil && !httputil.IsAccepted(err) {
			continue
		}
		break
	}
	return mi, err
}

// ...

//
func (c *HTTPClient) GetMetaInfo(namespace string, d core.Digest) (*core.MetaInfo, error) {
	r, err := httputil.Get(
		fmt.Sprintf("http://%s/internal/namespace/%s/blobs/%s/metainfo",
			c.addr, url.PathEscape(namespace), d),
		httputil.SendTimeout(15*time.Second),
		httputil.SendTLS(c.tls))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %s", err)
	}
	mi, err := core.DeserializeMetaInfo(raw)
	if err != nil {
		return nil, fmt.Errorf("deserialize metainfo: %s", err)
	}
	return mi, nil
}

// ...
// get 请求blobserver 对应的handler
r.Get("/internal/namespace/{namespace}/blobs/{digest}/metainfo", handler.Wrap(s.getMetaInfoHandler))
```

#### 4.3.2.4 ClusterClient.Stat

检查 blob 在 origin cluster 的可用性

- 调用 client`client.Stat(namespace, d)`;
- client.Stat()返回一个 blobinfo,调用 blobserver 对应的 API;

```go
//
func (c *clusterClient) Stat(namespace string, d core.Digest) (bi *core.BlobInfo, err error) {
	clients, err := c.resolver.Resolve(d)
	if err != nil {
		return nil, fmt.Errorf("resolve clients: %s", err)
	}

	shuffle(clients)
	for _, client := range clients {
		bi, err = client.Stat(namespace, d)
		if err != nil {
			continue
		}
		break
	}

	return bi, err
}

// get 请求blobserver 对应的handler
r.Get("/internal/namespace/{namespace}/blobs/{digest}/metainfo", handler.Wrap(s.getMetaInfoHandler))
```

#### 4.3.2.5 ClusterClient.OverwriteMetaInfo

暂无有效用途，略。

#### 4.3.2.6 ClusterClient.Owners

返回 d 所在的 origin peers

```go
// Owners returns the origin peers which own d.
func (c *clusterClient) Owners(d core.Digest) ([]core.PeerContext, error) {
	clients, err := c.resolver.Resolve(d)
	if err != nil {
		return nil, fmt.Errorf("resolve clients: %s", err)
	}

	var mu sync.Mutex
	var peers []core.PeerContext
	var errs []error

	var wg sync.WaitGroup
	for _, client := range clients {
		wg.Add(1)
		go func(client Client) {
			defer wg.Done()
			pctx, err := client.GetPeerContext()
			mu.Lock()
			if err != nil {
				errs = append(errs, err)
			} else {
				peers = append(peers, pctx)
			}
			mu.Unlock()
		}(client)
	}
	wg.Wait()

	err = errutil.Join(errs)

	if len(peers) == 0 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("no origin peers found")
	}

	if err != nil {
		log.With("blob", d.Hex()).Errorf("Error getting all origin peers: %s", err)
	}
	return peers, nil
}
```

#### 4.3.2.7 ClusterClient.ReplicateToRemote

将 digest 复制到远程 origin cluster。

```go
// ReplicateToRemote replicates d to a remote origin cluster.
func (c *clusterClient) ReplicateToRemote(namespace string, d core.Digest, remoteDNS string) error {
	// Re-use download backoff since replicate may download blobs.
	return Poll(c.resolver, c.defaultPollBackOff(), d, func(client Client) error {
		return client.ReplicateToRemote(namespace, d, remoteDNS)
	})
}

//
func (c *HTTPClient) ReplicateToRemote(namespace string, d core.Digest, remoteDNS string) error {
	_, err := httputil.Post(
		fmt.Sprintf("http://%s/namespace/%s/blobs/%s/remote/%s",
			c.addr, url.PathEscape(namespace), d, remoteDNS),
		httputil.SendTLS(c.tls))
	return err
}

// post 请求blobserver 对应的handler
r.Post("/namespace/{namespace}/blobs/{digest}/remote/{remote}", handler.Wrap(s.replicateToRemoteHandler))
```

# 5: kraken-tracker

- 追踪对等节点 peer 有哪些数据
- 提供一个对等节点 peer 连接任何给定文件的有序列表

## 5.1 kraken-tracker 启动流程

![](./statics/017_G4参考数据中心设计图-tracker.svg)

## 5.2 tracker-server

```go
	r.Get("/health", handler.Wrap(s.healthHandler))
	r.Get("/announce", handler.Wrap(s.announceHandlerV1))
	r.Post("/announce/{infohash}", handler.Wrap(s.announceHandlerV2))
	r.Get("/namespace/{namespace}/blobs/{digest}/metainfo", handler.Wrap(s.getMetaInfoHandler))
```

### 5.2.1 tracker-server 健康检查`r.Get("/health", ...)`

- 接口说明：服务健康检查
- 接口地址：`localhost:{{tracker-server-port}}/health`

- 返回格式：`text/plain`

- 请求方式：`GET`

- 请求示例：`localhost:15003/health`

- 接口备注：

- 请求参数说明：

  | 名称                | 类型 | 必填 | 说明                                                 |
  | ------------------- | ---- | ---- | ---------------------------------------------------- |
  | tracker-server-port | int  | true | 启动的命令行参数 port ，即 tracker-server 服务的端口 |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：
    ```json
    OK
    ```

### 5.2.2 宣告并获取 peers V1`r.Get("/announce", ...)`

V1 版本暂无使用，略。

### 5.2.3 宣告并获取 peers V2`r.Post("/announce/{infohash}", ...)`

- 接口说明：宣告并获取 peer list
- 接口地址：`localhost:{{tracker-server-port}}/announce/{infohash}`

- 返回格式：`text/plain`

- 请求方式：`Post`

- 请求示例：`http://localhost:15003/announce/59f442ba854913196f5da7e3281a1508f579b9d5`

- 接口备注：

- 请求参数说明：

  | 名称                | 类型   | 必填 | 说明                                                                                                                                                                                                                  |
  | ------------------- | ------ | ---- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
  | tracker-server-port | int    | true | 启动的命令行参数 port ，即 tracker-server 服务的端口                                                                                                                                                                  |
  | infohash            | string | true | 某个 digest 文件的`_torrentmeta` 取 sha1 值。比如`/var/cache/kraken/kraken-origin/cache//0c/49/0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86` 取 sha1 得`59f442ba854913196f5da7e3281a1508f579b9d5` |

- Request.Body（具体格式参见`/tracker/announceclient/client.go/Request{}`）:

  ```json
  {
    "name": "0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86",
    "digest": "sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86",
    "peer": {
      "ip": "localhost",
      "port": 17001,
      "origin": true,
      "complete": false
    }
  }
  ```

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  {
    // peer 列表
    "peers": [
      {
        "peer_id": [
          48, 112, 58, 121, 175, 35, 1, 202, 7, 228, 196, 125, 69, 114, 132,
          159, 87, 26, 187, 152
        ],
        "ip": "localhost",
        "port": 15001,
        "origin": true,
        "complete": true
      },
      {
        "peer_id": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
        "ip": "localhost",
        "port": 17001,
        "origin": false,
        "complete": false
      }
    ],
    "interval": 3000000000
  }
  ```

### 5.2.3 获取 metainfo`r.Get("/namespace/{namespace}/blobs/{digest}/metainfo", ...)`

- 接口说明：根据 namespace 和 digest 获取 metainfo,也就是`_torrentmeta`文件的 content。
- 接口地址：`localhost:{{tracker-server-port}}/namespace/{namespace}/blobs/{digest}/metainfo`

- 返回格式：`json`

- 请求方式：`Get`

- 请求示例：`http://localhost:15003/namespace/debian/blobs/sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86/metainfo`

- 接口备注：

- 请求参数说明：

  | 名称                | 类型   | 必填 | 说明                                                 |
  | ------------------- | ------ | ---- | ---------------------------------------------------- |
  | tracker-server-port | int    | true | 启动的命令行参数 port ，即 tracker-server 服务的端口 |
  | namespace           | string | true | 指 docker repository 的 namespace,repositories       |
  | digest              | string | true | 指的是 Image layer 的摘要 digest                     |

- Request.Body（具体格式参见`/tracker/announceclient/client.go/Request{}`）:

  ```json
  {
    "name": "0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86",
    "digest": "sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86",
    "peer": {
      "ip": "localhost",
      "port": 17001,
      "origin": true,
      "complete": false
    }
  }
  ```

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  {
    "Info": {
      "PieceLength": 4194304,
      "PieceSums": [3655704501],
      "Name": "0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86",
      "Length": 529
    }
  }
  ```

# 6: kraken-build-index

- 映射人们可读的 tag 给文件摘要
- Build-Index 不保证一致性：clientCLI 应该使用唯一的 tags。
- 增强 file 在 cluster 中的复制(具有重试功能的简单重复队列)
- 将 tags 作为文件存储在文件系统

## 6.1 kraken-build-index 启动流程

![](./statics/017_G4参考数据中心设计图-build-index.svg)

## 6.2 tagserver

tagserver Server 给 kraken-build-index 提供了 tag 操作。

```go
    // 服务健康检查
	r.Get("/health", handler.Wrap(s.healthHandler))

    // put tag 和对应的 digest 到本地 disk 和复制到远程 backend
	r.Put("/tags/{tag}/digest/{digest}", handler.Wrap(s.putTagHandler))
    // 判断该 tag 是否已存在
	r.Head("/tags/{tag}", handler.Wrap(s.hasTagHandler))
    // 根据tag获取对应的digest
	r.Get("/tags/{tag}", handler.Wrap(s.getTagHandler))

    // list repositry里面的tags
	r.Get("/repositories/{repo}/tags", handler.Wrap(s.listRepositoryHandler))

    // list image
	r.Get("/list/*", handler.Wrap(s.listHandler))
    // 复制tag到远程存储
	r.Post("/remotes/tags/{tag}", handler.Wrap(s.replicateTagHandler))
    // 获取一个origin地址
	r.Get("/origin", handler.Wrap(s.getOriginHandler))

    // 重复 复制tag到远程的操作
	r.Post("/internal/duplicate/remotes/tags/{tag}/digest/{digest}",handler.Wrap(s.duplicateReplicateTagHandler))

    // 重复 put tag 操作
	r.Put("/internal/duplicate/tags/{tag}/digest/{digest}",handler.Wrap(s.duplicatePutTagHandler))
```

## 6.3 tagclient

```go
// Client wraps tagserver endpoints.
type Client interface {
	Put(tag string, d core.Digest) error
	PutAndReplicate(tag string, d core.Digest) error
	Get(tag string) (core.Digest, error)
	Has(tag string) (bool, error)
	List(prefix string) ([]string, error)
	ListWithPagination(prefix string, filter ListFilter) (tagmodels.ListResponse, error)
	ListRepository(repo string) ([]string, error)
	ListRepositoryWithPagination(repo string, filter ListFilter) (tagmodels.ListResponse, error)
	Replicate(tag string) error
	Origin() (string, error)

	DuplicateReplicate(tag string, d core.Digest, dependencies core.DigestList, delay time.Duration) error
	DuplicatePut(tag string, d core.Digest, delay time.Duration) error
}
```

### 6.3.1 Put `r.Put("/tags/{tag}/digest/{digest}", ...)`

- 接口说明： put tag 和对应的 digest 到本地 disk
- 接口地址：`http://localhost:{{port}}/tags/{{tag}}/digest/{{digest}}`

- 返回格式：`text/plain`

- 请求方式：`Put`

- 请求示例：`http://localhost:15004/tags/debian:latest/digest/sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86`

- 接口备注：

- 请求参数说明：

  | 名称   | 类型   | 必填 | 说明                                                     |
  | ------ | ------ | ---- | -------------------------------------------------------- |
  | port   | int    | true | 启动的命令行参数 port ，即 build-index-server 服务的端口 |
  | tag    | string | true | 某个 Image 的 tag ,比如: debian:latest                   |
  | digest | string | true | 指的是 Image layer 的摘要 digest                         |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  1
  ```

### 6.3.2 PutAndReplicate `r.Put("/tags/{tag}/digest/{digest}", ...)`

- 接口说明： put tag 和对应的 digest 到本地 disk 和复制到远程 backend
- 接口地址：`http://localhost:{{port}}/tags/{{tag}}/digest/{{digest}}?replicate=true`

- 返回格式：`text/plain`

- 请求方式：`Put`

- 请求示例：`http://localhost:15004/tags/debian:latest/digest/sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86?replicate=true`

- 接口备注：

- 请求参数说明：

  | 名称      | 类型   | 必填 | 说明                                                     |
  | --------- | ------ | ---- | -------------------------------------------------------- |
  | port      | int    | true | 启动的命令行参数 port ，即 build-index-server 服务的端口 |
  | tag       | string | true | 某个 Image 的 tag ,比如: debian:latest                   |
  | digest    | string | true | 指的是 Image layer 的摘要 digest                         |
  | replicate | bool   | true | replicate 是否复制 digest 到 remote                      |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  1
  ```

### 6.3.3 根据 tag 获取对应的 digest `r.Get("/tags/{tag}", ...)`

- 接口说明： 根据 tag 获取对应的 digest
- 接口地址：`http://localhost:{{port}}/tags/{{tag}}`

- 返回格式：`text/plain`

- 请求方式：`Get`

- 请求示例：`http://localhost:15004/tags/debian:latest`

- 接口备注：

- 请求参数说明：

  | 名称 | 类型   | 必填 | 说明                                                     |
  | ---- | ------ | ---- | -------------------------------------------------------- |
  | port | int    | true | 启动的命令行参数 port ，即 build-index-server 服务的端口 |
  | tag  | string | true | 某个 Image 的 tag ,比如: debian:latest                   |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86
  ```

### 6.3.4 判断该 tag 是否已存在 `r.Head("/tags/{tag}", ...)`

- 接口说明： 判断该 tag 是否已存在
- 接口地址：`http://localhost:{{port}}/tags/{{tag}}`

- 返回格式：`text/plain`

- 请求方式：`Head`

- 请求示例：`http://localhost:15004/tags/debian:latest`

- 接口备注：

- 请求参数说明：

  | 名称 | 类型   | 必填 | 说明                                                     |
  | ---- | ------ | ---- | -------------------------------------------------------- |
  | port | int    | true | 启动的命令行参数 port ，即 build-index-server 服务的端口 |
  | tag  | string | true | 某个 Image 的 tag ,比如: debian:latest ('repo:tag')      |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  1
  ```

### 6.3.5 根据前缀获取 tag 列表 `r.Get("/list/*", ...)`

- 接口说明： 根据前缀获取 tag 列表
- 接口地址：`http://localhost:{{port}}/list/{{prefix}}`

- 返回格式：`json`

- 请求方式：`Head`

- 请求示例：`http://localhost:15004/list/debian`

- 接口备注：

- 请求参数说明：

  | 名称   | 类型   | 必填 | 说明                                                     |
  | ------ | ------ | ---- | -------------------------------------------------------- |
  | port   | int    | true | 启动的命令行参数 port ，即 build-index-server 服务的端口 |
  | prefix | string | true | tag 的前缀 ,比如: debian                                 |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  {
    "Links": {
      "next": "",
      "self": "/list/debian"
    },
    "size": 1,
    "result": ["debian:latest"]
  }
  ```

### 6.3.6 list 某个 repositry 下面所有的 tags

- 接口说明： 根据前缀获取 tag 列表
- 接口地址：`http://localhost:{{port}}/list/{{prefix}}`

- 返回格式：`json`

- 请求方式：`Head`

- 请求示例：`http://localhost:15004/list/debian`

- 接口备注：

- 请求参数说明：

  | 名称   | 类型   | 必填 | 说明                                                     |
  | ------ | ------ | ---- | -------------------------------------------------------- |
  | port   | int    | true | 启动的命令行参数 port ，即 build-index-server 服务的端口 |
  | prefix | string | true | tag 的前缀 ,比如: debian                                 |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  {
    "Links": {
      "next": "",
      "self": "/list/debian"
    },
    "size": 1,
    "result": ["debian:latest"]
  }
  ```

### 6.3.7 复制 digest 到远程`r.Post("/remotes/tags/{tag}", ...)`

- 接口说明： 复制 digest 到远程
- 接口地址：`http://localhost:{{port}}/remotes/tags/{{tag}}`

- 返回格式：`text/plain`

- 请求方式：`Post`

- 请求示例：`http://localhost:15004/remotes/tags/debian:latest`

- 接口备注：

- 请求参数说明：

  | 名称 | 类型   | 必填 | 说明                                                     |
  | ---- | ------ | ---- | -------------------------------------------------------- |
  | port | int    | true | 启动的命令行参数 port ，即 build-index-server 服务的端口 |
  | tag  | string | true | 某个 Image 的 tag ,比如: debian:latest ('repo:tag')      |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  1
  ```

### 6.3.8 重新 Put `r.Put("/internal/duplicate/tags/{tag}/digest/{digest}",...)`

- 接口说明： 复制 digest 到远程
- 接口地址：`http://localhost:{{port}}/internal/duplicate/tags/{tag}/digest/{digest}`

- 返回格式：`text/plain`

- 请求方式：`Post`

- 请求示例：`http://localhost:15004/internal/duplicate/tags/debian:latest/digest/sha256:0c493c517180a20aecaa80fe12ec594476df79deabbed42296169380df962e86`

- 接口备注：

- 请求参数说明：

  | 名称   | 类型   | 必填 | 说明                                                     |
  | ------ | ------ | ---- | -------------------------------------------------------- |
  | port   | int    | true | 启动的命令行参数 port ，即 build-index-server 服务的端口 |
  | tag    | string | true | 某个 Image 的 tag ,比如: debian:latest ('repo:tag')      |
  | digest | string | true | 指的是 tag 对应 Image 的摘要 digest                      |

- 返回参数说明：

  - response Header
    默认
    | 名称 | 类型 | 说明 |
    | ---- | ------ | -------- |
    | 200 | int | 状态码 |
  - response Body 返回示例：

  ```json
  1
  ```
