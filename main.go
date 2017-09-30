package main

import (
	"fmt"
	"github.com/log-yarder/yarder/appender"
	"github.com/log-yarder/yarder/discovery"
	"github.com/log-yarder/yarder/ingester"
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
	appenders := []*appender.Appender{
		&appender.Appender{
			Storage:            diskStorage,
			MaxEntriesPerChunk: maxEntriesPerChunk,
		},
		&appender.Appender{
			Storage:            diskStorage,
			MaxEntriesPerChunk: maxEntriesPerChunk,
		},
		&appender.Appender{
			Storage:            diskStorage,
			MaxEntriesPerChunk: maxEntriesPerChunk,
		},
	}
	discovery := &discovery.Discovery{Appenders: appenders}
	ingester := &ingester.Ingester{Discovery: discovery}

	for i := 0; i < 40; i++ {
		err := ingester.HandleIngest(fmt.Sprintf("entry-%d", i))
		if err != nil {
			log.Panicf(fmt.Sprintf("Failed to handle ingest, %v", err))
		}
	}
}
