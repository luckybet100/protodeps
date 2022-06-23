package cmd

import (
	"github.com/luckybet100/protodeps/pkg/target/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints list of all available targets",
	Run: func(cmd *cobra.Command, _ []string) {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
		log.SetOutput(cmd.OutOrStdout())
		if err := list.Exec(cmd, fileFlag); err != nil {
			log.Errorln(err)
		}
	},
}

func init() {
	targetCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(
		&fileFlag,
		list.FileFlag,
		list.FileFlagShort,
		list.FileFlagDefault,
		list.FileFlagUsage,
	)
}
