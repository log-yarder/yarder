package main

import (
	"fmt"
	"github.com/log-yarder/yarder/appender"
	"github.com/log-yarder/yarder/storage"
	"io/ioutil"
	"path"
)

func main() {
	tmpDir, err := ioutil.TempDir(path.Join("/", "tmp"), "almanac-dev")
	if err != nil {
		panic(fmt.Sprintf("Unable to create temp dir, %v", err))
	}

	diskStorage := storage.NewDiskStorage(tmpDir)
	appender := appender.New(diskStorage)

	chunk, err := diskStorage.CreateChunk()
	if err != nil {
		panic(fmt.Sprintf("Unable to create chunk, %v", err))
	}

	err = chunk.Append("foo")
	if err != nil {
		panic(fmt.Sprintf("Unable to write first log entry, %v", err))
	}

	err = chunk.Append("bar")
	if err != nil {
		panic(fmt.Sprintf("Unable to write second log entry, %v", err))
	}

	err = chunk.Close()
	if err != nil {
		panic(fmt.Sprintf("Unable to close log chunk, %v", err))
	}

	appender.HandleRequest("fooentry")
}
