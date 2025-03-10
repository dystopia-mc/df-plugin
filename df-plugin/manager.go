package plugin

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/gameparrot/goquery"
	"github.com/gameparrot/goqueryraknet"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"golang.org/x/exp/maps"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var m *Manager

func M() *Manager {
	return m
}

type ManagerConfig struct {
	Logger     *slog.Logger
	UserConfig server.UserConfig

	SubName string
	Packs   []*resource.Pack

	Protocols    []minecraft.Protocol
	QueryHandler QueryHandler
}

func (c ManagerConfig) fillDefaults() ManagerConfig {
	if c.Logger == nil {
		c.Logger = slog.Default()
	}
	if _, err := c.UserConfig.Config(c.Logger); err != nil {
		c.UserConfig = server.DefaultConfig()
	}
	if c.Packs == nil {
		c.Packs = []*resource.Pack{}
	}
	if c.Protocols == nil {
		c.Protocols = []minecraft.Protocol{}
	}
	if c.QueryHandler == nil {
		c.QueryHandler = NopQueryHandler{}
	}

	return c
}

func NewManager(cfg ManagerConfig) *Manager {
	cfg = cfg.fillDefaults()
	l := cfg.Logger

	m = &Manager{l: l, plugins: make(map[string]Plugin), cp: newConnPool()}

	port, ok := parsePort(cfg.UserConfig.Network.Address)
	if !ok {
		port = 19132
		l.Error("Failed to parse port in Network.Address, 19132 is now default")
	}

	onLogin := func(conn *minecraft.Conn) {
		for _, h := range m.allHandlers() {
			ctx := event.C(session.Conn(conn))
			if h.HandleLogin(ctx); ctx.Cancelled() {
				_ = conn.Close()
			}
		}
	}

	srvConf, err := cfg.UserConfig.Config(l)
	if err != nil {
		panic(err)
	}
	m.srvConf = srvConf
	m.srvConf.Listeners = []func(conf server.Config) (server.Listener, error){
		func(conf server.Config) (server.Listener, error) {
			li := m.createListener(port, cfg, onLogin)
			m.li = li

			return li, nil
		},
	}

	m.srv = m.srvQuery(m.srvConf)
	return m
}

type Manager struct {
	l   *slog.Logger
	srv *server.Server

	srvConf server.Config
	li      *Listener

	plugins   map[string]Plugin
	pluginsMu sync.RWMutex

	cp *connPool

	createNetworkOnce sync.Once

	started   atomic.Bool
	statusCmd atomic.Bool

	qh QueryHandler
}

func (m *Manager) HandleQuery(h QueryHandler) {
	if h == nil {
		h = NopQueryHandler{}
	}

	m.qh = h
}

func (m *Manager) createQueryNetwork(name string, q *goquery.QueryServer) {
	m.createNetworkOnce.Do(func() {
		goqueryraknet.CreateGophertunnelNetwork(name, q)
	})
}

func (m *Manager) srvQuery(c server.Config) *server.Server {
	q := goquery.New(map[string]string{}, []string{})
	m.createQueryNetwork("queryraknet", q)
	srv := c.New()

	q.SetInfoFunc(handleQuery(srv, m, m.qh))
	return srv
}

func (m *Manager) Online(i uuid.UUID) bool {
	for p := range m.Srv().Players(nil) {
		if p.UUID() == i {
			return true
		}
	}

	return false
}

func (m *Manager) WithPlayerProvider(p player.Provider) *Manager {
	m.srvConf.PlayerProvider = p
	m.srv = m.srvQuery(m.srvConf)
	return m
}

func (m *Manager) NewTask(f TaskFunc) *Task {
	if f == nil {
		f = func(*Manager) {}
	}
	return &Task{m: m, f: f}
}

func (m *Manager) ServerConfig() server.Config {
	return m.srvConf
}

func (m *Manager) createListener(port uint16, cfg ManagerConfig, onLogin func(conn2 *minecraft.Conn)) *Listener {
	li, err := newListener(
		port,
		cfg.UserConfig.Server.AuthEnabled,
		minecraft.NewStatusProvider(cfg.UserConfig.Server.Name, cfg.SubName),
		cfg.Logger,
		onLogin,
		nil,
		cfg.Packs,
		cfg.Protocols...,
	)
	if err != nil {
		panic(err)
	}

	return li
}

func (m *Manager) ToggleStatusCommand() {
	m.statusCmd.Store(!m.statusCmd.Load())
}

func (m *Manager) Plugins() map[string]Plugin {
	m.pluginsMu.RLock()
	defer m.pluginsMu.RUnlock()
	return m.plugins
}

func (m *Manager) Register(pl ...Plugin) {
	for _, p := range pl {
		m.registerPlugin(p)
	}
}

func (m *Manager) ListenServer() {
	if m.statusCmd.Load() {
		registerStatusCommand()
	}

	for _, t := range m.allTasks() {
		t.do()
	}

	m.srv.Listen()
	m.started.Store(true)

	for pl := range m.srv.Accept() {
		h := m.PlayerHandler(pl)

		pl.Handle(h)
		intercept(pl, h, m)

		for _, h := range m.allHandlers() {
			h.HandleSpawn(pl)
		}
	}
}

func (m *Manager) Started() bool {
	return m.started.Load()
}

func (m *Manager) PlayerHandler(p *player.Player) PlayerHandler {
	return NewPlayerHandler(p, m)
}

func (m *Manager) Srv() *server.Server {
	return m.srv
}

func (m *Manager) Conn(name string) (session.Conn, bool) {
	c := m.cp.get(name)
	return c, c != nil
}

func (m *Manager) PluginByIndex(i int) (Plugin, bool) {
	for _, p := range m.Plugins() {
		if p.Index() == i {
			return p, true
		}
	}

	return Plugin{}, false
}

func (m *Manager) Listener() *Listener {
	return m.li
}

func (m *Manager) PluginsSlice() []Plugin {
	return maps.Values(m.Plugins())
}

func (m *Manager) Logger() *slog.Logger {
	return m.l
}

func (m *Manager) Plugin(name string) (Plugin, bool) {
	m.pluginsMu.RLock()
	defer m.pluginsMu.RUnlock()

	p, ok := m.plugins[name]
	return p, ok
}

func (m *Manager) registerPlugin(pl Plugin) {
	m.pluginsMu.Lock()
	defer m.pluginsMu.Unlock()

	m.plugins[pl.Name()] = pl
	m.l.Info("Loaded plugin.", "plugin", pl.Name(), "src", "df-plugin")
}

func (m *Manager) allHandlers() (handlers []PlayerHandler) {
	m.pluginsMu.RLock()
	defer m.pluginsMu.RUnlock()

	for _, pl := range m.plugins {
		for _, h := range pl.handlers {
			handlers = append(handlers, h)
		}
	}

	return
}

func (m *Manager) allTasks() (tasks []*Task) {
	m.pluginsMu.RLock()
	defer m.pluginsMu.RUnlock()

	for _, pl := range m.plugins {
		tasks = append(tasks, pl.task)
	}

	return
}

func parsePort(a string) (uint16, bool) {
	addr, err := strconv.Atoi(strings.Split(a, ":")[1])
	if err != nil {
		return 0, false
	}

	return uint16(addr), true
}
