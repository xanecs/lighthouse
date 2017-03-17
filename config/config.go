package config

import (
	"io/ioutil"

	"github.com/naoina/toml"
)

//TomlConfig represents a parsed config file
type TomlConfig struct {
	Redis  RedisConfig
	Boards map[string]BoardConfig
}

// BoardConfig represents configuration options for a connected arduino
type BoardConfig struct {
	Serial string
	Dev    map[string]DeviceConfig
}

// DeviceConfig represents configuration options of a single device
type DeviceConfig struct {
	Mode     string
	Inverted bool
	Pins     []string
}

// RedisConfig represents configuration options for the redis broker
type RedisConfig struct {
	Host  string
	Topic string
}

//LoadConfig loads and parses a config file
func LoadConfig(filename string) (TomlConfig, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return TomlConfig{}, err
	}
	var config TomlConfig
	if err := toml.Unmarshal(buf, &config); err != nil {
		return TomlConfig{}, err
	}
	return config, nil
}
