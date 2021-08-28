// Copyright (c) 2016-2019 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lppgo/nova/utils/stringset"

	"github.com/docker/go-units"

	"github.com/pkg/errors"

	"github.com/lppgo/nova/build-index/tagclient"
	"github.com/lppgo/nova/core"
	"github.com/lppgo/nova/lib/dockerregistry/transfer"
	"github.com/lppgo/nova/lib/healthcheck"
	"github.com/lppgo/nova/lib/store"

	"github.com/lppgo/nova/lib/upstream"
	"github.com/lppgo/nova/metrics"
	"github.com/lppgo/nova/origin/blobclient"
	"github.com/lppgo/nova/utils/configutil"
	"github.com/lppgo/nova/utils/log"

	"github.com/flylog/colorstyle"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

const (
	Default_Namespace         = ".*"
	Build_Index_Put_Tag_Error = "cannot upload tag, missing dependency "
)

// Flags defines push-cli CLI flags.
type Flags struct {
	ConfigFile    string
	KrakenCluster string
	SecretsFile   string
	UploadFile    string //待上传的文件(绝对路径/filename)
	Tag           string //待上传Tag,默认是"latest"
}

// ParseFlags parses push-cli CLI flags.
func ParseFlags() *Flags {
	var flags Flags
	flag.StringVar(
		&flags.ConfigFile, "config", "", "configuration file path")
	flag.StringVar(
		&flags.KrakenCluster, "cluster", "", "cluster name (e.g. prod01-zone1)")
	flag.StringVar(
		&flags.SecretsFile, "secrets", "", "path to a secrets YAML file to load into configuration")
	flag.StringVar(
		&flags.UploadFile, "file", "", "file to be upload")
	flag.StringVar(
		&flags.Tag, "tag", "latest", "file to be upload")
	flag.Parse()

	fmt.Printf("\n%+v\n\n", flags)
	return &flags
}

type options struct {
	config  *Config
	metrics tally.Scope
	logger  *zap.Logger
}

// Option defines an optional Run parameter.
type Option func(*options)

// WithConfig ignores config/secrets flags and directly uses the provided config
// struct.
func WithConfig(c Config) Option {
	return func(o *options) { o.config = &c }
}

// WithMetrics ignores metrics config and directly uses the provided tally scope.
func WithMetrics(s tally.Scope) Option {
	return func(o *options) { o.metrics = s }
}

// WithLogger ignores logging config and directly uses the provided logger.
func WithLogger(l *zap.Logger) Option {
	return func(o *options) { o.logger = l }
}

// Run runs the push-cli.
func Run(flags *Flags, opts ...Option) {
	if flags.UploadFile == "" {
		panic("must specify a absolute path file")
	}

	var overrides options
	for _, o := range opts {
		o(&overrides)
	}

	var config Config
	if overrides.config != nil {
		config = *overrides.config
	} else {
		if err := configutil.Load(flags.ConfigFile, &config); err != nil {
			log.Fatalf("Load Config is Error:%s", err.Error())
		}
		if flags.SecretsFile != "" {
			if err := configutil.Load(flags.SecretsFile, &config); err != nil {
				panic(err)
			}
		}
	}

	if overrides.logger != nil {
		log.SetGlobalLogger(overrides.logger.Sugar())
	} else {
		zlog := log.ConfigureLogger(config.ZapLogging)
		zlog.Sync()
	}

	stats := overrides.metrics
	if stats == nil {
		s, closer, err := metrics.New(config.Metrics, flags.KrakenCluster)
		if err != nil {
			log.Fatalf("Failed to init metrics: %s", err)
		}
		stats = s
		defer closer.Close()
	}

	go metrics.EmitVersion(stats)

	cas, err := store.NewCAStore(config.CAStore, stats)
	if err != nil {
		log.Fatalf("Failed to create store: %s", err)
	}

	tls, err := config.TLS.BuildClient()
	if err != nil {
		log.Fatalf("Error building client tls config: %s", err)
	}

	origins, err := config.Origin.Build(upstream.WithHealthCheck(healthcheck.Default(tls)))
	if err != nil {
		log.Fatalf("Error building origin host list: %s", err)
	}

	r := blobclient.NewClientResolver(blobclient.NewProvider(blobclient.WithTLS(tls)), origins)
	originCluster := blobclient.NewClusterClient(r)

	buildIndexes, err := config.BuildIndex.Build(upstream.WithHealthCheck(healthcheck.Default(tls)))
	if err != nil {
		log.Fatalf("Error building build-index host list: %s", err)
	}

	tagClient := tagclient.NewClusterClient(buildIndexes, tls)

	namespace := Default_Namespace
	if config.Namespace != "" {
		namespace = config.Namespace
	}

	// TODO update ...
	var (
		dir      string
		filename string
		tag      string
	)
	oneSeconds := time.Second * 1
	startAt := time.Now()
	openFileTick := time.NewTicker(oneSeconds)
	fromReaderTick := time.NewTicker(oneSeconds)
	uploadTick := time.NewTicker(oneSeconds)

	dir, filename = filepath.Split(flags.UploadFile)
	tag = filename + ":" + flags.Tag
	log.Infof(" start operate for blobs [%s] and tags [%s] ...", flags.UploadFile, tag)

	digester := core.NewDigester()
	transferer := transfer.NewReadWriteTransferer(stats, tagClient, originCluster, cas)

	// if checkTagExist(transferer, filename, tag) {
	// 	log.Fatalf("tag :[%s] has exist !", tag)
	// }

	fileInfo, err := os.Lstat(flags.UploadFile)
	if err != nil {
		log.Fatalf("Lstat Err:%s", err.Error())
	}

	go openFileTicker(openFileTick)
	file, err := OpenFile(flags.UploadFile, uint64(fileInfo.Size()))
	openFileTick.Stop()
	if err != nil {
		log.Fatalf("OpenAndReadFile Err:%s", err.Error())
	}
	defer file.Close()

	go calculateDigestFromReaderTicker(fromReaderTick)
	digest, err := digester.FromReader(file)
	fromReaderTick.Stop()
	if err != nil {
		log.Fatalf("FromReader Err:%s", err.Error())
	} else {
		log.Infof("[%s] digest is [%s]", flags.UploadFile, digest.String())
	}

	//
	offset, whence := 0, io.SeekStart
	_, err = file.Seek(int64(offset), whence)
	if err != nil {
		log.Fatalf("file.Seek Err:%s", err.Error())
	}

	go uploadTicker(uploadTick)
	var cfer CustomFiler = CustomFile{file}
	err = uploadBlob(transferer, namespace, tag, digest, cfer)
	uploadTick.Stop()
	if err != nil {
		log.Fatalf("uploadBlob Err:%s", err.Error())
	} else {
		fmt.Printf("upload: dir:[%s] tag:[%s] succeed ; 耗时:[%s] !\n",
			colorstyle.Green(dir), colorstyle.Green(tag), colorstyle.Green(time.Since(startAt).String()))
	}

	err = putTag(transferer, tag, digest)
	if err != nil {
		log.Fatalf("puTag Err:%s", err.Error())
	}

	tags, err := getTag(transferer, filename)
	if err != nil {
		log.Errorf("getTag Err:%s", err.Error())
	}
	log.Infof("tags:%s,digest:%s", colorstyle.Green(tags), colorstyle.Green(digest.String()))
}

//

// OpenFile .
func OpenFile(path string, size uint64) (file *os.File, err error) {
	file, err = os.Open(path)
	if err != nil {
		log.Fatalf("Open File Err:%s", err.Error())
		err = errors.Wrap(err, "Open File Error")
		return
	}

	fmt.Printf("\r%s", strings.Repeat(" ", 5))
	fmt.Printf("\r======> Please waiting for opening file [%s] [%s] ...  <======\n", colorstyle.Green(path), colorstyle.Green(units.HumanSize(float64(size))))
	return
}

// checkTagExist 检查tag是否已存在.
func checkTagExist(transferer *transfer.ReadWriteTransferer, prefix, tag string) bool {
	log.Info(colorstyle.Yellow("----------------checkTagExist--------------------"))
	tags, err := transferer.ListTags(prefix)
	if err == nil {
		return stringset.FromSlice(tags).Has(tag)
	}
	return false
}

// puTag put digest对应的tag.
func putTag(transferer *transfer.ReadWriteTransferer, tag string, digest core.Digest) error {
	log.Info(colorstyle.Yellow("----------------putTag--------------------"))
	err := transferer.PutTag(tag, digest)
	// NOTE : 去origin upload和cache目录查看文件是否已存在.(首次上传不存在，返回 cannot upload tag, missing dependency sha256:xxx)
	ignorableError := Build_Index_Put_Tag_Error + digest.String()
	if err != nil {
		if !strings.Contains(err.Error(), ignorableError) {
			return err
		}
	}
	return nil
}

// getTag 根据prefix获取对应的tags.
func getTag(transferer *transfer.ReadWriteTransferer, prefix string) ([]string, error) {
	log.Info(colorstyle.Yellow("----------------getTag--------------------"))
	return transferer.ListTags(prefix)
}

// uploadBlob 上传一个blob 给origin.
func uploadBlob(transferer *transfer.ReadWriteTransferer, namespace, tag string, digest core.Digest, cfer CustomFiler) error {
	log.Info(colorstyle.Yellow("----------------uploadBlob----------------"))
	fmt.Printf("\r%s", strings.Repeat(" ", 5))
	fmt.Printf("\r======> Uploading [tag:%s]... %s  <======", colorstyle.Green(tag), colorstyle.Green(units.HumanSize(float64(cfer.Size()))))
	return transferer.Upload(namespace, digest, cfer)
}

//  openFileTicker .
func openFileTicker(tick *time.Ticker) {
	duration := 0
	for {
		select {
		case <-tick.C:
			duration++
			fmt.Printf("\r%s", strings.Repeat(" ", 5))
			fmt.Printf("\r======> OpenFile spend time  %s s ... <======", colorstyle.Green(duration))
		// case <-close:
		// 	tick.Stop()
		default:
		}
	}
}

// calculateDigestFromReaderTicker .
func calculateDigestFromReaderTicker(tick *time.Ticker) {
	duration := 0
	for {
		select {
		case <-tick.C:
			duration++
			fmt.Printf("\r%s", strings.Repeat(" ", 5))
			fmt.Printf("\r======> FromReader spend time  %s s ... <======\n", colorstyle.Green(duration))
		// case <-close:
		// 	tick.Stop()
		default:
		}
	}
}

// uploadTicker .
func uploadTicker(tick *time.Ticker) {
	duration := 0
	for {
		select {
		case <-tick.C:
			duration++
			fmt.Printf("\r%s", strings.Repeat(" ", 5))
			fmt.Printf("\r======> upLoadTicker spend time  %s s ... <======\n", colorstyle.Green(duration))
		// case <-close:
		// 	tick.Stop()
		default:
		}
	}
}
