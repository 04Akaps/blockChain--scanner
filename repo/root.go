package repo

import (
	"scanner/env"
	"scanner/repo/db"
	"scanner/repo/node"
)

type Repo struct {
	env *env.Env

	DB   *db.DB
	Node *node.Node
}

func NewRepo(env *env.Env) (*Repo, error) {
	r := &Repo{
		env: env,
	}

	var err error

	if r.DB, err = db.NewDB(env); err != nil {
		return nil, err
	} else if r.Node, err = node.NewNode(env); err != nil {
		return nil, err
	} else {
		return r, nil
	}

}
