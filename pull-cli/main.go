/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import "github.com/lppgo/nova/pull-cli/cmd"

func main() {
	cmd.Execute()
}

/*
// NOTE 下载文件脚本描述

./pull-cli \
--agent-server-port=16002 \
--tag=sse1.zip:latest \
--download-dir=/mnt/e/sse/

./pull-cli.exe \
--agent-server-port=18002 \
--tag=sse1.zip:latest \
--download-dir=E:\\sse
*/
