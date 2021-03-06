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
package peerhandoutpolicy

import "github.com/lppgo/nova/core"

const _defaultPolicy = "default"

// defaultAssignmentPolicy is a NO-OP policy that assigns all peers
// the highest priority.
type defaultAssignmentPolicy struct{}

func newDefaultAssignmentPolicy() assignmentPolicy {
	return &defaultAssignmentPolicy{}
}

func (p *defaultAssignmentPolicy) assignPriority(peer *core.PeerInfo) (int, string) {
	return 0, "default"
}
