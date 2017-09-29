package appender

import (
	"fmt"
	"github.com/log-yarder/yarder/storage"
	"log"
)

func New(storage storage.Storage, maxEntriesPerChunk int) *appender {
	return &appender{
		storage:            storage,
		maxEntriesPerChunk: maxEntriesPerChunk,
		openChunk:          nil,
	}
}

// HandleRequest processes a request to append a single entry to the logs.
func (a *appender) HandleRequest(entry string) error {
	// Make sure we have a chunk to write to.
	if a.openChunk == nil {
		chunk, err := a.storage.CreateChunk()
		if err != nil {
			return fmt.Errorf("Unable to create chunk: %v", err)
		}
		a.openChunk = chunk
		log.Println("Appender opened new chunk")
	}

	// Write the entry.
	err := a.openChunk.Append(entry)
	if err != nil {
		return fmt.Errorf("Unable to append entry: %v", err)
	}

	// Close up the current chunk if necessary.
	if a.openChunk.Size() > a.maxEntriesPerChunk {
		err := a.openChunk.Close()
		if err != nil {
			return fmt.Errorf("Unable to close chunk: %v", err)
		}

		a.openChunk = nil
	}

	return nil
}

type appender struct {
	storage            storage.Storage
	maxEntriesPerChunk int
	openChunk          storage.LogChunk
}
