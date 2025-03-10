package plugin

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"math"
	"reflect"
	"sync"
	"unsafe"
)

var handlers = struct {
	v  map[string]PacketHandler
	mu sync.RWMutex
}{
	v: make(map[string]PacketHandler),
}

type PacketHandler interface {
	HandleClientPacket(ctx *player.Context, pk packet.Packet)
	HandleServerPacket(ctx *player.Context, pk packet.Packet)
}

func hook(name string, handler PacketHandler) {
	handlers.mu.Lock()
	defer handlers.mu.Unlock()

	handlers.v[name] = handler
}

func unHook(name string) {
	handlers.mu.Lock()
	defer handlers.mu.Unlock()

	delete(handlers.v, name)
}

type conn struct {
	session.Conn
	p *player.Player
}

func (c *conn) ReadPacket() (packet.Packet, error) {
	pkt, err := c.Conn.ReadPacket()
	if err != nil {
		return pkt, err
	}

	ctx := event.C(c.p)
	handlers.mu.RLock()
	for _, h := range handlers.v {
		h.HandleClientPacket(ctx, pkt)
	}
	handlers.mu.RUnlock()

	if ctx.Cancelled() {
		return NopPacket{}, nil
	}
	return pkt, nil
}

func (c *conn) WritePacket(pk packet.Packet) error {
	ctx := event.C(c.p)
	handlers.mu.RLock()
	for _, h := range handlers.v {
		h.HandleServerPacket(ctx, pk)
	}
	handlers.mu.RUnlock()

	if ctx.Cancelled() {
		return nil
	}
	return c.Conn.WritePacket(pk)
}

func intercept(p *player.Player, h PacketHandler, m *Manager) {
	hook(p.Name(), h)
	s := p.Data().Session

	c := fetchPrivateField[session.Conn](s, "conn")
	cn := &conn{c, p}
	updatePrivateField[session.Conn](s, "conn", cn)

	m.cp.add(cn)
}

type NopPacket struct{}

func (NopPacket) ID() uint32 {
	return math.MaxUint32
}

func (NopPacket) Marshal(protocol.IO) {}

// updatePrivateField sets a private field of a session to the value passed.
func updatePrivateField[T any](s *session.Session, name string, value T) {
	reflectedValue := reflect.ValueOf(s).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)

	privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

	privateFieldValue.Set(reflect.ValueOf(value))
}

// fetchPrivateField fetches a private field of a session.
func fetchPrivateField[T any](s *session.Session, name string) T {
	reflectedValue := reflect.ValueOf(s).Elem()
	privateFieldValue := reflectedValue.FieldByName(name)
	privateFieldValue = reflect.NewAt(privateFieldValue.Type(), unsafe.Pointer(privateFieldValue.UnsafeAddr())).Elem()

	return privateFieldValue.Interface().(T)
}
