// package shared contains the shared interface definition
package iExec

import (
	"fmt"
	"net/rpc"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// // Create an hclog.Logger
// logger := hclog.New(&hclog.LoggerOptions{
// 	Name:   "plugin",
// 	Output: os.Stdout,
// 	Level:  hclog.Debug,
// })

// CdExecutor represents the interface for executing commands in a plugin.
type CdExecutor interface {
	// CdExec is a method that executes a command and returns the result.
	CdExec(jsonInput string) (string, error)
}

// Here is an implementation that talks over RPC
type CdExecutorRPCClient struct {
	logger hclog.Logger
	client *rpc.Client
}

func (g *CdExecutorRPCClient) CdExec(req string) (string, error) {
	g.logger.Info("CdExecutorRPCClient::CdExec()", "req", req)
	var resp string
	err := g.client.Call("Plugin.CdExec", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp, err
}

// Here is the RPC server that CdExecutorRPC talks to, conforming to
// the requirements of net/rpc
type CdExecutorRPCServer struct {
	logger hclog.Logger
	// This is the real implementation
	Impl CdExecutor
}

func (s *CdExecutorRPCServer) CdExec(args interface{}, resp *string) error {
	req := fmt.Sprintf("%v", args)
	s.logger.Info("CdExecutorRPCServer::CdExec()", "req", req)
	*resp, _ = s.Impl.CdExec(req)
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a CdExecutorRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return CdExecutorRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type CdExecutorPlugin struct {
	// Impl Injection
	Impl CdExecutor
}

func (p *CdExecutorPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CdExecutorRPCServer{Impl: p.Impl}, nil
}

func (CdExecutorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CdExecutorRPCClient{client: c}, nil
}
