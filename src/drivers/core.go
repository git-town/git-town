package drivers

import (
	"fmt"
	"log"

	"github.com/Originate/git-town/src/git"
	"github.com/fatih/color"
)

// Core provides the public API for the drivers subsystem

var registry = Registry{}

var activeDriver CodeHostingDriver

// GetActiveDriver returns the code hosting driver to use based on the git config
func GetActiveDriver() CodeHostingDriver {
	if activeDriver == nil {
		var err error
		activeDriver, err = GetDriver(git.GetRemoteOriginURL())
		if err != nil {
			log.Fatal(err)
		}
	}
	return activeDriver
}

// GetDriver returns the code hosting driver to use based on the git config
func GetDriver(originURL string) (CodeHostingDriver, error) {
	return registry.DetermineActiveDriver(originURL)
}

func printLog(message string) {
	fmt.Println()
	color.New(color.Bold).Println(message)
}
