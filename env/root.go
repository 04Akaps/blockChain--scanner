package env

import (
	"github.com/naoina/toml"
	"os"
)

type Env struct {
	DB struct {
		Uri string
		DB  string

		Block string
		Tx    string
	}

	Node struct {
		Dial       string
		StartBlock uint64
	}

	Log struct {
		LogName string
	}
}

func NewEnv(f string) *Env {
	c := new(Env)

	if f, err := os.Open(f); err != nil {
		panic(err)
	} else if err = toml.NewDecoder(f).Decode(c); err != nil {
		panic(err)
	} else {
		return c
	}

}
