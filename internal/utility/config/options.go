package config

// Options are the config options for the application
type Options struct {
	Address       string `json:"address"`
	Port          uint16 `json:"port"`
	InsecurePort  uint16 `json:"insecurePort"`
	KeyFile       string `json:"keyFile"`
	CertFile      string `json:"certFile"`
	StaticContent string `json:"staticContent"`
}

// DefaultOptions are the default options
func DefaultOptions() Options {
	return Options{
		Address:       "localhost",
		Port:          443,
		InsecurePort:  80,
		KeyFile:       defaultKeyFile,
		CertFile:      defaultCertFile,
		StaticContent: "build",
	}
}

// GetAddress returns the address
func (o Options) GetAddress() string {
	if o.Address == "" {
		return DefaultOptions().GetAddress()
	}
	return o.Address
}

// GetPort returns the port
func (o Options) GetPort() uint16 {
	if o.Port == 0 {
		return DefaultOptions().GetPort()
	}
	return o.Port
}

// GetInsecurePort returns the insecure port
func (o Options) GetInsecurePort() uint16 {
	if o.InsecurePort == 0 {
		return DefaultOptions().GetInsecurePort()
	}
	return o.InsecurePort
}

// GetKeyFile returns the key file location
func (o Options) GetKeyFile() string {
	if o.KeyFile == "" {
		return DefaultOptions().GetKeyFile()
	}
	return o.KeyFile
}

// GetCertFile returns the cert file location
func (o Options) GetCertFile() string {
	if o.CertFile == "" {
		return DefaultOptions().GetCertFile()
	}
	return o.CertFile
}

// GetStaticContent returns the static content location
func (o Options) GetStaticContent() string {
	if o.StaticContent == "" {
		return DefaultOptions().GetStaticContent()
	}
	return o.StaticContent
}
