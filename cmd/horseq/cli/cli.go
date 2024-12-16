package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"github.com/sundbassen/horseq/internal/util"
)

func Exec(
	ctx context.Context,
	version string,
	args []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (err error) {
	root := NewRootCmd(stdin, stdout, stderr)
	_ = NewDPCmd(root)
	defer func() {
		if errors.Is(err, util.ErrCliRequiredFlags) {
			fmt.Fprintf(stderr, "\nerror: %s\n", err)
		}

		if errors.Is(err, ff.ErrHelp) {
			fmt.Fprintf(stderr, "\n%s\n", ffhelp.Command(root.Command))
		}
	}()

	if err = root.Command.Parse(
		args,
		ff.WithEnvVars(),
	); err != nil {
		return err
	}

	// defer func() {
	// 	_ = tracingShutdown(ctx)
	// }()

	if err = root.Command.Run(ctx); err != nil {
		slog.Error("failed to run command", slog.Any("error", err))
		return err
	}

	return nil
}
