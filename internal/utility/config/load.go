package config

import (
	"os"

	"github.com/LLKennedy/goconfig"
	"golang.org/x/tools/godoc/vfs"
)

var envGetter = os.Getenv

// Load loads the default config location on the provided file system, returning defaults for any keys not present in the config file.
// Options start as defaults, then load from environment variables, then a config file, then runtime flags, each overriding the previous.
func Load(fs vfs.FileSystem, flags []string) (Options, error) {
	opts := DefaultOptions()
	err := goconfig.Load(&opts, "webserver", goconfig.ParseArgs(flags), fs, nil)
	return opts, err
}
