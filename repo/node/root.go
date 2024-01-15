package node

import (
	"github.com/ethereum/go-ethereum/ethclient"
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
