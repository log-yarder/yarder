package ingester

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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

// HandleIngest ingests a single log entry.
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
		return nil, fmt.Errorf("Unable to unmarshal json: %v", err)
	}

	value, ok := jsonMap[timestampKey]
	if !ok {
		return nil, fmt.Errorf("Could not find key %s in map", timestampKey)
	}

	timeString, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("Unable to interpret timestamp as string")
	}

	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert %s to a timestamp", timeString)
	}

	return &storage.LogEntry{
		Timestamp: time.Unix(timestamp, 0),
		Raw:       rawEntry,
	}, nil
}
