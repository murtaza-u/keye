package srv

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Delta456/box-cli-maker/v2"
)

type OptFunc func(*opts) error

type opts struct {
	// port for the database server.
	port string
	// wpi (watcher ping interval) is the duration between two keepalive
	// pings for the watcher.
	wpi time.Duration
	// size of the event queue
	eventQueueSize uint
	// path to the database file
	path string
	// enable gRPC reflection
	reflect bool
	// log in JSON
	jsonLogger bool
	// enable debug logs
	debug bool
}

// WithPort sets the port for the database server. Ensure that the
// chosen port is available and not in use.
//
// Default: 23023
func WithPort(port uint16) OptFunc {
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
func WithWatcherPingInterval(interval time.Duration) OptFunc {
	return func(o *opts) error {
		o.wpi = interval
		return nil
	}
}

// WithEventQueueSize sets the event queue size used by the watcher to
// notify subscribers about events in their occurrence order.
//
// Default size: 10
func WithEventQueueSize(size uint) OptFunc {
	return func(o *opts) error {
		o.eventQueueSize = size
		return nil
	}
}

// WithDBPath sets the path to the database file. It throws an error if
// the specified file is a directory or not writable.
//
// Default: ./data.db
func WithDBPath(path string) OptFunc {
	return func(o *opts) error {
		info, err := os.Stat(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if os.IsNotExist(err) {
			o.path = path
			return nil
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
func WithReflection() OptFunc {
	return func(o *opts) error {
		o.reflect = true
		return nil
	}
}

// WithJSONLogger configures logger to use JSON.
//
// Default: disabled
func WithJSONLogger() OptFunc {
	return func(o *opts) error {
		o.jsonLogger = true
		return nil
	}
}

// WithDebug enables debug logs.
//
// Default: disabled
func WithDebug() OptFunc {
	return func(o *opts) error {
		o.debug = true
		return nil
	}
}

// print displays server options in a box format
func (o opts) print() {
	box := box.New(box.Config{
		Px:            2,
		Py:            2,
		Type:          "Bold",
		AllowWrapping: true,
		ContentAlign:  "Left",
		TitlePos:      "Bottom",
	})
	params := strings.Join([]string{
		fmt.Sprintf("Address=%s", o.port),
		fmt.Sprintf("DB=%s", o.path),
		fmt.Sprintf("WPI=%s", o.wpi),
		fmt.Sprintf("Event Queue Size=%d", o.eventQueueSize),
		fmt.Sprintf("Reflection=%v", o.reflect),
		fmt.Sprintf("JSON logger=%v", o.jsonLogger),
		fmt.Sprintf("Debug=%v", o.debug),
	}, "\n")
	box.Print("Keye", params)
}

func defaultOpts() opts {
	return opts{
		port:           ":23023",
		wpi:            time.Second * 10,
		eventQueueSize: 10,
		path:           "data.db",
		reflect:        false,
		jsonLogger:     false,
		debug:          false,
	}
}
