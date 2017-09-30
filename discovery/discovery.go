package discovery

import (
	"github.com/log-yarder/yarder/appender"
)

// Discovery provides access to the available services.
type Discovery struct {
	Appenders []*appender.Appender
}
