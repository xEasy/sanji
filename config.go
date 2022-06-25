package sanji

import (
	"net"
	"net/http"
	"time"

	"github.com/bitleak/lmstfy/client"
)

const (
	DefaultNamespace    = "default"
	DefaultPort         = 7777
	DefaultPollInterval = 5
	maxReadTimeout      = 600
)

type Config struct {
	client       *client.LmstfyClient
	Namespace    string
	Token        string
	Host         string
	Port         int
	PollInterval int
}

func Configure(c *Config) *Config {
	c.assignDefaultIfBlank()
	c.validate()
	c.initClient()
	return c
}

func (c *Config) initClient() {
	cli := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        128,
			MaxIdleConnsPerHost: 32,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: time.Minute,
			}).DialContext,
		},
		Timeout: maxReadTimeout * time.Second,
	}

	c.client = client.NewLmstfyWithClient(cli, c.Host, c.Port, c.Namespace, c.Token)
}

func (c *Config) validate() {
	if c.Token == "" {
		panic("missing LmstfyClient token config!")
	}
	if c.Host == "" {
		panic("missing LmstfyClient host config!")
	}
}

func (c *Config) assignDefaultIfBlank() {
	if c.Namespace == "" {
		c.Namespace = "default"
	}
	if c.Port == 0 {
		c.Port = DefaultPort
	}
	if c.PollInterval == 0 {
		c.PollInterval = DefaultPollInterval
	}
}
