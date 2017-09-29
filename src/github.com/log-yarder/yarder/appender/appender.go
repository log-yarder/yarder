package appender

import (
	"fmt"
	"github.com/log-yarder/yarder/storage"
	"log"
)

const (
	maxEntriesPerChunk = 1000
	maxChunkAgeMs      = 1000 * 60 * 30
)

func New(storage storage.Storage) *appender {
	return &appender{
		storage: storage,
	}
}

// HandleRequest processes a request to append a single entry to the logs.
func (a *appender) HandleRequest(entry string) error {
	log.Println("Appender handling request")

	// Make sure we have a chunk to write to.
	if a.openChunk == nil {
		chunk, err := a.storage.CreateChunk()
		if err != nil {
			return fmt.Errorf("Unable to create chunk: %v", err)
		}
		a.openChunk = chunk
		log.Println("Appender opened new chunk")
	}

	err := a.openChunk.Append(entry)
	if err != nil {
		return fmt.Errorf("Unable to append entry: %v", err)
	}
	return nil
}

type appender struct {
	storage   storage.Storage
	openChunk storage.LogChunk
}
