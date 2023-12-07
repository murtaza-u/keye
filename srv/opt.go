package srv

import (
	"fmt"
	"os"
	"time"
)

type optFunc func(*opts) error

type opts struct {
	// port for the database server.
	port string
	// wpi (watcher ping interval) is the duration between two keepalive
	// pings for the watcher.
	wpi time.Duration
	// path to the database file
	path string
	// enable gRPC reflection
	reflect bool
}

// WithPort sets the port for the database server. Ensure that the
// chosen port is available and not in use.
//
// Default: 23023
func WithPort(port uint16) optFunc {
	return func(o *opts) error {
		o.port = fmt.Sprintf(":%d", port)
		return nil
	}
}

// WithWatcherPingInterval sets the interval between two keepalive pings
// for the watcher. Keepalive messages are periodically sent to maintain
// the watcher's connection. Adjust this value if you encounter
// connection issues.
//
// Default: 10s
func WithWatcherPingInterval(interval time.Duration) optFunc {
	return func(o *opts) error {
		o.wpi = interval
		return nil
	}
}

// WithDBPath sets the path to the database file. It throws an error if
// the specified file is a directory or not writable.
//
// Default: ./data.db
func WithDBPath(path string) optFunc {
	return func(o *opts) error {
		info, err := os.Stat(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if info.IsDir() {
			return fmt.Errorf("%q is a directory", path)
		}
		o.path = path
		return nil
	}
}

// WithReflection enables gRPC reflection
//
// Default: disabled
func WithReflection() optFunc {
	return func(o *opts) error {
		o.reflect = true
		return nil
	}
}

func defaultOpts() opts {
	return opts{
		port:    ":23023",
		wpi:     time.Second * 10,
		path:    "data.db",
		reflect: false,
	}
}
