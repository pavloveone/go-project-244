package main

import (
	"code"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "format",
		Aliases:     []string{"f"},
		Usage:       "output format",
		DefaultText: `"stylish"`,
	},
}

func main() {
	command := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: flags,
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("file paths are required")
			}
			paths := c.Args().Slice()
			out, err := code.ParseByPaths(paths)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		},
	}
	if err := command.Run(context.Background(), os.Args); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
