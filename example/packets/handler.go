package packets

import (
	"github.com/df-mc/dragonfly/server/player"
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type handler struct {
	plugin.NopPlayerHandler
}

func (h *handler) HandleClientPacket(ctx *player.Context, pk packet.Packet) {
	switch pk := pk.(type) {
	case *packet.Text:
		if pk.Message == ".test" {
			ctx.Cancel()
			ctx.Val().Messagef("test")
		}
	}
}
