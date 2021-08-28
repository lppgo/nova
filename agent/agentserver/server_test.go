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
package agentserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"
	"time"

	"github.com/lppgo/nova/agent/agentclient"
	"github.com/lppgo/nova/build-index/tagclient"
	"github.com/lppgo/nova/core"
	"github.com/lppgo/nova/lib/store"
	"github.com/lppgo/nova/lib/torrent/scheduler"
	"github.com/lppgo/nova/lib/torrent/scheduler/connstate"
	mocktagclient "github.com/lppgo/nova/mocks/build-index/tagclient"
	mockdockerdaemon "github.com/lppgo/nova/mocks/lib/dockerdaemon"
	mockscheduler "github.com/lppgo/nova/mocks/lib/torrent/scheduler"
	"github.com/lppgo/nova/utils/httputil"
	"github.com/lppgo/nova/utils/log"
	"github.com/lppgo/nova/utils/testutil"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
)

type serverMocks struct {
	cads      *store.CADownloadStore
	sched     *mockscheduler.MockReloadableScheduler
	tags      *mocktagclient.MockClient
	dockerCli *mockdockerdaemon.MockDockerClient
	cleanup   *testutil.Cleanup
}

func newServerMocks(t *testing.T) (*serverMocks, func()) {
	var cleanup testutil.Cleanup

	cads, c := store.CADownloadStoreFixture()
	cleanup.Add(c)

	ctrl := gomock.NewController(t)
	cleanup.Add(ctrl.Finish)

	sched := mockscheduler.NewMockReloadableScheduler(ctrl)

	tags := mocktagclient.NewMockClient(ctrl)

	dockerCli := mockdockerdaemon.NewMockDockerClient(ctrl)

	return &serverMocks{cads, sched, tags, dockerCli, &cleanup}, cleanup.Run
}

func (m *serverMocks) startServer() string {
	// s := New(Config{}, tally.NoopScope, m.cads, m.sched, m.tags, m.dockerCli)
	s := New(Config{}, tally.NoopScope, m.cads, m.sched, m.tags)
	addr, stop := testutil.StartServer(s.Handler())
	m.cleanup.Add(stop)
	return addr
}

func TestGetTag(t *testing.T) {
	require := require.New(t)

	mocks, cleanup := newServerMocks(t)
	defer cleanup()

	tag := core.TagFixture()
	d := core.DigestFixture()

	mocks.tags.EXPECT().Get(tag).Return(d, nil)

	c := agentclient.New(mocks.startServer())

	result, err := c.GetTag(tag)
	log.Infof("%v\n", result)
	require.NoError(err)
	require.Equal(d, result)
}

func TestGetTagNotFound(t *testing.T) {
	require := require.New(t)

	mocks, cleanup := newServerMocks(t)
	defer cleanup()

	tag := core.TagFixture()

	mocks.tags.EXPECT().Get(tag).Return(core.Digest{}, tagclient.ErrTagNotFound)

	c := agentclient.New(mocks.startServer())

	_, err := c.GetTag(tag)
	// log.Infof("%v\n", d)
	require.Error(err)
	require.Equal(agentclient.ErrTagNotFound, err)
}

func TestDownload(t *testing.T) {
	require := require.New(t)

	mocks, cleanup := newServerMocks(t)
	defer cleanup()

	namespace := core.TagFixture()
	blob := core.NewBlobFixture()

	mocks.sched.EXPECT().Download(namespace, blob.Digest).DoAndReturn(
		func(namespace string, d core.Digest) error {
			return store.RunDownload(mocks.cads, d, blob.Content)
		})

	addr := mocks.startServer()
	c := agentclient.New(addr)

	r, err := c.Download(namespace, blob.Digest)
	require.NoError(err)
	result, err := ioutil.ReadAll(r)
	require.NoError(err)
	require.Equal(string(blob.Content), string(result))
}

func TestDownloadNotFound(t *testing.T) {
	require := require.New(t)

	mocks, cleanup := newServerMocks(t)
	defer cleanup()

	namespace := core.TagFixture()
	blob := core.NewBlobFixture()

	mocks.sched.EXPECT().Download(namespace, blob.Digest).Return(scheduler.ErrTorrentNotFound)

	addr := mocks.startServer()
	c := agentclient.New(addr)

	_, err := c.Download(namespace, blob.Digest)
	require.Error(err)
	require.True(httputil.IsNotFound(err))
}

func TestDownloadUnknownError(t *testing.T) {
	require := require.New(t)

	mocks, cleanup := newServerMocks(t)
	defer cleanup()

	namespace := core.TagFixture()
	blob := core.NewBlobFixture()

	mocks.sched.EXPECT().Download(namespace, blob.Digest).Return(fmt.Errorf("test error"))

	addr := mocks.startServer()
	c := agentclient.New(addr)

	_, err := c.Download(namespace, blob.Digest)
	require.Error(err)
	require.True(httputil.IsStatus(err, 500))
}

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		desc     string
		probeErr error
	}{
		{"probe error", errors.New("some probe error")},
		{"healthy", nil},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			require := require.New(t)

			mocks, cleanup := newServerMocks(t)
			defer cleanup()

			mocks.sched.EXPECT().Probe().Return(test.probeErr)

			addr := mocks.startServer()

			_, err := httputil.Get(fmt.Sprintf("http://%s/health", addr))
			if test.probeErr != nil {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}

func TestPatchSchedulerConfigHandler(t *testing.T) {
	require := require.New(t)

	mocks, cleanup := newServerMocks(t)
	defer cleanup()

	addr := mocks.startServer()

	config := scheduler.Config{
		ConnTTI: time.Minute,
	}
	b, err := json.Marshal(config)
	require.NoError(err)

	mocks.sched.EXPECT().Reload(config)

	_, err = httputil.Patch(
		fmt.Sprintf("http://%s/x/config/scheduler", addr),
		httputil.SendBody(bytes.NewReader(b)))
	require.NoError(err)
}

func TestGetBlacklistHandler(t *testing.T) {
	require := require.New(t)

	mocks, cleanup := newServerMocks(t)
	defer cleanup()

	blacklist := []connstate.BlacklistedConn{{
		PeerID:    core.PeerIDFixture(),
		InfoHash:  core.InfoHashFixture(),
		Remaining: time.Second,
	}}
	mocks.sched.EXPECT().BlacklistSnapshot().Return(blacklist, nil)

	addr := mocks.startServer()

	resp, err := httputil.Get(fmt.Sprintf("http://%s/x/blacklist", addr))
	require.NoError(err)

	var result []connstate.BlacklistedConn
	require.NoError(json.NewDecoder(resp.Body).Decode(&result))
	require.Equal(blacklist, result)
}

func TestDeleteBlobHandler(t *testing.T) {
	require := require.New(t)

	mocks, cleanup := newServerMocks(t)
	defer cleanup()

	d := core.DigestFixture()

	addr := mocks.startServer()

	mocks.sched.EXPECT().RemoveTorrent(d).Return(nil)

	_, err := httputil.Delete(fmt.Sprintf("http://%s/blobs/%s", addr, d))
	require.NoError(err)
}

func TestPreloadHandler(t *testing.T) {
	require := require.New(t)

	mocks, cleanup := newServerMocks(t)
	defer cleanup()

	addr := mocks.startServer()

	tag := url.PathEscape("repo1:tag1")

	mocks.dockerCli.EXPECT().PullImage(context.Background(), "repo1", "tag1").Return(nil)

	_, err := httputil.Get(fmt.Sprintf("http://%s/preload/tags/%s", addr, tag))
	require.NoError(err)
}
