package appender

import (
	"fmt"

	"github.com/log-yarder/yarder/storage"
)

// Appender handles requets to append log entries to storage-backed chunks.
type Appender struct {
	Storage            storage.Storage
	MaxEntriesPerChunk int
	openChunk          storage.LogChunk
}

// HandleAppend processes a request to append a single entry to the logs.
func (a *Appender) HandleAppend(entry *storage.LogEntry) error {
	// Make sure we have a chunk to write to.
	if a.openChunk == nil {
		chunk, err := a.Storage.CreateChunk()
		if err != nil {
			return fmt.Errorf("Unable to create chunk: %v", err)
		}
		a.openChunk = chunk
	}

	// Write the entry.
	err := a.openChunk.Append(entry)
	if err != nil {
		return fmt.Errorf("Unable to append entry: %v", err)
	}

	// Close up the current chunk if necessary.
	if a.openChunk.Size() > a.MaxEntriesPerChunk {
		err := a.openChunk.Close()
		if err != nil {
			return fmt.Errorf("Unable to close chunk: %v", err)
		}

		a.openChunk = nil
	}

	return nil
}
