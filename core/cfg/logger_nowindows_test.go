//go:build !windows

package cfg

import (
	"errors"
	"log/syslog"
	"testing"

	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

func TestInitLoggersDoesNotKeepFailedSyslogHook(t *testing.T) {
	previousHook := syslogHook
	previousFactory := newSyslogHook
	t.Cleanup(func() {
		syslogHook = previousHook
		newSyslogHook = previousFactory
	})

	syslogHook = nil
	newSyslogHook = func(network, raddr string, priority syslog.Priority, tag string) (*logrus_syslog.SyslogHook, error) {
		return &logrus_syslog.SyslogHook{}, errors.New("syslog unavailable")
	}

	InitLoggers("syslog")

	if syslogHook != nil {
		t.Fatal("failed syslog hook must not be published")
	}
	logger := NewLogger("after-failed-syslog-init")
	if len(logger.Hooks) != 0 {
		t.Fatalf("new logger has %d hooks after failed syslog initialization, want 0", len(logger.Hooks))
	}
	logger.Warning("logging must fall back without panicking")
}
