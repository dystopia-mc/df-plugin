package knockback

import (
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"time"

	_ "embed"
)

var (
	//go:embed plugin.toml
	config []byte

	height, force atomic.Float64
	immunity      atomic.Duration
)

func init() {
	height.Store(0.46)
	force.Store(0.2)

	immunity.Store(420 * time.Millisecond)
}

type handler struct {
	plugin.NopPlayerHandler
}

func (h handler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	switch src.(type) {
	case entity.AttackDamageSource:
		*attackImmunity = immunity.Load()
	default:
		ctx.Cancel()
	}
}

func (h handler) HandleAttackEntity(ctx *player.Context, e world.Entity, f, he *float64, critical *bool) {
	*he = height.Load()
	*f = force.Load()
}

func Plugin() plugin.Plugin {
	m := plugin.M()

	task := m.NewTask(func(m *plugin.Manager) {
		cmd.Register(cmd.New("knockback", "Allows you to customize knockback", nil, knockbackHeight{}, knockbackForce{}, knockbackImmunity{}))
	})

	return plugin.New(plugin.MustUnmarshalConfig(config), task).WithHandlers(handlers()...)
}

func handlers() []plugin.PlayerHandler {
	return []plugin.PlayerHandler{
		handler{},
	}
}
