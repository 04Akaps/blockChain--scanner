package app

import (
	"fmt"
	"scanner/env"
	"scanner/repo"
)

type App struct {
	env *env.Env

	repo *repo.Repo
}

func NewApp(e *env.Env) {
	a := &App{
		env: e,
	}

	var err error

	if a.repo, err = repo.NewRepo(e); err != nil {
		// TODO Log
		fmt.Println(err)
	}

	fmt.Println(a)
}
