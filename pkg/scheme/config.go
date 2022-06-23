package scheme

import (
	"github.com/luckybet100/protodeps/pkg/constants"
	"github.com/luckybet100/protodeps/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type ProtoDepsConfig struct {
	Project    string          `yaml:"project" json:"project"`
	Version    string          `yaml:"version" json:"version"`
	Targets    []*Target       `yaml:"targets" json:"targets"`
	Sources    []*Path         `yaml:"src" json:"src"`
	Imports    []*Path         `yaml:"imports" json:"imports"`
	Deps       []*Dependency   `yaml:"deps" json:"deps"`
	Plugins    []*CustomPlugin `yaml:"plugins" json:"plugins"`
	CreateDirs []string        `yaml:"create_dirs" json:"create_dirs"`
}

func (config *ProtoDepsConfig) validateVersion() error {
	if config.Version == "" {
		return errors.ValidationError.Newf("config version should be specified")
	}
	parts := strings.Split(config.Version, ".")
	if len(parts) != 2 {
		return errors.ValidationError.Newf("invalid version")
	}
	if value, err := strconv.Atoi(parts[0]); err != nil {
		return errors.ValidationError.Newf("invalid version")
	} else if value != constants.VersionMajor {
		return errors.InvalidArgument.Newf("unsupported config version %s, supported: %d.0 <= version < %d.%d",
			config.Version,
			constants.VersionMajor,
			constants.VersionMajor,
			constants.VersionMinor+1,
		)
	}
	if value, err := strconv.Atoi(parts[1]); err != nil {
		return errors.ValidationError.Newf("invalid version")
	} else if value > constants.VersionMinor {
		return errors.InvalidArgument.Newf("unsupported config version %s, supported: %d.0 <= version < %d.%d",
			config.Version,
			constants.VersionMajor,
			constants.VersionMajor,
			constants.VersionMinor+1,
		)
	}
	return nil
}

func (config *ProtoDepsConfig) Validate() error {
	if config.Project == "" {
		return errors.ValidationError.Newf("project name should be specified")
	}
	if err := config.validateVersion(); err != nil {
		return err
	}
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
