package cmd

import (
	"github.com/luckybet100/protodeps/pkg/target/build"
	"github.com/luckybet100/protodeps/pkg/target/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds target",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
		log.SetOutput(cmd.OutOrStdout())
		if len(args) == 0 {
			log.Errorln("Expected positional argument <target>")
			log.Println("List of available targets:")
			if err := list.Exec(cmd, fileFlag); err != nil {
				log.Errorln(err)
			}
			return
		}
		target := args[0]
		if err := build.Exec(cmd, fileFlag, target); err != nil {
			log.Errorln(err)
		}
	},
}

func init() {
	targetCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringVarP(
		&fileFlag,
		build.FileFlag,
		build.FileFlagShort,
		build.FileFlagDefault,
		build.FileFlagUsage,
	)
}
