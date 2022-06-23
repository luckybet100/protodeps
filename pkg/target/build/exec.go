package build

import (
	"fmt"
	"github.com/luckybet100/protodeps/pkg/analyze"
	"github.com/luckybet100/protodeps/pkg/deps"
	"github.com/luckybet100/protodeps/pkg/errors"
	"github.com/luckybet100/protodeps/pkg/scheme"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildTarget(
	cmd *cobra.Command,
	config *scheme.ProtoDepsConfig,
	target *scheme.Target,
) error {
	command := []string{"protoc"}

	for _, importPath := range config.Imports {
		command = append(command, "-I", importPath.Absolute())
	}

	for _, plugin := range target.Plugins {
		command = append(command, fmt.Sprintf("--%s_out=%s", plugin.Name, plugin.Output))
		if plugin.Opts != "" {
			command = append(command, fmt.Sprintf("--%s_opt=%s", plugin.Name, plugin.Opts))
		}
	}

	for _, source := range config.Sources {
		files, err := filepath.Glob(source.Absolute())
		if err != nil {
			return errors.RuntimeError.Wrapf(err, "Failed to list files by patten: %s", source)
		}
		log.Infof("Found %d matches by path %s\n", len(files), source.Absolute())
		command = append(command, files...)

	}

	log.Println(strings.Join(command, " "))
	buildCmd := exec.Command("/bin/bash", "-c", strings.Join(command, " "))
	buildCmd.Stdout = cmd.OutOrStdout()
	buildCmd.Stderr = cmd.ErrOrStderr()
	if err := buildCmd.Run(); err != nil {
		return errors.RuntimeError.Wrapf(err, "failed to build target %s", target.Name)
	}
	log.Println("Build finished successfully!")
	return nil
}

func Exec(cmd *cobra.Command, configFile string, targetName string) error {
	log.SetLevel(log.ErrorLevel)
	if config, err := analyze.Exec(configFile); err != nil {
		return err
	} else {
		log.SetLevel(log.InfoLevel)
		for _, dep := range config.Deps {
			if err := deps.ResolveDep(cmd, dep); err != nil {
				return err
			}
		}
		for _, dir := range config.CreateDirs {
			if err := os.MkdirAll(dir, 0777); err != nil {
				return errors.RuntimeError.Wrapf(err, "Failed to create required directory")
			}
		}
		for _, target := range config.Targets {
			for _, plugin := range target.Plugins {
				if ok, err := plugin.IsInstalled(); err != nil {
					return err
				} else if !ok {
					if ok, err := plugin.Install(cmd); err != nil {
						return err
					} else if ok {
						log.Infof("Plugin <%s> successfully installed", plugin.Name)
					} else {
						return errors.RuntimeError.Newf("error, missing plugin %s", plugin.Name)
					}
				}
			}
			if target.Name == targetName {
				return buildTarget(cmd, config, target)
			}
		}
		return errors.InvalidArgument.Newf(" invalid target: <%s>.\nUse `protodeps target list` to list available targets", targetName)
	}
}
