package server

import (
	"errors"
	"net/rpc"

	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/hashicorp/go-plugin"
)

// RPCPlugin is an implementation of the github.com/hashicorp/go-plugin#Plugin interface indicating how
// to serve a Shell Plugin over RPC using go-plugin.
type RPCPlugin struct {
	RPCPlugin func() (schema.Plugin, error)
}

// Server registers the RPC provider server with the RPC server that
// go-plugin is setting up.
func (p *RPCPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	pl, err := p.RPCPlugin()
	if err != nil {
		return nil, err
	}

	return newServer(pl), nil
}

// Client always returns an error; we're only implementing a server.
func (p *RPCPlugin) Client(*plugin.MuxBroker, *rpc.Client) (interface{}, error) {
	return nil, errors.New("only server is implemented")
}
