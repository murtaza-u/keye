package main

import (
	"context"
	"fmt"

	"github.com/murtaza-u/keye/client"
	"github.com/murtaza-u/keye/watch"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

var watchCmd = &cli.Command{
	Name:      "watch",
	Usage:     "Performs WATCH operation on the database",
	UsageText: "key [--regex]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "regex",
			Usage: "treat the key as a regex",
			Value: false,
		},
	},
	Action: func(ctx *cli.Context) error {
		key := ctx.Args().First()
		if key == "" {
			return fmt.Errorf("missing key")
		}

		c, err := client.New(client.Config{Addr: address})
		if err != nil {
			return err
		}
		defer c.Close()

		var opts []client.OptFunc
		if ctx.Bool("regex") {
			opts = append(opts, client.WithRegex())
		}
		listen(c, key, opts...)

		return nil
	},
}

func listen(c *client.C, k string, opts ...client.OptFunc) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w, err := c.Watch(ctx, k, opts...)
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			return fmt.Errorf(stat.Message())
		}
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-w.Error():
			return err
		case ev := <-w.NextEvent():
			var typ string
			switch ev.Type {
			case watch.EventPut:
				typ = "PUT"
			case watch.EventDel:
				typ = "DEL"
			default:
				typ = "UNKNOWN"
			}

			k := ev.KV.K
			v := ev.KV.V
			fmt.Printf("[EVENT] typ=%s | k=%s | v=%s\n", typ, k, v)
		}
	}
}
