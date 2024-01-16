package app

import (
	"scanner/env"
	"scanner/log"
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
		log.ErrLog(err.Error())
	} else {

	}

}
