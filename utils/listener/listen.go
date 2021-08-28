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
package listener

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lppgo/nova/utils/log"
	"github.com/pkg/errors"
)

// Serve serves h on a listener configured by config. Useful for easily
// swapping tcp / unix servers.
func Serve(config Config, h http.Handler) error {
	if config.Net == "unix" {
		if _, err := os.Lstat(config.Addr); !os.IsNotExist(err) {
			err = os.Remove(config.Addr)
			if err != nil {
				err = errors.Wrap(err, "RemoveErr")
				return err
			}
		}
	}
	l, err := net.Listen(config.Net, config.Addr)
	if err != nil {
		return err
	}
	return http.Serve(l, h)
}

//
func ListenAndServe(config Config, h http.Handler) (err error) {
	server := &http.Server{
		Addr:         config.Addr,
		Handler:      h,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	go func(config Config, h http.Handler) {
		serverErr := Serve(config, h)
		if serverErr != nil {
			log.Errorf("serverErr %s:", serverErr.Error())
		}
	}(config, h)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()

	stop()
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(timeoutCtx); err != nil {
		err = errors.WithMessage(err, "utils listener shutdown error")
	}
	return err
}
