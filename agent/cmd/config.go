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
	"github.com/lppgo/nova/agent/agentserver"
	"github.com/lppgo/nova/core"
	"github.com/lppgo/nova/lib/dockerdaemon"
	"github.com/lppgo/nova/lib/dockerregistry"
	"github.com/lppgo/nova/lib/store"
	"github.com/lppgo/nova/lib/torrent/networkevent"
	"github.com/lppgo/nova/lib/torrent/scheduler"
	"github.com/lppgo/nova/lib/upstream"
	"github.com/lppgo/nova/metrics"
	"github.com/lppgo/nova/nginx"
	"github.com/lppgo/nova/utils/httputil"

	"go.uber.org/zap"
)

// Config defines agent configuration.
type Config struct {
	ZapLogging      zap.Config                     `yaml:"zap"`
	Metrics         metrics.Config                 `yaml:"metrics"`
	CADownloadStore store.CADownloadStoreConfig    `yaml:"store"`
	Registry        dockerregistry.Config          `yaml:"registry"`
	Scheduler       scheduler.Config               `yaml:"scheduler"`
	PeerIDFactory   core.PeerIDFactory             `yaml:"peer_id_factory"`
	NetworkEvent    networkevent.Config            `yaml:"network_event"`
	Tracker         upstream.PassiveHashRingConfig `yaml:"tracker"`
	BuildIndex      upstream.PassiveConfig         `yaml:"build_index"`
	AgentServer     agentserver.Config             `yaml:"agentserver"`
	RegistryBackup  string                         `yaml:"registry_backup"`
	Nginx           nginx.Config                   `yaml:"nginx"`
	TLS             httputil.TLSConfig             `yaml:"tls"`
	AllowedCidrs    []string                       `yaml:"allowed_cidrs"`
	DockerDaemon    dockerdaemon.Config            `yaml:"docker_daemon"`
}
