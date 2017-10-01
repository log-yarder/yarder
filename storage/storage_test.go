package storage

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSize(t *testing.T) {
	storage := createStorage(t)
	chunk, err := storage.CreateChunk()
	require.NoError(t, err)

	chunk.Append(&LogEntry{Timestamp: time.Unix(6, 0)})
	chunk.Append(&LogEntry{Timestamp: time.Unix(4, 0)})
	chunk.Append(&LogEntry{Timestamp: time.Unix(5, 0)})

	require.Equal(t, 3, chunk.Size())
}

func TestSortsEntries(t *testing.T) {
	storage := createStorage(t)
	chunk, err := storage.CreateChunk()
	require.NoError(t, err)

	chunk.Append(&LogEntry{Timestamp: time.Unix(6, 0)})
	chunk.Append(&LogEntry{Timestamp: time.Unix(4, 0)})
	chunk.Append(&LogEntry{Timestamp: time.Unix(5, 0)})

	err = chunk.Close()
	require.NoError(t, err)

	// TODO(dino): Add an actual read API to storage. For now, we use the fact
	// that we are dealing with a file storage and use its path.
	files, err := ioutil.ReadDir(storage.Path)
	require.NoError(t, err)
	require.Equal(t, 1, len(files))

	contents, err := ioutil.ReadFile(path.Join(storage.Path, files[0].Name()))
	require.NoError(t, err)

	var p persistedChunk
	err = json.Unmarshal(contents, &p)
	require.NoError(t, err)

	// Check that the entries are in order.
	require.Equal(t, 3, len(p.Entries))
	require.Equal(t, time.Unix(4, 0), p.Entries[0].Timestamp)
	require.Equal(t, time.Unix(5, 0), p.Entries[1].Timestamp)
	require.Equal(t, time.Unix(6, 0), p.Entries[2].Timestamp)
}

func createStorage(t *testing.T) *DiskStorage {
	tmpDir, err := ioutil.TempDir("", "yarder-dev")
	require.NoError(t, err)
	return &DiskStorage{Path: tmpDir}
}
