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

import "github.com/lppgo/nova/agent/cmd"

func main() {
	cmd.Run(cmd.ParseFlags())
}

/*
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
  --agent-server-port=18002
*/
