package plugin

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
	"time"
)

// Handler is used internally by plugin manager.
type Handler struct {
	p *player.Player
	m *Manager
}

func NewPlayerHandler(p *player.Player, m *Manager) Handler {
	return Handler{p, m}
}

type PlayerHandler interface {
	player.Handler
	HandleLogin(ctx *event.Context[session.Conn])
	HandleSpawn(p *player.Player)
	HandleClientPacket(ctx *player.Context, pk packet.Packet)
	HandleServerPacket(ctx *player.Context, pk packet.Packet)
}

func (h Handler) HandleMove(ctx *player.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	for _, h := range h.m.allHandlers() {
		h.HandleMove(ctx, newPos, newRot)
	}
}
func (h Handler) HandleJump(p *player.Player) {
	for _, h := range h.m.allHandlers() {
		h.HandleJump(p)
	}
}
func (h Handler) HandleTeleport(ctx *player.Context, pos mgl64.Vec3) {
	for _, h := range h.m.allHandlers() {
		h.HandleTeleport(ctx, pos)
	}
}
func (h Handler) HandleChangeWorld(p *player.Player, before, after *world.World) {
	for _, h := range h.m.allHandlers() {
		h.HandleChangeWorld(p, before, after)
	}
}
func (h Handler) HandleToggleSprint(ctx *player.Context, after bool) {
	for _, h := range h.m.allHandlers() {
		h.HandleToggleSprint(ctx, after)
	}
}
func (h Handler) HandleToggleSneak(ctx *player.Context, after bool) {
	for _, h := range h.m.allHandlers() {
		h.HandleToggleSneak(ctx, after)
	}
}
func (h Handler) HandleChat(ctx *player.Context, message *string) {
	for _, h := range h.m.allHandlers() {
		h.HandleChat(ctx, message)
	}
}
func (h Handler) HandleFoodLoss(ctx *player.Context, from int, to *int) {
	for _, h := range h.m.allHandlers() {
		h.HandleFoodLoss(ctx, from, to)
	}
}
func (h Handler) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
	for _, h := range h.m.allHandlers() {
		h.HandleHeal(ctx, health, src)
	}
}
func (h Handler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	for _, h := range h.m.allHandlers() {
		h.HandleHurt(ctx, damage, immune, attackImmunity, src)
	}
}
func (h Handler) HandleDeath(p *player.Player, src world.DamageSource, keepInv *bool) {
	for _, h := range h.m.allHandlers() {
		h.HandleDeath(p, src, keepInv)
	}
}
func (h Handler) HandleRespawn(p *player.Player, pos *mgl64.Vec3, w **world.World) {
	for _, h := range h.m.allHandlers() {
		h.HandleRespawn(p, pos, w)
	}
}
func (h Handler) HandleSkinChange(ctx *player.Context, skin *skin.Skin) {
	for _, h := range h.m.allHandlers() {
		h.HandleSkinChange(ctx, skin)
	}
}
func (h Handler) HandleFireExtinguish(ctx *player.Context, pos cube.Pos) {
	for _, h := range h.m.allHandlers() {
		h.HandleFireExtinguish(ctx, pos)
	}
}
func (h Handler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	for _, h := range h.m.allHandlers() {
		h.HandleStartBreak(ctx, pos)
	}
}
func (h Handler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	for _, h := range h.m.allHandlers() {
		h.HandleBlockBreak(ctx, pos, drops, xp)
	}
}
func (h Handler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	for _, h := range h.m.allHandlers() {
		h.HandleBlockPlace(ctx, pos, b)
	}
}
func (h Handler) HandleBlockPick(ctx *player.Context, pos cube.Pos, b world.Block) {
	for _, h := range h.m.allHandlers() {
		h.HandleBlockPick(ctx, pos, b)
	}
}
func (h Handler) HandleItemUse(ctx *player.Context) {
	for _, h := range h.m.allHandlers() {
		h.HandleItemUse(ctx)
	}
}
func (h Handler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	for _, h := range h.m.allHandlers() {
		h.HandleItemUseOnBlock(ctx, pos, face, clickPos)
	}
}
func (h Handler) HandleItemUseOnEntity(ctx *player.Context, e world.Entity) {
	for _, h := range h.m.allHandlers() {
		h.HandleItemUseOnEntity(ctx, e)
	}
}
func (h Handler) HandleItemRelease(ctx *player.Context, item item.Stack, dur time.Duration) {
	for _, h := range h.m.allHandlers() {
		h.HandleItemRelease(ctx, item, dur)
	}
}
func (h Handler) HandleItemConsume(ctx *player.Context, item item.Stack) {
	for _, h := range h.m.allHandlers() {
		h.HandleItemConsume(ctx, item)
	}
}
func (h Handler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	for _, h := range h.m.allHandlers() {
		h.HandleAttackEntity(ctx, e, force, height, critical)
	}
}
func (h Handler) HandleExperienceGain(ctx *player.Context, amount *int) {
	for _, h := range h.m.allHandlers() {
		h.HandleExperienceGain(ctx, amount)
	}
}
func (h Handler) HandlePunchAir(ctx *player.Context) {
	for _, h := range h.m.allHandlers() {
		h.HandlePunchAir(ctx)
	}
}
func (h Handler) HandleSignEdit(ctx *player.Context, pos cube.Pos, frontSide bool, oldText, newText string) {
	for _, h := range h.m.allHandlers() {
		h.HandleSignEdit(ctx, pos, frontSide, oldText, newText)
	}
}
func (h Handler) HandleLecternPageTurn(ctx *player.Context, pos cube.Pos, oldPage int, newPage *int) {
	for _, h := range h.m.allHandlers() {
		h.HandleLecternPageTurn(ctx, pos, oldPage, newPage)
	}
}
func (h Handler) HandleItemDamage(ctx *player.Context, i item.Stack, damage int) {
	for _, h := range h.m.allHandlers() {
		h.HandleItemDamage(ctx, i, damage)
	}
}
func (h Handler) HandleItemPickup(ctx *player.Context, i *item.Stack) {
	for _, h := range h.m.allHandlers() {
		h.HandleItemPickup(ctx, i)
	}
}
func (h Handler) HandleHeldSlotChange(ctx *player.Context, from, to int) {
	for _, h := range h.m.allHandlers() {
		h.HandleHeldSlotChange(ctx, from, to)
	}
}
func (h Handler) HandleItemDrop(ctx *player.Context, s item.Stack) {
	for _, h := range h.m.allHandlers() {
		h.HandleItemDrop(ctx, s)
	}
}
func (h Handler) HandleTransfer(ctx *player.Context, addr *net.UDPAddr) {
	for _, h := range h.m.allHandlers() {
		h.HandleTransfer(ctx, addr)
	}
}
func (h Handler) HandleCommandExecution(ctx *player.Context, command cmd.Command, args []string) {
	for _, h := range h.m.allHandlers() {
		h.HandleCommandExecution(ctx, command, args)
	}
}
func (h Handler) HandleQuit(p *player.Player) {
	for _, h := range h.m.allHandlers() {
		h.HandleQuit(p)
	}
	unHook(p.Name())
	h.m.cp.del(p.Name())
}
func (h Handler) HandleDiagnostics(p *player.Player, d session.Diagnostics) {
	for _, h := range h.m.allHandlers() {
		h.HandleDiagnostics(p, d)
	}
}
func (h Handler) HandleSpawn(p *player.Player) {
	for _, h := range h.m.allHandlers() {
		h.HandleSpawn(p)
	}
}
func (h Handler) HandleLogin(ctx *event.Context[session.Conn]) {
	for _, h := range h.m.allHandlers() {
		h.HandleLogin(ctx)
	}
}
func (h Handler) HandleClientPacket(ctx *player.Context, pk packet.Packet) {
	for _, h := range h.m.allHandlers() {
		h.HandleClientPacket(ctx, pk)
	}
}
func (h Handler) HandleServerPacket(ctx *player.Context, pk packet.Packet) {
	for _, h := range h.m.allHandlers() {
		h.HandleServerPacket(ctx, pk)
	}
}

type NopPlayerHandler struct {
	player.NopHandler
}

func (h NopPlayerHandler) HandleLogin(*event.Context[session.Conn])          {}
func (h NopPlayerHandler) HandleSpawn(*player.Player)                        {}
func (h NopPlayerHandler) HandleClientPacket(*player.Context, packet.Packet) {}
func (h NopPlayerHandler) HandleServerPacket(*player.Context, packet.Packet) {}
