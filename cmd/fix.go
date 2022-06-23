package cmd

import (
	"github.com/luckybet100/protodeps/pkg/fix"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Installs missing known plugins",
	Run: func(cmd *cobra.Command, _ []string) {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
		log.SetOutput(cmd.OutOrStdout())
		if err := fix.Exec(cmd, fileFlag); err != nil {
			log.Errorln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
	fixCmd.Flags().StringVarP(
		&fileFlag,
		fix.FileFlag,
		fix.FileFlagShort,
		fix.FileFlagDefault,
		fix.FileFlagUsage,
	)
}
