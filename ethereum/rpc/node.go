package rpc

import (
	"github.com/ethereum/go-ethereum/node"

	"github.com/artela-network/artela/ethereum/rpc/types"
)

// Node Wrapers Ethereum Node
type Node struct {
	*node.Node
}

// Node is an implement of NetworkingStack
var _ types.NetworkingStack = (*Node)(nil)

// Node creates a new NetworkingStack instance.
func NewNode(config *node.Config) (types.NetworkingStack, error) {
	node, err := node.New(config)
	if err != nil {
		return nil, err
	}

	return &Node{
		Node: node,
	}, nil
}

// ExtRPCEnabled returns whether or not the external RPC service is enabled.
func (n *Node) ExtRPCEnabled() bool {
	return n.Node.Config().ExtRPCEnabled()
}

// Start starts the networking stack.
func (n *Node) Start() error {
	return n.Node.Start()
}

// DefaultConfig returns the default configuration for the provider.
func DefaultGethNodeConfig() *node.Config {
	nodeCfg := node.DefaultConfig
	nodeCfg.P2P.NoDiscovery = true
	nodeCfg.P2P.MaxPeers = 0
	nodeCfg.Name = clientIdentifier
	nodeCfg.HTTPModules = append(nodeCfg.HTTPModules, "eth", "web3", "net", "txpool", "debug")
	nodeCfg.WSModules = append(nodeCfg.WSModules, "eth")
	nodeCfg.HTTPHost = "0.0.0.0"
	nodeCfg.WSHost = ""
	nodeCfg.WSOrigins = []string{"*"}
	nodeCfg.HTTPCors = []string{"*"}
	nodeCfg.HTTPVirtualHosts = []string{"*"}
	nodeCfg.GraphQLCors = []string{"*"}
	nodeCfg.GraphQLVirtualHosts = []string{"*"}
	return &nodeCfg
}
