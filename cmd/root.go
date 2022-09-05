package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var (
	appName = "reth"
)

const (
	flagKeystorePath = "keystore_path"
	flagLogLevel     = "log_level"

	defaultKeystorePath = "./keys"
)

// NewRootCmd returns the root command.
func NewRootCmd() *cobra.Command {
	// RootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "reth service",
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {

		return nil
	}

	rootCmd.AddCommand(
		genAccountCmd(),
		versionCmd(),
	)
	return rootCmd
}

func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	ctx := context.Background()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
