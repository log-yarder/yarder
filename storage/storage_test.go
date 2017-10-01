package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSize(t *testing.T) {
	storage := createStorage(t)
	defer os.RemoveAll(storage.Path)

	chunk, err := storage.CreateChunk()
	require.NoError(t, err)

	chunk.Append(&LogEntry{TimestampMs: 6})
	chunk.Append(&LogEntry{TimestampMs: 4})
	chunk.Append(&LogEntry{TimestampMs: 5})

	require.Equal(t, 3, chunk.Size())
}

func TestSortsEntries(t *testing.T) {
	storage := createStorage(t)
	defer os.RemoveAll(storage.Path)

	chunk, err := storage.CreateChunk()
	require.NoError(t, err)

	chunk.Append(&LogEntry{TimestampMs: 6})
	chunk.Append(&LogEntry{TimestampMs: 4})
	chunk.Append(&LogEntry{TimestampMs: 5})

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
	require.Equal(t, int64(4), p.Entries[0].TimestampMs)
	require.Equal(t, int64(5), p.Entries[1].TimestampMs)
	require.Equal(t, int64(6), p.Entries[2].TimestampMs)
}

func createStorage(t *testing.T) *DiskStorage {
	tmpDir, err := ioutil.TempDir("", "yarder-dev")
	require.NoError(t, err)
	return &DiskStorage{Path: tmpDir}
}
