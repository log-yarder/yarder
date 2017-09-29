package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

const (
	dirPerms  = 0777
	filePerms = 777
)

// LogChunk represents a sequence of log entries.
type LogChunk interface {
	// Append add a single entry to an open chunk.
	Append(entry string) error

	// Close declares an open chunk as done and performs necessary cleanup tasks.
	Close() error
}

// Storage represents a single machine's interface to a corpus of almanac data.
type Storage interface {
	// CreateChunk allocates a new empty chunk of logs.
	CreateChunk() (LogChunk, error)
}

// NewDiskStorage returns an instance of Strage backed by files on disk.
func NewDiskStorage(rootPath string) Storage {
	log.Printf("Creating new storage with root: %s\n", rootPath)
	return &diskStorage{
		chunkCounter: 0,
		path:         rootPath,
	}
}

// diskLogChunk is an implementation of LogChunk backed by a directory on disk.
type diskLogChunk struct {
	id             string
	chunkDir       string
	chunkFileRaw   string
	chunkFileIndex string
	entryCounter   int
	closed         bool
}

func (c *diskLogChunk) Append(entry string) error {
	entryId := fmt.Sprintf("entry-%d", c.entryCounter)
	c.entryCounter++

	entryPath := path.Join(c.chunkDir, entryId)
	err := ioutil.WriteFile(entryPath, []byte(entry), filePerms)
	if err != nil {
		return fmt.Errorf("Unable to write file: %v", err)
	}

	log.Printf("[%s] wrote [%s] to file: %s\n", c.id, entryId, entryPath)
	return nil
}

// index is the format used for the index file of a chunk.
type index struct {
	Count int
}

// raw is the format used for the raw file of a chunk.
type raw struct {
	Entries []string
}

func (c *diskLogChunk) Close() error {
	if c.closed {
		return fmt.Errorf("Cannot close already closed %s", c.id)
	}
	c.closed = true

	// Go through all the entries and construct the contents of the two files.
	rawEntries := []string{}
	indexCount := 0
	err := filepath.Walk(c.chunkDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// Skip the directory itself.
			return nil
		}

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Unable to read file: %s", path)
		}
		rawEntries = append(rawEntries, string(bytes))
		indexCount++
		return nil
	})
	if err != nil {
		return fmt.Errorf("Unable to walk files in chunk: %v", err)
	}

	// Write the raw file for this chunk before anything else.
	err = writeJson(c.chunkFileRaw, &raw{Entries: rawEntries})
	if err != nil {
		return fmt.Errorf("Unable to write raw chunk: %v", err)
	}
	log.Printf("Wrote raw file: %s", c.chunkFileRaw)

	// Write the index file corresponding to this chunk.
	err = writeJson(c.chunkFileIndex, &index{Count: indexCount})
	if err != nil {
		return fmt.Errorf("Unable to write chunk index: %v", err)
	}
	log.Printf("Wrote index file: %s", c.chunkFileIndex)

	// Delete the directory containing the individual files.
	err = os.RemoveAll(c.chunkDir)
	if err != nil {
		return fmt.Errorf("Unable to delete chunk directory %s: %v", c.chunkDir, err)
	}
	log.Printf("Removed chunk directory %s", c.chunkDir)

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

// diskStorage is an implementation of Storage backed by a directory on disk.
type diskStorage struct {
	chunkCounter int
	path         string
}

func (s *diskStorage) CreateChunk() (LogChunk, error) {
	chunkId := fmt.Sprintf("chunk-%d", s.chunkCounter)
	s.chunkCounter++
	log.Printf("Allocating [%s]\n", chunkId)

	chunkPath := path.Join(s.path, chunkId)
	err := os.Mkdir(chunkPath, dirPerms)
	if err != nil {
		return nil, fmt.Errorf("Unable to create file: %v", err)
	}

	return &diskLogChunk{
		id:             chunkId,
		closed:         false,
		entryCounter:   0,
		chunkDir:       chunkPath,
		chunkFileRaw:   path.Join(s.path, fmt.Sprintf("%s-raw", chunkId)),
		chunkFileIndex: path.Join(s.path, fmt.Sprintf("%s-idx", chunkId)),
	}, nil
}
