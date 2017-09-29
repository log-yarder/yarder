package main

import (
	"fmt"
	"github.com/log-yarder/yarder/appender"
	"github.com/log-yarder/yarder/storage"
	"io/ioutil"
	"path"
)

const (
	maxEntriesPerChunk = 10
)

func main() {
	tmpDir, err := ioutil.TempDir(path.Join("/", "tmp"), "almanac-dev")
	if err != nil {
		panic(fmt.Sprintf("Unable to create temp dir, %v", err))
	}

	diskStorage := storage.NewDiskStorage(tmpDir)
	appender := appender.New(diskStorage, maxEntriesPerChunk)

	for i := 0; i < 40; i++ {
		appender.HandleRequest(fmt.Sprintf("entry-%d", i))
	}
}
