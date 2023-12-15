package main

import (
	"time"

	"github.com/murtaza-u/keye/srv"

	"github.com/caarlos0/env/v10"
	"github.com/urfave/cli/v2"
)

func parseEnv() (*vars, error) {
	v := new(vars)
	if err := env.Parse(v); err != nil {
		return nil, err
	}
	return v, nil
}

type vars struct {
	Path        string        `env:"KEYE_DB_FILE_PATH,expand"`
	Port        uint16        `env:"KEYE_PORT"`
	WPI         time.Duration `env:"KEYE_WATCHER_PING_INTERVAL"`
	EvQueueSize uint          `env:"KEYE_EVENT_QUEUE_SIZE"`
	Reflect     bool          `env:"KEYE_ENABLE_REFLECTION"`
	JSONLogger  bool          `env:"KEYE_USE_JSON_LOGGER"`
	Debug       bool          `env:"KEYE_DEBUG"`
}

func (v *vars) overFlags(ctx *cli.Context) {
	for _, f := range ctx.FlagNames() {
		switch f {
		case "db-file":
			v.Path = ctx.Path(f)
		case "port":
			v.Port = uint16(ctx.Uint(f))
		case "watcher-ping-interval":
			v.WPI = ctx.Duration(f)
		case "event-queue-size":
			v.EvQueueSize = ctx.Uint(f)
		case "enable-reflection":
			v.Reflect = ctx.Bool(f)
		case "json-logger":
			v.JSONLogger = ctx.Bool(f)
		case "debug":
			v.Debug = ctx.Bool(f)
		}
	}
}

func (v vars) toOpts() []srv.OptFunc {
	var opts []srv.OptFunc
	if v.Path != "" {
		opts = append(opts, srv.WithDBPath(v.Path))
	}
	if v.Port != 0 {
		opts = append(opts, srv.WithPort(v.Port))
	}
	if v.WPI != 0 {
		opts = append(opts, srv.WithWatcherPingInterval(v.WPI))
	}
	if v.EvQueueSize != 0 {
		opts = append(opts, srv.WithEventQueueSize(v.EvQueueSize))
	}
	if v.Reflect {
		opts = append(opts, srv.WithReflection())
	}
	if v.JSONLogger {
		opts = append(opts, srv.WithJSONLogger())
	}
	if v.Debug {
		opts = append(opts, srv.WithDebug())
	}
	return opts
}
