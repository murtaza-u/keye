package main

import (
	"fmt"
	"log"
	"os"

	"github.com/murtaza-u/keye/client"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

var backupCmd = &cli.Command{
	Name:      "backup",
	Usage:     "Takes a backup of the database",
	UsageText: "FILE",
	Aliases:   []string{"snapshot"},
	Action: func(ctx *cli.Context) error {
		out := ctx.Args().First()
		if out == "" {
			return fmt.Errorf("missing output file")
		}

		c, err := client.New(client.Config{
			Addr:    address,
			Timeout: timeout,
		})
		if err != nil {
			return err
		}

		snap, err := c.Backup(128)
		if err != nil {
			if stat, ok := status.FromError(err); ok {
				return fmt.Errorf(stat.Message())
			}
			return err
		}

		err = os.WriteFile(out, snap, 0600)
		if err != nil {
			return err
		}
		log.Printf("written backup to %q", out)

		return nil
	},
}
