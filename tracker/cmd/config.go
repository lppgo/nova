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
	"go.uber.org/zap"

	"github.com/lppgo/nova/lib/upstream"
	"github.com/lppgo/nova/metrics"
	"github.com/lppgo/nova/nginx"
	"github.com/lppgo/nova/tracker/originstore"
	"github.com/lppgo/nova/tracker/peerhandoutpolicy"
	"github.com/lppgo/nova/tracker/peerstore"
	"github.com/lppgo/nova/tracker/trackerserver"
	"github.com/lppgo/nova/utils/httputil"
)

// Config defines tracker configuration.
type Config struct {
	ZapLogging        zap.Config               `yaml:"zap"`
	PeerStore         peerstore.Config         `yaml:"peerstore"`
	OriginStore       originstore.Config       `yaml:"originstore"`
	TrackerServer     trackerserver.Config     `yaml:"trackerserver"`
	PeerHandoutPolicy peerhandoutpolicy.Config `yaml:"peerhandoutpolicy"`
	Origin            upstream.ActiveConfig    `yaml:"origin"`
	Metrics           metrics.Config           `yaml:"metrics"`
	Nginx             nginx.Config             `yaml:"nginx"`
	TLS               httputil.TLSConfig       `yaml:"tls"`
}
