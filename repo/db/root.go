package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"scanner/env"
	. "scanner/util"
)

type DB struct {
	env *env.Env

	client *mongo.Client
	db     *mongo.Database

	block *mongo.Collection
	tx    *mongo.Collection
}

func NewDB(env *env.Env) (*DB, error) {
	d := &DB{
		env: env,
	}

	ctx := Context()
	var err error

	if d.client, err = mongo.Connect(ctx, options.Client().ApplyURI(env.DB.Uri)); err != nil {
		return nil, err
	} else if err = d.client.Ping(ctx, nil); err != nil {
		return nil, err
	} else {
		d.db = d.client.Database(env.DB.DB)

		d.block = d.db.Collection(env.DB.Block)
		d.tx = d.db.Collection(env.DB.Tx)

		return d, nil
	}

}
