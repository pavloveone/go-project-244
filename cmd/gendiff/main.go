package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	command := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Action: func(ctx context.Context, c *cli.Command) error {
			return nil
		},
	}
	if err := command.Run(context.Background(), os.Args); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
