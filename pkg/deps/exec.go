package deps

import (
	"fmt"
	"github.com/luckybet100/protodeps/pkg/analyze"
	"github.com/luckybet100/protodeps/pkg/constants"
	"github.com/luckybet100/protodeps/pkg/errors"
	"github.com/luckybet100/protodeps/pkg/scheme"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

func cloneDep(cmd *cobra.Command, name string, repo string, location string) error {
	command := []string{"git", "clone", repo, location}
	cloneCmd := exec.Command("/bin/bash", "-c", strings.Join(command, " "))
	cloneCmd.Stdout = cmd.OutOrStdout()
	cloneCmd.Stderr = cmd.ErrOrStderr()
	if err := cloneCmd.Run(); err != nil {
		return errors.RuntimeError.Wrapf(err, "failed to clone dependency: %s", name)
	} else {
		log.Infof("Clone %s finished successfully!\n", name)
		return nil
	}
}

func checkoutDep(cmd *cobra.Command, name string, location string, ref string) error {
	command := []string{"git", "checkout", ref}
	cloneCmd := exec.Command("/bin/bash", "-c", strings.Join(command, " "))
	cloneCmd.Stdout = cmd.OutOrStdout()
	cloneCmd.Stderr = cmd.ErrOrStderr()
	cloneCmd.Dir = location
	if err := cloneCmd.Run(); err != nil {
		return errors.RuntimeError.Wrapf(err, "failed to checkout dependency: %s %s", name, ref)
	} else {
		log.Infof("Checkout %s finished successfully!\n", name)
		return nil
	}
}

func ResolveDep(cmd *cobra.Command, dep *scheme.Dependency) error {
	location := fmt.Sprintf("%s/%s", constants.DepsFolder, dep.Name)
	if _, err := os.Stat(location); err != nil && !os.IsNotExist(err) {
		return errors.InternalError.Wrapf(err, "failed to deps file info <%s>", location)
	} else if err != nil {
		if err := cloneDep(cmd, dep.Name, dep.GitRepo, location); err != nil {
			return err
		}
	}
	if dep.Ref != "" {
		return checkoutDep(cmd, dep.Name, location, dep.Ref)
	}
	return nil
}

func Exec(cmd *cobra.Command, configFile string, upgrade bool) error {
	log.SetLevel(log.ErrorLevel)
	if config, err := analyze.Exec(configFile); err != nil {
		return err
	} else {
		log.SetLevel(log.InfoLevel)
		if upgrade {
			if err := os.RemoveAll(constants.DepsFolder); err != nil {
				return errors.RuntimeError.Newf("Failed to clean %s directory", constants.DepsFolder)
			}
		}
		for _, dep := range config.Deps {
			if err := ResolveDep(cmd, dep); err != nil {
				return err
			}
		}
		return nil
	}
}
