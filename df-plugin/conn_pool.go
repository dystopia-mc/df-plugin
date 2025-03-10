package plugin

import (
	"github.com/df-mc/dragonfly/server/session"
	"sync"
)

type connPool struct {
	conn   map[string]session.Conn
	connMu sync.RWMutex
}

func newConnPool() *connPool {
	return &connPool{conn: make(map[string]session.Conn)}
}

func (cp *connPool) add(conn session.Conn) {
	cp.connMu.Lock()
	defer cp.connMu.Unlock()

	cp.conn[conn.IdentityData().DisplayName] = conn
}

func (cp *connPool) del(name string) {
	cp.connMu.Lock()
	defer cp.connMu.Unlock()
	delete(cp.conn, name)
}

func (cp *connPool) get(name string) session.Conn {
	cp.connMu.RLock()
	defer cp.connMu.RUnlock()
	return cp.conn[name]
}
