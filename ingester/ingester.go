package ingester

import (
	"encoding/json"
	"fmt"

	"github.com/log-yarder/yarder/discovery"
	"github.com/log-yarder/yarder/storage"
)

const (
	timestampKey = "timestamp"
)

// Ingester is the entry point for log entry ingestion. It fans the entries out
// to an appropriate number of appenders.
type Ingester struct {
	Discovery *discovery.Discovery
}

// HandleIngest ingests a single log entry. Then entry is treated as a json
// blob. The only requirement is that the blob have a 'timestamp' field.
func (i *Ingester) HandleIngest(rawEntry []byte) error {
	entry, err := createEntry(rawEntry)
	if err != nil {
		return fmt.Errorf("unable to create entry from raw entry: %v", err)
	}

	// For now, just append to all appenders.
	for i, appender := range i.Discovery.Appenders {
		err := appender.HandleAppend(entry)
		if err != nil {
			return fmt.Errorf("unable to handle append on appender %d: %v", i, err)
		}
	}
	return nil
}

func createEntry(rawEntry []byte) (*storage.LogEntry, error) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(rawEntry, &jsonMap)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal json: %v", err)
	}

	value, ok := jsonMap[timestampKey]
	if !ok {
		return nil, fmt.Errorf("could not find key %s in map", timestampKey)
	}

	timestampMs, ok := value.(float64)
	if !ok {
		return nil, fmt.Errorf("unable to interpret timestamp as int64")
	}

	return &storage.LogEntry{
		TimestampMs: int64(timestampMs),
		Raw:         rawEntry,
	}, nil
}
