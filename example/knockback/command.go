package knockback

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

type knockbackHeight struct {
	Height cmd.SubCommand `cmd:"height"`
	Set    cmd.SubCommand `cmd:"set"`
	Value  float64        `cmd:"value"`
}

type knockbackForce struct {
	Force cmd.SubCommand `cmd:"force"`
	Set   cmd.SubCommand `cmd:"set"`
	Value float64        `cmd:"value"`
}

type knockbackImmunity struct {
	Force cmd.SubCommand `cmd:"immunity"`
	Set   cmd.SubCommand `cmd:"set"`
	Value int            `cmd:"milliseconds"`
}

func (c knockbackHeight) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	height.Store(c.Value)
	o.Printf("knockback height is now %v", c.Value)
}

func (c knockbackForce) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	force.Store(c.Value)
	o.Printf("knockback force is now %v", c.Value)
}

func (c knockbackImmunity) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	immunity.Store(time.Millisecond * time.Duration(c.Value))
	o.Printf("knockback immunity duration (in ms) is now %v", c.Value)
}
