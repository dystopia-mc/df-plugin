package npc

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/npc"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"math"

	_ "embed"
)

var (
	//go:embed plugin.toml
	config []byte
)

type command struct {
	m *plugin.Manager
}

func (c command) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p := npc.Create(npc.Settings{
		Name:       s.(*player.Player).NameTag(),
		Skin:       s.(*player.Player).Skin(),
		Position:   s.(*player.Player).Position(),
		Scale:      s.(*player.Player).Scale(),
		Vulnerable: true,
	}, tx, nil)

	p.SetMaxHealth(math.MaxFloat64)
	p.Heal(math.MaxFloat64, nil)

	p.Handle(c.m.PlayerHandler(p))
	o.Printf("created actor")
}

func Plugin() plugin.Plugin {
	m := plugin.M()

	task := m.NewTask(func(m *plugin.Manager) {
		cmd.Register(cmd.New("npc", "Creates an actor that you can punch", nil, command{m: m}))
	})

	return plugin.New(plugin.MustUnmarshalConfig(config), task)
}
