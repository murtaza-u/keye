package main

import (
	"fmt"

	"github.com/murtaza-u/keye/client"
	"google.golang.org/grpc/status"

	"github.com/urfave/cli/v2"
)

var putCmd = &cli.Command{
	Name:      "put",
	Usage:     "Performs PUT operation on the database",
	UsageText: "key val [--regex]",
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
		val := ctx.Args().Get(1)
		if val == "" {
			return fmt.Errorf("missing val")
		}

		c, err := client.New(client.Config{
			Addr:    address,
			Timeout: timeout,
		})
		if err != nil {
			return err
		}
		defer c.Close()

		var opts []client.OptFunc
		if ctx.Bool("regex") {
			opts = append(opts, client.WithRegex())
		}
		keys, err := c.Put(key, []byte(val), opts...)
		if err != nil {
			if stat, ok := status.FromError(err); ok {
				return fmt.Errorf(stat.Message())
			}
			return err
		}

		fmt.Println("modified keys:")
		for _, k := range keys {
			fmt.Println(k)
		}

		return nil
	},
}
