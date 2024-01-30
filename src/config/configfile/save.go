package configfile

import (
	"os"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

func Save(config *configdomain.FullConfig) error {
	toml := RenderTOML(config)
	return os.WriteFile(FileName, []byte(toml), 0o600)
}

func RenderTOML(config *configdomain.FullConfig) string {
	return ""
}
