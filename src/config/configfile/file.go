package configfile

import (
	"log"

	"github.com/BurntSushi/toml"
)

func load() {
	_, err := toml.Decode(blob, &contacts)
	if err != nil {
		log.Fatal(err)
	}
}
