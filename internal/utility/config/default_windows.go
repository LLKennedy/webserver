package config

import (
	"fmt"
	"os"
)

var (
	defaultKeyFile  = fmt.Sprintf("%s/pki/server.key", os.Getenv("PROGRAMDATA"))
	defaultCertFile = fmt.Sprintf("%s/pki/server.crt", os.Getenv("PROGRAMDATA"))
)
