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
package main

import (
	"github.com/lppgo/nova/build-index/cmd"

	// Import all backend client packages to register them with backend manager.
	_ "github.com/lppgo/nova/lib/backend/gcsbackend"
	_ "github.com/lppgo/nova/lib/backend/hdfsbackend"
	_ "github.com/lppgo/nova/lib/backend/httpbackend"
	_ "github.com/lppgo/nova/lib/backend/registrybackend"
	_ "github.com/lppgo/nova/lib/backend/s3backend"
	_ "github.com/lppgo/nova/lib/backend/shadowbackend"
	_ "github.com/lppgo/nova/lib/backend/sqlbackend"
	_ "github.com/lppgo/nova/lib/backend/testfs"
)

func main() {
	cmd.Run(cmd.ParseFlags())
}
