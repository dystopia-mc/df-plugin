package plugin

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"log/slog"
)

type IncomingConnectionFunc func(c *minecraft.Conn)
type DisconnectFunc func(c *minecraft.Conn)

// Listener is an implementation of server.Listener interface.
type Listener struct {
	l *minecraft.Listener

	f IncomingConnectionFunc
	d DisconnectFunc
}

func (l *Listener) Accept() (session.Conn, error) {
	conn, err := l.l.Accept()
	if err != nil {
		return nil, err
	}

	l.f(conn.(*minecraft.Conn))
	return conn.(session.Conn), err
}

func (l *Listener) AddResourcePack(p *resource.Pack) {
	l.l.AddResourcePack(p)
}

func (l *Listener) Disconnect(conn session.Conn, reason string) error {
	l.d(conn.(*minecraft.Conn))
	return l.l.Disconnect(conn.(*minecraft.Conn), reason)
}

func (l *Listener) Close() error {
	return l.l.Close()
}

func newListener(port uint16, authEnabled bool, prov minecraft.ServerStatusProvider, l *slog.Logger, f IncomingConnectionFunc, d DisconnectFunc, packs []*resource.Pack, acceptedProtocols ...minecraft.Protocol) (*Listener, error) {
	if f == nil {
		f = func(c *minecraft.Conn) {}
	}
	if d == nil {
		d = func(c *minecraft.Conn) {}
	}

	li, err := minecraft.ListenConfig{
		ErrorLog:               l,
		AuthenticationDisabled: !authEnabled,
		StatusProvider:         prov,
		ResourcePacks:          packs,
		TexturePacksRequired:   len(packs) > 0,
		AcceptedProtocols:      acceptedProtocols,
	}.Listen("queryraknet", fmt.Sprintf(":%v", port))

	if err != nil {
		return nil, err
	}

	return &Listener{l: li, f: f, d: d}, nil
}
