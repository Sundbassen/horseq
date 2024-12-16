package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v4"
	"github.com/sundbassen/horseq/cmd/horseq/cli"
)

const version = "v0.0.0"

func main() {
	var (
		ctx    = context.Background()
		args   = os.Args[1:]
		stdin  = os.Stdin
		stdout = os.Stdout
		stderr = os.Stderr
	)

	err := cli.Exec(ctx, version, args, stdin, stdout, stderr)

	switch {
	case errors.Is(err, ff.ErrNoExec):
	case errors.Is(err, ff.ErrHelp):
	case err != nil:
		fmt.Printf("horseq exited with error:\n% +-.3v\n", err)
		os.Exit(1)
	}
}
