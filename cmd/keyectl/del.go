package main

import (
	"fmt"

	"github.com/murtaza-u/keye/client"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

var delCmd = &cli.Command{
	Name:      "del",
	Usage:     "Performs DEL operation on the database",
	UsageText: "key [--regex]",
	Aliases:   []string{"delete"},
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

		parent := ctx.Lineage()[0]
		c, err := client.New(client.Config{
			Addr:    parent.String("address"),
			Timeout: parent.Duration("timeout"),
		})
		if err != nil {
			return err
		}
		defer c.Close()

		var opts []client.OptFunc
		if ctx.Bool("regex") {
			opts = append(opts, client.WithRegex())
		}
		keys, err := c.Del(key, opts...)
		if err != nil {
			if stat, ok := status.FromError(err); ok {
				return fmt.Errorf(stat.Message())
			}
			return err
		}

		fmt.Println("deleted keys:")
		for _, k := range keys {
			fmt.Println(k)
		}

		return nil
	},
}
