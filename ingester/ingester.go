package ingester

import (
	"fmt"
	"github.com/log-yarder/yarder/discovery"
)

// Ingester is the entry point for log entry ingestion. It fans the entries out
// to an appropriate number of appenders.
type Ingester struct {
	Discovery *discovery.Discovery
}

// HandleIngest ingests a single log entry.
func (i *Ingester) HandleIngest(entry string) error {
	// For now, just append to all appenders.
	for _, appender := range i.Discovery.Appenders {
		err := appender.HandleAppend(entry)
		if err != nil {
			return fmt.Errorf("Unable to handle append: %v", err)
		}
	}
	return nil
}
