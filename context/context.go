package context

import (
	"context"
)

//Context holds the connection context struct, cancel func and user value
type Context struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	User   string
}

//NewContext creates a new Context struct
func NewContext(user string) Context {
	ctx, cancel := context.WithCancel(context.Background())
	return Context{
		ctx, cancel, user,
	}
}
