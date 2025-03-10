package plugin

import (
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/player"
	"maps"
	"slices"
	"sync"
)

type WithHandler interface {
	Plugin() Plugin
	Handlers() []PlayerHandler
}

type WithoutHandler interface {
	Plugin() Plugin
}

type TaskFunc func(m *Manager)

type Task struct {
	f TaskFunc
	m *Manager
}

func (t *Task) do() {
	go t.f(t.m)
}

func New(c Config, task *Task, handlers ...PlayerHandler) Plugin {
	pool.amount.Add(1)
	index := int(pool.amount.Load())

	pl := Plugin{task: task, handlers: handlers, Config: c, index: index}

	pool.pluginsMu.Lock()
	pool.plugins[index] = pl
	pool.pluginsMu.Unlock()

	return pl
}

var pool = struct {
	amount atomic.Int64

	plugins   map[int]Plugin
	pluginsMu sync.RWMutex
}{
	plugins: make(map[int]Plugin),
}

func RegisteredPlugins() []Plugin {
	pool.pluginsMu.Lock()
	defer pool.pluginsMu.Unlock()
	return slices.Collect(maps.Values(pool.plugins))
}

type Plugin struct {
	Config

	task     *Task
	handlers []PlayerHandler

	index int
}

func (p Plugin) HasHandlers() bool {
	return len(p.handlers) > 0
}

func (p Plugin) WithHandlers(handlers ...PlayerHandler) Plugin {
	p.handlers = append(p.handlers, handlers...)
	return p
}

func (p Plugin) Contains(handler player.Handler) bool {
	for _, h := range p.handlers {
		if h == handler {
			return true
		}
	}

	return false
}

func (p Plugin) Equals(pl Plugin) bool {
	return p.Config == pl.Config &&
		p.task == pl.task // cannot equal handlers
}

func (p Plugin) Index() int {
	return p.index
}

func (p Plugin) Name() string {
	return p.Config.Plugin.Name
}

func (p Plugin) Author() string {
	return p.Config.Plugin.Author
}

func (p Plugin) Description() string {
	return p.Config.Plugin.Description
}
