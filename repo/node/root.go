package node

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"scanner/env"
	. "scanner/util"
)

type Node struct {
	env *env.Env

	client *ethclient.Client
}

func NewNode(env *env.Env) (*Node, error) {
	n := &Node{
		env: env,
	}

	var err error

	if n.client, err = ethclient.DialContext(Context(), env.Node.Dial); err != nil {
		return nil, err
	} else {
		return n, nil
	}

}

func (n *Node) GetLatestBlock() (uint64, error) {
	return n.client.BlockNumber(Context())
}

func (n *Node) GetBlockByNumber(block *big.Int) (*types.Block, error) {
	return n.client.BlockByNumber(Context(), block)
}

func (n *Node) GetChainID() (*big.Int, error) {
	return n.client.ChainID(Context())
}
