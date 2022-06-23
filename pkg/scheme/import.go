package scheme

import (
	"github.com/luckybet100/protodeps/pkg/constants"
	"github.com/luckybet100/protodeps/pkg/errors"
)

type Path struct {
	Path string `yaml:"path" json:"path"`
	From string `yaml:"from" json:"from"`
}

func (path *Path) Validate() error {
	if path.Path == "" {
		return errors.ValidationError.New("import path should not be empty")
	}
	return nil
}

func (path *Path) Absolute() string {
	result := path.Path
	if path.From != "" {
		if result == "." {
			result = constants.DepsFolder + "/" + path.From
		} else {
			result = constants.DepsFolder + "/" + path.From + "/" + result
		}
	}
	return result
}
