package main

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/k4ties/df-plugin/example/announcement"
	"github.com/k4ties/df-plugin/example/cps"
	"github.com/k4ties/df-plugin/example/knockback"
	"github.com/k4ties/df-plugin/example/npc"
	"github.com/k4ties/df-plugin/example/packets"
	"github.com/k4ties/df-plugin/example/simple"
	"github.com/pelletier/go-toml"
	"log/slog"
	"os"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	m := plugin.NewManager(slog.Default(), "", mustReadConfig("config.toml"))
	m.Register(allPlugins()...)

	m.ToggleStatusCommand()
	m.ListenServer()
}

func allPlugins() []plugin.Plugin {
	return []plugin.Plugin{
		announcement.Plugin(),
		knockback.Plugin(),
		npc.Plugin(),
		cps.Plugin(),
		simple.Plugin(),
		packets.Plugin(),
	}
}

func mustReadConfig(path string) server.UserConfig {
	c := server.DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		b, err := toml.Marshal(c)
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(path, b, 0644); err != nil {
			panic(err)
		}
	}

	if err := toml.Unmarshal(data, &c); err != nil {
		panic(err)
	}

	return c
}
