package cmd

import (
	"github.com/spf13/cobra"
)

var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "subcommand to work with targets",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Printf("Please specify command\n")
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(targetCmd)
}
