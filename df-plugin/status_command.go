package plugin

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"runtime"
	"strconv"
	"strings"
)

type status struct {
	m      *Manager
	Plugin cmd.Optional[pluginEnum] `cmd:"plugin"`
}

func (c status) Run(src cmd.Source, _ *cmd.Output, _ *world.Tx) {
	if selected, ok := c.Plugin.Load(); ok {
		if s, ok := src.(form.Submitter); ok {
			i, _ := strconv.Atoi(string(selected))

			plugin, _ := c.m.PluginByIndex(i)
			formatText := func(s string, a ...any) string {
				return fmt.Sprintf("<dark-grey>|</dark-grey> %s", fmt.Sprintf(s, a...))
			}

			format := "Plugin %s\n\n%s\n%s"

			s.SendForm(form.NewMenu(
				statusForm{
					Close: form.Button{Text: "Close"},
				},
				"Status of a plugin",
			).WithBody(text.Colourf(format, plugin.Name(), formatText("Description: <grey>%s</grey>\n", plugin.Description()), formatText("Author: <grey>%s</grey>", plugin.Author()))))
		}

		return
	}

	format := "Server is running <grey>df-plugin.</grey>\nGoroutines: <grey>" + strconv.Itoa(runtime.NumGoroutine()) + "</grey>\n\nCurrently have <grey>%d</grey> plugin(s):\n\n%s"

	var names []string
	for i := 1; i <= len(c.m.Plugins()); i++ {
		if pl, ok := c.m.PluginByIndex(i); ok {
			names = append(names, text.Colourf("<dark-grey>%d)</dark-grey> %s, <dark-grey>by %s</dark-grey>", i, pl.Name(), pl.Author()))
		}
	}

	if s, ok := src.(form.Submitter); ok {
		s.SendForm(form.NewMenu(statusForm{
			Close: form.Button{
				Text: "Close",
			},
		}, "Status of the server").WithBody(text.Colourf(format, len(names), strings.Join(names, "\n"))))
	}
}

type statusForm struct {
	Close form.Button
}

func (f statusForm) Submit(form.Submitter, form.Button, *world.Tx) {}

type pluginEnum string

func (pluginEnum) Type() string {
	return "plugin"
}

func (pluginEnum) Options(_ cmd.Source) []string {
	var ops []string

	for i := 1; i <= len(M().Plugins()); i++ {
		ops = append(ops, strconv.Itoa(i))
	}

	return ops
}

func registerStatusCommand() {
	if _, ok := cmd.ByAlias("status"); ok {
		panic("there is already a status command. please delete it or disable /status command")
	}

	cmd.Register(cmd.New("status", "Status of the server", nil, status{m: m}))
}
