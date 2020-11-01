package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

//GeneralConfig - contains general configuration values for clooudwatch search
type GeneralConfig struct {
	RunInterval string `toml:"run_interval"`
	Template    string `toml:"sns_message_template"`
}

//Rule - contains rules for clooudwatch search
type Rule struct {
	Name      string `toml:"name"`
	LogGroup  string `toml:"log_group"`
	NumEvents int64  `toml:"num_events"`
	Timeframe string `toml:"timeframe"`
	Query     string `toml:"query"`
	SnsTopic  string `toml:"sns_topic"`
	NotifyOk  bool   `toml:"notify_ok,omitempty"`
}

//Config - config structure
type Config struct {
	General GeneralConfig   `toml:"general"`
	Rules   map[string]Rule `toml:"rules"`
}

//ParseConfig parses toml config and returns it back as struct
func ParseConfig(path string) (*Config, error) {
	filename, _ := filepath.Abs(path)
	tomlfile, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = toml.Unmarshal(tomlfile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
