package cmd

import (
	"github.com/luckybet100/protodeps/pkg/deps"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Downloads all dependecies",
	Run: func(cmd *cobra.Command, _ []string) {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
		log.SetOutput(cmd.OutOrStdout())
		if err := deps.Exec(cmd, fileFlag, false); err != nil {
			log.Errorln(err)
		} else {
			log.Infoln("Dependencies successfully resolved!")
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(
		&fileFlag,
		deps.FileFlag,
		deps.FileFlagShort,
		deps.FileFlagDefault,
		deps.FileFlagUsage,
	)
}
