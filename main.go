package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/log-yarder/yarder/appender"
	"github.com/log-yarder/yarder/discovery"
	"github.com/log-yarder/yarder/ingester"
	"github.com/log-yarder/yarder/storage"
)

const (
	maxEntriesPerChunk = 10
)

func main() {
	tmpDir, err := ioutil.TempDir("", "yarder-dev")
	if err != nil {
		log.Panicf("unable to create temp dir, %v", err)
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
		name := fmt.Sprintf("entry-%d", i)
		entryBlob, err := json.Marshal(&entry{
			Name:      name,
			Timestamp: time.Now().Unix(),
		})
		if err != nil {
			log.Panicf("failed to marshal json: %v", err)
		}

		err = ingester.HandleIngest(entryBlob)
		if err != nil {
			log.Panicf("failed to ingest entry [%s], %v", name, err)
		}
	}
}

// entry is used only to serialize test log entries to json.
type entry struct {
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}
