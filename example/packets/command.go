package packets

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type command struct {
	m *plugin.Manager
}

func (c command) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	p := s.(*player.Player)

	if conn, ok := c.m.Conn(p.Name()); ok {
		_ = conn.WritePacket(&packet.LevelEvent{EventType: 16404, Position: mgl64to32(p.Position())}) // spawns heart particle
		o.Printf("success")

		return
	}

	o.Error("cannot get your connection")
}

func (c command) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

func mgl64to32(s mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{
		float32(s.X()),
		float32(s.Y()),
		float32(s.Z()),
	}
}
