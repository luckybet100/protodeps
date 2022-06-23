package list

import (
	"github.com/luckybet100/protodeps/pkg/analyze"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Exec(cmd *cobra.Command, configFile string) error {
	log.SetLevel(log.ErrorLevel)
	if config, err := analyze.Exec(configFile); err != nil {
		return err
	} else {
		log.SetLevel(log.InfoLevel)
		for _, target := range config.Targets {
			log.Println(target.Name)
		}
		return nil
	}
}
