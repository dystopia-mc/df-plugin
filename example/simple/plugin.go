package simple

import (
	plugin "github.com/k4ties/df-plugin/df-plugin"
	"log"
)

var config = func() plugin.Config {
	c := plugin.Config{}
	c.Plugin.Name = "Simple example"
	c.Plugin.Author = "k4ties"
	c.Plugin.Description = "Does nothing."

	return c
}()

func Plugin() plugin.Plugin {
	m := plugin.M()

	task := m.NewTask(func(m *plugin.Manager) {
		log.Printf("Hello from plugin! Nether spawn position is %v", m.Srv().Nether().Spawn().String())
	})

	return plugin.New(config, task)
}
