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
	"github.com/lppgo/nova/lib/dockerregistry"
	"github.com/lppgo/nova/lib/store"
	"github.com/lppgo/nova/lib/upstream"
	"github.com/lppgo/nova/metrics"
	"github.com/lppgo/nova/nginx"
	"github.com/lppgo/nova/push-cli/registryoverride"
	"github.com/lppgo/nova/utils/httputil"

	"go.uber.org/zap"
)

// Config defines push-cli configuration
type Config struct {
	CAStore          store.CAStoreConfig     `yaml:"castore"`
	Registry         dockerregistry.Config   `yaml:"registry"`
	BuildIndex       upstream.ActiveConfig   `yaml:"build_index"`
	Origin           upstream.ActiveConfig   `yaml:"origin"`
	ZapLogging       zap.Config              `yaml:"zap"`
	Metrics          metrics.Config          `yaml:"metrics"`
	RegistryOverride registryoverride.Config `yaml:"registryoverride"`
	Nginx            nginx.Config            `yaml:"nginx"`
	TLS              httputil.TLSConfig      `yaml:"tls"`
	Namespace        string                  `yaml:"namespace"`
}
