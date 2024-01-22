package db

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"scanner/env"
	"scanner/log"
	"scanner/types"
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

func (d *DB) SaveBlock(b *types.CBlock) error {
	filter := bson.M{"blockNumber": b.Number}

	opt := options.Update().SetUpsert(true)

	if j, err := toJson(b); err != nil {
		return err
	} else if result, err := d.block.UpdateOne(Context(), filter, bson.M{"$set": j}, opt); err != nil {
		return err
	} else {
		log.InfoLog("success to upsert block, modified : " + ToString(result.ModifiedCount) + " upserted : " + ToString(result.UpsertedCount))
		return nil
	}

}

func toJson(t interface{}) (interface{}, error) {
	var v interface{}
	if bytes, err := json.Marshal(t); err != nil {
		return nil, err
	} else if err := json.Unmarshal(bytes, &v); err != nil {
		return nil, err
	} else {
		jsonMap := v.(map[string]interface{})
		for key, value := range jsonMap {
			if reflect.TypeOf(value) == reflect.TypeOf(float64(0)) {
				jsonMap[key] = int64(value.(float64))
			}
		}

		return jsonMap, nil
	}
}
