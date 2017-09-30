package discovery

import (
	"github.com/log-yarder/yarder/appender"
)

type Discovery struct {
	Appenders []*appender.Appender
}
