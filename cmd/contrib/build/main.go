package main

import (
	"errors"

	"github.com/1Password/shell-plugins/plugins"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/hashicorp/go-plugin"

	"github.com/1Password/shell-plugins/sdk/rpc/proto"
	"github.com/1Password/shell-plugins/sdk/rpc/server"
)

// PluginName is the name of the plugin to serve. This should be set during building using
// -ldflags="-X 'main.PluginName=<plugin-name>'"
var PluginName string

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  proto.Version,
			MagicCookieKey:   proto.MagicCookieKey,
			MagicCookieValue: proto.MagicCookieValue,
		},
		Plugins: plugin.PluginSet{
			"plugin": &server.RPCPlugin{RPCPlugin: func() (schema.Plugin, error) {
				if PluginName == "" {
					return schema.Plugin{}, errors.New("plugin not set")
				}
				p, err := plugins.Get(PluginName)
				if err != nil {
					return schema.Plugin{}, err
				}
				return p, nil
			}},
		},
	})
}
