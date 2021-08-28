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
	"github.com/lppgo/nova/core"
	"github.com/lppgo/nova/lib/backend"
	"github.com/lppgo/nova/lib/blobrefresh"
	"github.com/lppgo/nova/lib/hashring"
	"github.com/lppgo/nova/lib/healthcheck"
	"github.com/lppgo/nova/lib/hostlist"
	"github.com/lppgo/nova/lib/metainfogen"
	"github.com/lppgo/nova/lib/persistedretry"
	"github.com/lppgo/nova/lib/store"
	"github.com/lppgo/nova/lib/torrent/networkevent"
	"github.com/lppgo/nova/lib/torrent/scheduler"
	"github.com/lppgo/nova/localdb"
	"github.com/lppgo/nova/metrics"
	"github.com/lppgo/nova/nginx"
	"github.com/lppgo/nova/origin/blobserver"
	"github.com/lppgo/nova/utils/httputil"

	"go.uber.org/zap"
)

// Config defines origin server configuration.
// TODO(evelynl94): consolidate cluster and hashring.
type Config struct {
	Verbose       bool
	ZapLogging    zap.Config               `yaml:"zap"`
	Cluster       hostlist.Config          `yaml:"cluster"`
	HashRing      hashring.Config          `yaml:"hashring"`
	HealthCheck   healthcheck.FilterConfig `yaml:"healthcheck"`
	BlobServer    blobserver.Config        `yaml:"blobserver"`
	CAStore       store.CAStoreConfig      `yaml:"castore"`
	Scheduler     scheduler.Config         `yaml:"scheduler"`
	NetworkEvent  networkevent.Config      `yaml:"network_event"`
	PeerIDFactory core.PeerIDFactory       `yaml:"peer_id_factory"`
	Metrics       metrics.Config           `yaml:"metrics"`
	MetaInfoGen   metainfogen.Config       `yaml:"metainfogen"`
	Backends      []backend.Config         `yaml:"backends"`
	Auth          backend.AuthConfig       `yaml:"auth"`
	BlobRefresh   blobrefresh.Config       `yaml:"blobrefresh"`
	LocalDB       localdb.Config           `yaml:"localdb"`
	WriteBack     persistedretry.Config    `yaml:"writeback"`
	Nginx         nginx.Config             `yaml:"nginx"`
	TLS           httputil.TLSConfig       `yaml:"tls"`
}
