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
package healthcheck

import (
	"testing"
	"time"

	"github.com/lppgo/nova/lib/hostlist"
	mockhealthcheck "github.com/lppgo/nova/mocks/lib/healthcheck"
	"github.com/lppgo/nova/utils/stringset"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestActiveMonitor(t *testing.T) {
	require := require.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	x := "x:80"
	y := "y:80"

	filter := mockhealthcheck.NewMockFilter(ctrl)

	filter.EXPECT().Run(stringset.New(x, y)).Return(stringset.New(x))

	m := NewMonitor(
		MonitorConfig{Interval: time.Second},
		hostlist.Fixture(x, y),
		filter)
	defer m.Stop()

	require.Equal(stringset.New(x, y), m.Resolve())

	time.Sleep(1250 * time.Millisecond)

	require.Equal(stringset.New(x), m.Resolve())
}
