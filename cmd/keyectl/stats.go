package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/murtaza-u/keye/client"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

var statsCmd = &cli.Command{
	Name:      "stats",
	Usage:     "Obtain database statistics over a time range",
	UsageText: "[--duration]",
	Aliases:   []string{"statistics"},
	Flags: []cli.Flag{
		&cli.DurationFlag{
			Name:  "duration",
			Value: time.Second * 10,
		},
	},
	Action: func(ctx *cli.Context) error {
		dur := ctx.Duration("duration")
		c, err := client.New(client.Config{
			Addr:    address,
			Timeout: timeout,
		})
		if err != nil {
			return err
		}

		stats, err := c.Stats(dur)
		if err != nil {
			if stat, ok := status.FromError(err); ok {
				return fmt.Errorf(stat.Message())
			}
			return err
		}

		err = json.NewEncoder(os.Stdout).Encode(stats)
		if err != nil {
			return err
		}

		return nil
	},
}
