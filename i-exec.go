// package shared contains the shared interface definition
package iExec

import (
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// CdExecutor represents the interface for executing commands in a plugin.
type CdExecutor interface {
	// CdExec is a method that executes a command and returns the result.
	CdExec(jsonInput string) (string, error)
}

// Here is an implementation that talks over RPC
type CdExecutorRPCClient struct {
	client *rpc.Client
}

func (g *CdExecutorRPCClient) CdExec(req string) (string, error) {
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
	// This is the real implementation
	Impl CdExecutor
}

func (s *CdExecutorRPCServer) CdExec(args interface{}, resp *string) (string, error) {
	fmt.Println("CdExecutorRPCServer::args:", args)
	// req := args.(string)
	*resp, _ = s.Impl.CdExec("xxxx")
	return *resp, nil
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
