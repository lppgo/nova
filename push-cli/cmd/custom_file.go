/**
 * @Author: lucas
 * @Description:
 * @File:  custom_file.go
 * @Version: 1.0.0
 * @Date: 2021/8/20 15:29
 */
package cmd

import (
	"os"

	"github.com/lppgo/nova/lib/store/base"
)

/*
type FileReader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
	Size() int64
}
*/
type CustomFile struct {
	*os.File
}

type CustomFiler = base.FileReader

func (c CustomFile) Size() int64 {
	info, _ := c.Stat()
	return info.Size()
}
