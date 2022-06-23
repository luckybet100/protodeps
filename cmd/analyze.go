package cmd

import (
	"github.com/luckybet100/protodeps/pkg/analyze"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Validates project config file",
	Run: func(cmd *cobra.Command, _ []string) {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
		log.SetOutput(cmd.OutOrStdout())
		if _, err := analyze.Exec(fileFlag); err != nil {
			log.Errorln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(
		&fileFlag,
		analyze.FileFlag,
		analyze.FileFlagShort,
		analyze.FileFlagDefault,
		analyze.FileFlagUsage,
	)
}
