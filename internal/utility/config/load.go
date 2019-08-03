package config

import "golang.org/x/tools/godoc/vfs"

// Load loads the default config location on the provided file system, returning defaults for any keys not present in the config file
func Load(fs vfs.FileSystem, flags map[string]interface{}) (Options, error) {
	opts := DefaultOptions()
	return opts, nil
}
