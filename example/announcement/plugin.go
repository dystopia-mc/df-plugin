package announcement

import (
	_ "embed"
	"github.com/df-mc/dragonfly/server/player/chat"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"time"
)

//go:embed plugin.toml
var config []byte

func Plugin() plugin.Plugin {
	m := plugin.M()

	task := m.NewTask(func(m *plugin.Manager) {
		// will send message to global chat every 10 seconds
		go func() {
			for range time.Tick(time.Second * 10) {
				_, _ = chat.Global.WriteString(getOrderedMessage())
			}
		}()
	})

	return plugin.New(plugin.MustUnmarshalConfig(config), task)
}

var current int

func getOrderedMessage() string {
	switch current {
	case 0:
		current++
		return "first"
	case 1:
		current++
		return "second"
	case 2:
		current++
		return "third"
	default:
		current = 0
		return "fourth, reset"
	}
}
