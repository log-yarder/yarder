package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
)

const (
	dirPerms  = 0777
	filePerms = 777
)

// LogChunk represents a sequence of log entries.
type LogChunk interface {
	// Append add a single entry to an open chunk.
	Append(entry string) error

	// Size returns the current number of entries in the chunk.
	Size() int

	// Close declares an open chunk as done and performs necessary cleanup tasks.
	Close() error
}

// Storage represents a single machine's interface to a corpus of almanac data.
type Storage interface {
	// CreateChunk allocates a new empty chunk of logs.
	CreateChunk() (LogChunk, error)
}

// DiskStorage is an implementation of Storage backed by a directory on disk.
type DiskStorage struct {
	Path         string
	chunkCounter int
}

// persistedChunk is the format used to serialize a chunk to bytes.
type persistedChunk struct {
	Entries []string
}

// diskLogChunk is an implementation of LogChunk backed by a directory on disk.
type diskLogChunk struct {
	id        string
	closed    bool
	chunkFile string
	entries   []string
}

func (c *diskLogChunk) Append(entry string) error {
	c.entries = append(c.entries, entry)
	return nil
}

func (c *diskLogChunk) Size() int {
	return len(c.entries)
}

func (c *diskLogChunk) Close() error {
	if c.closed {
		return fmt.Errorf("Cannot close already closed %s", c.id)
	}
	c.closed = true

	// Write the raw file for this chunk before anything else.
	err := writeJson(c.chunkFile, &persistedChunk{Entries: c.entries})
	if err != nil {
		return fmt.Errorf("Unable to write raw chunk: %v", err)
	}
	log.Printf("Wrote chunk file: %s", c.chunkFile)

	return nil
}

func writeJson(path string, content interface{}) error {
	bytes, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("Unable to marshal json: %v", err)
	}

	err = ioutil.WriteFile(path, bytes, filePerms)
	if err != nil {
		return fmt.Errorf("Unable to write file to file %s: %v", path, err)
	}

	return nil
}

func (s *DiskStorage) CreateChunk() (LogChunk, error) {
	chunkId := fmt.Sprintf("chunk-%d", s.chunkCounter)
	s.chunkCounter++

	return &diskLogChunk{
		id:        chunkId,
		closed:    false,
		chunkFile: path.Join(s.Path, chunkId),
		entries:   []string{},
	}, nil
}
