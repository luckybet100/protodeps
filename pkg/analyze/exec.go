package analyze

import (
	"github.com/luckybet100/protodeps/pkg/errors"
	"github.com/luckybet100/protodeps/pkg/scheme"
	"os"
)

func Exec(configFile string) (*scheme.ProtoDepsConfig, error) {
	if file, err := os.Open(configFile); err != nil {
		return nil, errors.ValidationError.Wrap(err, "failed to open file")
	} else if parser, err := scheme.GetParser(file); err != nil {
		return nil, err
	} else if config, err := parser.Parse(); err != nil {
		return nil, err
	} else {
		for _, plugin := range config.Plugins {
			if err := plugin.MakeKnown(); err != nil {
				return nil, err
			}
		}
		return config, config.Validate()
	}
}
