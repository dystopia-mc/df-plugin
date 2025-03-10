package plugin

import (
	"github.com/pelletier/go-toml"
	"os"
)

type Config struct {
	Plugin struct {
		Name,
		Description,
		Author string
	}
}

func NewConfig(name, author, description string) (c Config) {
	c.Plugin.Name = name
	c.Plugin.Author = author
	c.Plugin.Description = description

	return
}

func MustUnmarshalConfig(b []byte) Config {
	var c Config
	if err := toml.Unmarshal(b, &c); err != nil {
		panic(err)
	}

	return c
}

func MustParseConfig(path string) Config {
	c, err := ParseConfig(path)
	if err != nil {
		panic(err)
	}

	return c
}

func ParseConfig(path string) (Config, error) {
	var c Config

	b, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}

	if err := toml.Unmarshal(b, &c); err != nil {
		return c, err
	}

	return c, nil
}
