package config

import (
	"errors"
	"fmt"
	"strconv"
)

type Offline struct {
	gitConfig *gitConfig
}

// Enabled indicates whether Git Town is currently in offline mode.
func (o *Offline) Enabled() bool {
	config := o.gitConfig.globalConfigValue("git-town.offline")
	if config == "" {
		return false
	}
	result, err := strconv.ParseBool(config)
	if err != nil {
		fmt.Printf("Invalid value for git-town.offline: %q. Please provide either true or false. Considering false for now.", config)
		fmt.Println()
		return false
	}
	return result
}

// Enable updates whether Git Town is in offline mode.
func (o *Offline) Enable(value bool) error {
	_, err := o.gitConfig.SetGlobalConfigValue("git-town.offline", strconv.FormatBool(value))
	return err
}

// Validate asserts that Git Town is not in offline mode.
func (o *Offline) Validate() error {
	if o.Enabled() {
		return errors.New("this command requires an active internet connection")
	}
	return nil
}
