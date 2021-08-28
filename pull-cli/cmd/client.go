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
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/lppgo/nova/utils/log"

	"github.com/docker/go-units"

	"github.com/flylog/colorstyle"

	"github.com/pkg/errors"

	"github.com/lppgo/nova/core"
	"github.com/lppgo/nova/utils/httputil"
)

// Client errors.
var (
	ErrTagNotFound = errors.New("tag not found")
)

// Client defines a client for accessing the agent server.
type Client interface {
	GetTag(tag string) (core.Digest, error)
	Download(namespace string, d core.Digest) (io.ReadCloser, error)
}

// HTTPClient provides a wrapper for HTTP operations on an agent.
type HTTPClient struct {
	addr string
}

// New creates a new client for an agent at addr.
func New(addr string) *HTTPClient {
	return &HTTPClient{addr}
}

// GetTag resolves tag into a digest. Returns ErrTagNotFound if the tag does
// not exist.
func (c *HTTPClient) GetDigestByTag(tag string) (core.Digest, error) {
	resp, err := httputil.Get(fmt.Sprintf("http://%s/tags/%s", c.addr, url.PathEscape(tag)))
	if err != nil {
		if httputil.IsNotFound(err) {
			return core.Digest{}, ErrTagNotFound
		}
		return core.Digest{}, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return core.Digest{}, fmt.Errorf("read body: %s", err)
	}
	d, err := core.ParseSHA256Digest(string(b))
	if err != nil {
		return core.Digest{}, fmt.Errorf("parse digest: %s", err)
	}
	return d, nil
}

// Download returns the blob of d. Callers should close the returned ReadCloser
// when done reading the blob.
func (c *HTTPClient) Download(namespace string, d core.Digest) (io.ReadCloser, error) {
	resp, err := httputil.Get(
		fmt.Sprintf("http://%s/namespace/%s/blobs/%s", c.addr, url.PathEscape(namespace), d))
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// PullFile pull file and save to local.
func PullFile(addr, tag, namespace, filename, pullDir string) error {
	startAt := time.Now()
	client := New(addr)
	digest, err := client.GetDigestByTag(tag)
	if err != nil {
		log.Errorf("GetDigestByTag Error:%s", err.Error())
		return err
	}
	rc, err := client.Download(namespace, digest)
	if err != nil {
		log.Errorf("Download Error:%s", err.Error())
		return err
	}
	defer rc.Close()

	filePath := path.Join(pullDir, filename)

	_, err = os.Stat(filePath)
	// exist
	if err == nil {
		newName, err := GenerateNewFileName(filename)
		if err != nil {
			return errors.Wrap(err, "GenerateNewFileName Error")
		}
		backUpFilePath := path.Join(pullDir, newName)
		log.Infof("backUpFilePath:%s", backUpFilePath)
		err = os.Rename(filePath, backUpFilePath)
		if err != nil {
			return errors.Wrapf(err, "rename [%s] > [%s] Error", filePath, backUpFilePath)
		}
	}

	if err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "Stat filePath :%s Error", filePath)
	}

	dstf, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "create file ")
	}
	defer dstf.Close()

	writeCounter := &WriteCounter{ObjPath: filePath}
	if _, err = io.Copy(dstf, io.TeeReader(rc, writeCounter)); err != nil {
		return errors.Wrapf(err, "Copy to dstFile Error")
	}

	fmt.Println()
	log.Infof("pull file [%s] is succeeded ;耗时[%s] ! ", colorstyle.Yellow(filePath), colorstyle.Yellow(time.Since(startAt).String()))
	return nil
}

// GenerateNewFileName 返回新文件的name
func GenerateNewFileName(oldName string) (newName string, err error) {
	base := 10
	exceptedLen := 2
	suffix := "_" + strconv.FormatInt(time.Now().UnixMilli(), base)

	extensionName := path.Ext(oldName)
	if extensionName == "" {
		newName = oldName + suffix
		return newName, nil
	}

	names := strings.Split(oldName, extensionName)
	if len(names) != exceptedLen {
		return "", errors.New(fmt.Sprintf("filename %s is invalid", oldName))
	}

	newName = names[0] + suffix + extensionName
	return newName, nil
}

type WriteCounter struct {
	Total   uint64
	ObjPath string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 5))
	fmt.Printf(colorstyle.Green(fmt.Sprintf("\r======> Downloading [%s]... %s complete <======", wc.ObjPath, units.HumanSize(float64(wc.Total)))))
}
