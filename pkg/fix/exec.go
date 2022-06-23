package fix

import (
	"github.com/luckybet100/protodeps/pkg/analyze"
	"github.com/luckybet100/protodeps/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/exec"
)

func checkProtoc() bool {
	_, err := exec.LookPath("protoc")
	return err == nil
}

func installProtoc(cmd *cobra.Command) error {
	return errors.RuntimeError.New("Please install protoc manually: https://grpc.io/docs/protoc-installation/")
}

func Exec(cmd *cobra.Command, configFile string) error {
	if config, err := analyze.Exec(configFile); err != nil {
		return err
	} else {
		log.Infof("Check protoc")
		if !checkProtoc() {
			if err := installProtoc(cmd); err != nil {
				return err
			}
		}
		log.Infof("OK")
		log.Infof("Start installing missing plugins")
		for _, target := range config.Targets {
			for _, plugin := range target.Plugins {
				if ok, err := plugin.IsInstalled(); err != nil {
					return err
				} else if !ok {
					if ok, err := plugin.Install(cmd); err != nil {
						return err
					} else if ok {
						log.Infof("Plugin <%s> successfully installed", plugin.Name)
					}
				}
			}
		}
		log.Infof("Installation of missing files finished")
		return nil
	}
}
