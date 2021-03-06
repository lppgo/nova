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
package config

// OriginTemplate is the default origin nginx tmpl.
const OriginTemplate = `
server {
  listen {{.port}};

    {{.client_verification}}

  client_max_body_size 10G;
  
  proxy_connect_timeout 300;
  proxy_read_timeout 300;
  proxy_send_timeout 300;

  access_log {{.access_log_path}} json;
  error_log {{.error_log_path}};

  gzip on;
  gzip_types text/plain test/csv application/json;

  location / {
    proxy_pass http://{{.server}};
  }
}
`
