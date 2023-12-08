package client

import "time"

// Config holds configurations for the Keye client.
type Config struct {
	// Addr is the Keye database server address.
	//
	// Default: localhost:23023
	Addr string
	// Timeout sets the request cancellation time.
	//
	// Default: 5s
	Timeout time.Duration
}

// impute replaces missing values in the configuration with defaults.
func (conf *Config) impute() {
	if conf.Addr == "" {
		conf.Addr = "localhost:23023"
	}

	if conf.Timeout == 0 {
		conf.Timeout = time.Second * 5
	}
}
