package ingester

import (
	"fmt"
	"github.com/log-yarder/yarder/discovery"
)

type Ingester struct {
	Discovery *discovery.Discovery
}

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
