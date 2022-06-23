package scheme

import (
	"github.com/luckybet100/protodeps/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Target struct {
	Name    string           `yaml:"name" json:"name"`
	Plugins []ProtobufPlugin `yaml:"plugins" json:"plugins"`
}

func (target *Target) Validate() error {
	log.Infof("Validating target <%s>...\n", target.Name)
	if target.Name == "" {
		return errors.ValidationError.New("Plugin <name> should not be empty")
	}
	for _, plugin := range target.Plugins {
		if err := plugin.Validate(); err != nil {
			return err
		}
	}
	log.Infoln("OK")
	return nil
}
