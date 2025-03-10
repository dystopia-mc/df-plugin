package packets

import (
	_ "embed"
	"github.com/df-mc/dragonfly/server/cmd"
	plugin "github.com/k4ties/df-plugin/df-plugin"
)

//go:embed plugin.toml
var config []byte

func Plugin() plugin.Plugin {
	m := plugin.M()

	task := m.NewTask(func(m *plugin.Manager) {
		cmd.Register(cmd.New("pk", "packets example", nil, command{m: m}))
	})

	return plugin.New(plugin.MustUnmarshalConfig(config), task).WithHandlers(&handler{})
}
