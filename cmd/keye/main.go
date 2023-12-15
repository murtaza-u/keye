package main

import (
	"log"
	"os"
	"time"

	"github.com/murtaza-u/keye/srv"

	"github.com/urfave/cli/v2"
)

const version = "23.12"

func main() {
	app := cli.NewApp()
	app.Name = "keye"
	app.Usage = "key-value DB with the ability to watch over keys"
	app.Version = version
	app.EnableBashCompletion = true
	app.Copyright = "Apache-2.0"
	app.Authors = []*cli.Author{
		{Name: "Murtaza Udaipurwala", Email: "murtaza@murtazau.xyz"},
	}
	app.Flags = []cli.Flag{
		&cli.PathFlag{
			Name:  "db-file",
			Usage: "path to the database file",
			Value: "data.db",
		},
		&cli.UintFlag{
			Name:  "port",
			Usage: "port for the database server.",
			Value: 23023,
		},
		&cli.DurationFlag{
			Name:  "watcher-ping-interval",
			Usage: "duration between two keepalive pings for the watcher",
			Value: time.Second * 10,
		},
		&cli.UintFlag{
			Name:  "event-queue-size",
			Usage: "size of the event queue",
			Value: 10,
		},
		&cli.BoolFlag{
			Name:  "enable-reflection",
			Usage: "enable gRPC reflection",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "json-logger",
			Usage: "use JSON logger",
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug logs",
			Value: false,
		},
	}
	app.Action = func(ctx *cli.Context) error {
		vars, err := parseEnv()
		if err != nil {
			log.Fatal(err)
		}
		vars.overFlags(ctx)
		srv, err := srv.New(vars.toOpts()...)
		if err != nil {
			log.Fatal(err)
		}
		return srv.Run()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
