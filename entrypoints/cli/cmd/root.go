package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zmoog/eh/entrypoints/cli/cmd/tail"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "eh",
		Short: "Azure Event Hub CLI",
	}

	rootCmd.AddCommand(tail.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
