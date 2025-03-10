package cps

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"

	_ "embed"
)

//go:embed plugin.toml
var config []byte

func Plugin() plugin.Plugin {
	m := plugin.M()

	task := m.NewTask(func(m *plugin.Manager) {
		// ...
	})

	return plugin.New(plugin.MustUnmarshalConfig(config), task, newCpsHandler())
}

type cpsHandler struct {
	plugin.NopPlayerHandler
	c *cps
}

func newCpsHandler() *cpsHandler {
	return &cpsHandler{c: &cps{}}
}

func (c cpsHandler) HandleAttackEntity(ctx *player.Context, _ world.Entity, _, _ *float64, _ *bool) {
	c.c.add()
	showCPS(ctx.Val(), c.c)
}

func (c cpsHandler) HandlePunchAir(ctx *player.Context) {
	c.c.add()
	showCPS(ctx.Val(), c.c)
}

func showCPS(p *player.Player, c *cps) {
	p.SendTip(c.CPS())
}
