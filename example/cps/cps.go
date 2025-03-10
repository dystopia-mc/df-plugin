package cps

import (
	"sync"
	"time"
)

type cps struct {
	clicks []time.Time
	mu     sync.Mutex
}

func (c *cps) add() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.clicks) >= 100 {
		c.clicks = c.clicks[1:]
	}

	c.clicks = append(c.clicks, time.Now())
}

func (c *cps) CPS() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calculate()
}

func (c *cps) calculate() int {
	var clicks int

	for _, t := range c.clicks {
		if time.Since(t) <= time.Second {
			clicks++
		}
	}

	return clicks
}
