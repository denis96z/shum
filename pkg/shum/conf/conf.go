package conf

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/maildealru/shum/pkg/shum/consts"
	"github.com/maildealru/shum/pkg/shum/errs"

	"github.com/akamensky/argparse"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Auth   AuthConfig   `yaml:"auth"`
	Shell  ShellConfig  `yaml:"shell"`
}

type AuthConfig struct {
	Clients []ClientConfig `yaml:"clients"`

	clients map[string]ClientConfig `yaml:"-"`
}

type ClientConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`

	TLS      ServerTLSConfig      `yaml:"tls"`
	Shutdown ServerShutdownConfig `yaml:"shutdown"`
}

type ServerTLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	KeyPath  string `yaml:"key_path"`
	CertPath string `yaml:"cert_path"`
}

type ServerShutdownConfig struct {
	Timeout time.Duration `yaml:"timeout"`
}

type ShellConfig struct {
	Bin  string   `yaml:"bin"`
	Args []string `yaml:"args"`

	Commands map[string]CommandConfig `yaml:"commands"`
}

type CommandConfig struct {
	Command      string `yaml:"command"`
	Async        bool   `yaml:"async"`
	RevealOutput bool   `yaml:"reveal_output"`
	Clients      string `yaml:"clients"`
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Addr: "127.0.0.1",
			Port: 8065,
			Shutdown: ServerShutdownConfig{
				Timeout: time.Second,
			},
		},
		Shell: ShellConfig{
			Bin:      "sh",
			Args:     []string{"-c"},
			Commands: make(map[string]CommandConfig),
		},
	}
}

func (c *Config) TryLoad(args []string) error {
	if isManagementCommand(args) {
		return nil
	}

	parser := argparse.NewParser(consts.Name, consts.Description)

	configFile := parser.String(
		"c", "config-file",
		&argparse.Options{
			Help:    "Path to config file",
			Default: "/etc/shum.yml",
		},
	)
	allowNoConf := parser.Flag(
		"", "allow-no-config-file",
		&argparse.Options{
			Help:    "If true, missing config file will not cause fail",
			Default: false,
		},
	)

	if err := parser.Parse(args); err != nil {
		return errs.Wrap(err, "failed to parse arguments")
	}

	if _, err := os.Stat(*configFile); err != nil {
		if *allowNoConf {
			return nil
		}
		return errs.Wrap(err, "config file is missing")
	}

	b, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return errs.Wrap(err, "failed to read config file")
	}
	if err = yaml.Unmarshal(b, c); err != nil {
		return errs.Wrap(err, "failed to parse config file")
	}

	if len(c.Auth.Clients) > 0 {
		c.Auth.clients = make(map[string]ClientConfig, len(c.Auth.Clients))
		for _, client := range c.Auth.Clients {
			c.Auth.clients[client.ClientID] = client
		}
	}

	return nil
}

func (c *AuthConfig) IsAuthOK(clientID, clientSecret string) bool {
	client, ok := c.clients[clientID]
	if !ok {
		return false
	}
	return client.ClientSecret == clientSecret
}

func isManagementCommand(args []string) bool {
	if len(args) != 2 {
		return false
	}

	v := args[1]
	return false ||
		v == "install" || v == "remove" ||
		v == "start" || v == "stop" || v == "status"
}
