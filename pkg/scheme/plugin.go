package scheme

import (
	"fmt"
	"github.com/luckybet100/protodeps/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

type CustomPlugin struct {
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
	Git     string `yaml:"git" json:"git"`
}

func (plugin *CustomPlugin) MakeKnown() error {
	if plugin.Name == "" {
		return errors.ValidationError.New("CustomPlugin <name> should not be empty")
	}
	if plugin.Version == "" {
		return errors.ValidationError.New("CustomPlugin <version> should not be empty")
	}
	if plugin.Git == "" {
		return errors.ValidationError.New("CustomPlugin <version> should not be empty")
	}
	knownPlugins[plugin.Name] = &knownPlugin{
		name:       plugin.Name,
		installCmd: fmt.Sprintf("go install %s@%s", plugin.Git, plugin.Version),
		checkInstalled: func() (bool, error) {
			return checkFilename(os.Getenv("GOBIN") + "/" + plugin.Name)
		},
	}
	return nil
}

type ProtobufPlugin struct {
	Name      string   `yaml:"name" json:"name"`
	Output    string   `yaml:"out" json:"out"`
	Opts      string   `yaml:"opts" json:"opts"`
	Arguments []string `yaml:"args" json:"args"`
}

func (plugin *ProtobufPlugin) Validate() error {
	log.Infof("Validating plugin <%s>...\n", plugin.Name)
	if plugin.Name == "" {
		return errors.ValidationError.New("Plugin <name> should not be empty")
	} else if ok, err := plugin.IsInstalled(); err != nil {
		return err
	} else if !ok {
		log.Warnf("Warning, plugin <%s> is not installed...\n", plugin.Name)
	}
	if plugin.Output == "" {
		return errors.ValidationError.New("Plugin <out> directory should not be empty")
	}
	log.Infoln("OK")
	return nil
}

func (plugin *ProtobufPlugin) IsInstalled() (bool, error) {
	if knownPlugin, ok := knownPlugins[plugin.Name]; !ok {
		return false, errors.ValidationError.New("Unknown plugin")
	} else {
		if ok, err := knownPlugin.checkInstalled(); err != nil {
			return false, err
		} else if !ok {
			return false, nil
		} else {
			return true, nil
		}
	}
}

func (plugin *ProtobufPlugin) Install(cmd *cobra.Command) (bool, error) {
	if ok, err := plugin.IsInstalled(); ok || err != nil {
		return ok, err
	}
	pluginInfo := knownPlugins[plugin.Name]
	installCmd := exec.Command("/bin/bash", "-c", pluginInfo.installCmd)
	installCmd.Stdout = cmd.OutOrStdout()
	installCmd.Stderr = cmd.ErrOrStderr()
	if err := installCmd.Run(); err != nil {
		return false, errors.RuntimeError.Wrapf(err, "failed to install plugin %s", pluginInfo.name)
	} else {
		return true, nil
	}
}

type knownPlugin struct {
	name           string
	installCmd     string
	checkInstalled func() (bool, error)
}

func checkFilename(filename string) (bool, error) {
	if _, err := os.Stat(filename); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, errors.InternalError.Wrapf(err, "failed to deps file info <%s>", filename)
	}
}

var knownPlugins = map[string]*knownPlugin{
	"go": {
		name:       "go",
		installCmd: "go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28",
		checkInstalled: func() (bool, error) {
			return checkFilename(os.Getenv("GOBIN") + "/protoc-gen-go")
		},
	},
	"go-grpc": {
		name:       "go-grpc",
		installCmd: "go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2",
		checkInstalled: func() (bool, error) {
			return checkFilename(os.Getenv("GOBIN") + "/protoc-gen-go")
		},
	},
	"dart": {
		name:       "dart",
		installCmd: "dart pub global activate protoc_plugin",
		checkInstalled: func() (bool, error) {
			return checkFilename(os.Getenv("HOME") + "/.pub-cache/bin/protoc-gen-dart")
		},
	},
}
