package cli

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"github.com/peterbourgon/ff/v4"
	"github.com/sundbassen/horseq/component/service"
	"github.com/sundbassen/horseq/component/transaction/transactionstore"
)

type DPCmdOptions struct {
	// GCP bucket name
	BucketName string

	// GCP project ID
	ProjectID string

	// GCP region
	Region string

	// CSV file path
	CsvPath string

	Root *RootCmdOptions
}

type DPCmd struct {
	Opts DPCmdOptions

	root    *RootCmd
	Flags   *ff.FlagSet
	Command *ff.Command
}

func NewDPCmd(parent *RootCmd) *DPCmd {
	var cmd DPCmd
	cmd.Opts.Root = &parent.Opts
	cmd.root = parent
	cmd.Flags = ff.NewFlagSet("datapipeline").SetParent(parent.Flags)
	cmd.Flags.StringVar(&cmd.Opts.BucketName, 'n', "name", "gs://example-bucket-1-bananas/", "Name of bucket to read from")
	cmd.Flags.StringVar(&cmd.Opts.ProjectID, 'p', "project-id", "", "GCP project ID")
	cmd.Flags.StringVar(&cmd.Opts.Region, 'r', "region", "europe-west1", "GCP region")
	cmd.Flags.StringVar(&cmd.Opts.CsvPath, 'c', "csv-path", "", "Path to CSV file in bucket")

	cmd.Command = &ff.Command{
		Name:  "datapipeline",
		Usage: CmdLabel + " datapipeline [flags]",
		Flags: cmd.Flags,
		Exec:  DPCmdExec(&cmd.Opts),
	}
	cmd.root.Command.Subcommands = append(cmd.root.Command.Subcommands, cmd.Command)

	return &cmd
}

func DPCmdExec(opts *DPCmdOptions) func(context.Context, []string) error {
	return func(ctx context.Context, args []string) error {
		client, err := storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create GCS client: %v", err)
		}

		defer client.Close()

		bucket := client.Bucket(opts.BucketName)
		transactionRead := transactionstore.NewBucket(bucket, opts.CsvPath)

		bqClient, err := bigquery.NewClient(ctx, opts.ProjectID)
		if err != nil {
			return fmt.Errorf("failed to create BigQuery client: %v", err)
		}

		transactionWrite := transactionstore.NewBQWriter(bqClient)

		transactionSVC := service.NewTransactionService(transactionRead, transactionWrite)

		return transactionSVC.MapToNew(ctx)
	}
}
