package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

const version = "23.12"

func init() {
	log.SetFlags(0)
}

func main() {
	app := cli.NewApp()
	app.Name = "keyectl"
	app.Usage = "command line tool for interacting with the Keye server"
	app.Version = version
	app.EnableBashCompletion = true
	app.Copyright = "Apache-2.0"
	app.Authors = []*cli.Author{
		{Name: "Murtaza Udaipurwala", Email: "murtaza@murtazau.xyz"},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Usage:   "address of the database server",
			Value:   "localhost:23023",
			Aliases: []string{"endpoint"},
		},
		&cli.DurationFlag{
			Name:  "timeout",
			Usage: "request timeout",
			Value: time.Second * 5,
		},
	}
	app.Commands = []*cli.Command{
		getCmd, putCmd, delCmd, watchCmd, backupCmd, statsCmd,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
