package wsbridge

import (
	"net/url"
	"time"
)

// Config is a wsbridge global config
type Config struct {
	Host             string
	Port             int
	Server           url.URL
	Proxy            *url.URL
	HandshakeTimeout time.Duration
}

// NewConfig creates a new config instance with default values
func NewConfig(host string, port int, server url.URL) Config {
	return Config{
		Host:             host,
		Port:             port,
		Server:           server,
		Proxy:            nil,
		HandshakeTimeout: time.Duration(1 * time.Minute),
	}
}
