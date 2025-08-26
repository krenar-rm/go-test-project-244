package main

import (
	"code"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage:   "output format (default: \"stylish\")",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Validate arguments
			if cmd.NArg() != 2 {
				return fmt.Errorf("exactly two file paths are required")
			}

			path1 := cmd.Args().Get(0)
			path2 := cmd.Args().Get(1)
			format := cmd.String("format")

			// Generate diff using the library function
			result, err := code.GenDiff(path1, path2, format)
			if err != nil {
				return fmt.Errorf("failed to generate diff: %w", err)
			}

			// Output the result
			fmt.Print(result)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
