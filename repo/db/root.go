package db

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

//0x9ff6712d37633e5b00a8cb9a86154db5e846602faafc0ea00bdab7ad0d6d0c84

func (d *DB) BulkSaveTx(model []mongo.WriteModel) error {
	if result, err := d.tx.BulkWrite(Context(), model); err != nil {
		return err
	} else {
		log.InfoLog("success to upsert block, modified : " + ToString(result.ModifiedCount) + " upserted : " + ToString(result.UpsertedCount))
		return nil
	}
}

func (d *DB) SaveTx(tx *types.CTx) error {
	filter := bson.M{"hash": hexutil.Encode(tx.Hash[:])}

	opt := options.Update().SetUpsert(true)

	if j, err := ToJson(tx); err != nil {
		return err
	} else if result, err := d.tx.UpdateOne(Context(), filter, bson.M{"$set": j}, opt); err != nil {
		return err
	} else {
		log.InfoLog("success to upsert block, modified : " + ToString(result.ModifiedCount) + " upserted : " + ToString(result.UpsertedCount))
		return nil
	}
}

func (d *DB) SaveBlock(b *types.CBlock) error {
	filter := bson.M{"blockNumber": b.Number}

	opt := options.Update().SetUpsert(true)

	if j, err := ToJson(b); err != nil {
		return err
	} else if result, err := d.block.UpdateOne(Context(), filter, bson.M{"$set": j}, opt); err != nil {
		return err
	} else {
		log.InfoLog("success to upsert block, modified : " + ToString(result.ModifiedCount) + " upserted : " + ToString(result.UpsertedCount))
		return nil
	}

}
