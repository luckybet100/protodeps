package scheme

import (
	"github.com/luckybet100/protodeps/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ProtoDepsConfig struct {
	Targets    []*Target       `yaml:"targets" json:"targets"`
	Sources    []*Path         `yaml:"src" json:"src"`
	Imports    []*Path         `yaml:"imports" json:"imports"`
	Deps       []*Dependency   `yaml:"deps" json:"deps"`
	Plugins    []*CustomPlugin `yaml:"plugins" json:"plugins"`
	CreateDirs []string        `yaml:"create_dirs" json:"create_dirs"`
}

func (config *ProtoDepsConfig) Validate() error {
	log.Infoln("Start validating targets...")
	for _, target := range config.Targets {
		if err := target.Validate(); err != nil {
			return err
		}
	}
	log.Infoln("Targets validation finished!")
	deps := map[string]bool{"": true}
	for _, dep := range config.Deps {
		deps[dep.Name] = true
	}
	log.Infoln("Start validating imports...")
	for _, importPath := range config.Imports {
		if err := importPath.Validate(); err != nil {
			return err
		}
		if !deps[importPath.From] {
			return errors.ValidationError.Newf("Unknown import from: %s", importPath.From)
		}
	}
	log.Infoln("Imports validation finished!")
	log.Infoln("Start validating sources...")
	for _, path := range config.Sources {
		if err := path.Validate(); err != nil {
			return err
		}
	}
	log.Infoln("Imports validation sources!")
	return nil
}
