package config

// Options are the config options for the application
type Options struct {
	Address      string `json:"address"`
	Port         uint16 `json:"port"`
	InsecurePort uint16 `json:"insecurePort"`
	KeyFile      string `json:"keyFile"`
	CertFile     string `json:"certFile"`
}

// DefaultOptions are the default options
func DefaultOptions() Options {
	return Options{
		Address:      "localhost",
		Port:         443,
		InsecurePort: 80,
		KeyFile:      defaultKeyFile,
		CertFile:     defaultCertFile,
	}
}
