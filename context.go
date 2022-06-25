package sanji

import (
	"math"

	"github.com/bitleak/lmstfy/client"
)

const abortIndex int8 = math.MaxInt8 >> 1

type Context struct {
	index       int8
	handlers    []JobFunc
	acknowledge bool

	Queue string
	Job   *client.Job
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers))+1 {
		c.handlers[c.index-1](c)
		c.index++
	}
}

func (c *Context) IsAbort() bool {
	return c.index >= abortIndex
}

func (c *Context) Abort() {
	c.index = abortIndex
}
