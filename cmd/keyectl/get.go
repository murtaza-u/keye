package main

import (
	"fmt"

	"github.com/murtaza-u/keye/client"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

var getCmd = &cli.Command{
	Name:      "get",
	Usage:     "Performs GET operation on the database",
	UsageText: "key [--regex|--keys-only]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "regex",
			Usage: "treat the key as a regex",
		},
		&cli.BoolFlag{
			Name:  "keys-only",
			Usage: "retrieve only keys, without their values",
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
		if ctx.Bool("keys-only") {
			opts = append(opts, client.WithKeysOnly())
		}
		kvs, err := c.Get(key, opts...)
		if err != nil {
			if stat, ok := status.FromError(err); ok {
				return fmt.Errorf(stat.Message())
			}
			return err
		}

		for _, kv := range kvs {
			fmt.Printf("key=%s | v=%s\n", kv.GetKey(), string(kv.GetVal()))
		}

		return nil
	},
}
