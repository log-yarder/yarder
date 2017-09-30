package main

import (
	"fmt"
	"github.com/log-yarder/yarder/appender"
	"github.com/log-yarder/yarder/storage"
	"io/ioutil"
	"log"
)

const (
	maxEntriesPerChunk = 10
)

func main() {
	tmpDir, err := ioutil.TempDir("", "yarder-dev")
	if err != nil {
		log.Panicf(fmt.Sprintf("Unable to create temp dir, %v", err))
	}

	diskStorage := &storage.DiskStorage{Path: tmpDir}
	appender := &appender.Appender{
		Storage:            diskStorage,
		MaxEntriesPerChunk: maxEntriesPerChunk,
	}

	for i := 0; i < 40; i++ {
		appender.HandleRequest(fmt.Sprintf("entry-%d", i))
	}
}
