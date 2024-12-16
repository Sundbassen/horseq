package cli

import (
	"context"
	"io"

	"github.com/peterbourgon/ff/v4"
)

type RootCmdOptions struct {
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

type RootCmd struct {
	Opts    RootCmdOptions
	Flags   *ff.FlagSet
	Command *ff.Command
}

const CmdLabel = "horseq"

func NewRootCmd(stdin io.Reader, stdout, stderr io.Writer) *RootCmd {
	cmd := RootCmd{
		Opts: RootCmdOptions{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		},
	}
	cmd.Flags = ff.NewFlagSet(CmdLabel)

	cmd.Command = &ff.Command{
		Name:      CmdLabel,
		ShortHelp: CmdLabel + " data pipeline for mapping transactions",
		Usage:     CmdLabel + " [flags] <subcommand> ...",
		Flags:     cmd.Flags,
		Exec: func(ctx context.Context, args []string) error {
			return ff.ErrHelp
		},
	}

	return &cmd
}
