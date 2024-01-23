package chain

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/big"
	"scanner/env"
	"scanner/log"
	"scanner/repo"
	. "scanner/types"
	. "scanner/util"
	"sync/atomic"
	"time"
)

type Chain struct {
	env     *env.Env
	chainID *big.Int
	signer  types.EIP155Signer

	repo *repo.Repo
}

func NewChain(env *env.Env, repo *repo.Repo, startBlock uint64) {
	c := &Chain{
		env:  env,
		repo: repo,
	}

	var err error
	if c.chainID, err = c.getChainId(); err != nil {
		log.CritLog(err.Error())
	} else {
		c.signer = types.NewEIP155Signer(c.chainID)
	}

	c.scanBlock(startBlock)
}

func (c *Chain) scanBlock(startBlock uint64) {

	st := startBlock

	if st == 0 {
		lBlock, err := c.latestBlock()
		if err != nil {
			log.CritLog(err.Error())
		} else {
			st = lBlock
		}
	}

	log.InfoLog("start block to scan")

	for {
		time.Sleep(3 * time.Second)

		ltBlock, err := c.repo.Node.GetLatestBlock()

		if err != nil {
			log.ErrLog(err.Error())
			continue
		} else if ltBlock < st {
			log.InfoLog(ErrToString(BlockNumberInvalid))
		} else {
			log.InfoLog("read Block Success : " + ToString(ltBlock))

			go c.readBlock(st, ltBlock)

			atomic.StoreUint64(&st, ltBlock)
		}
	}

}

func (c *Chain) readBlock(startBlock uint64, endBlock uint64) {
	for i := startBlock; i <= endBlock; i++ {
		// 현재 블록부터, 끝 블록까지 읽는다.
		if blockToRead, err := c.getBlockByNumber(big.NewInt(int64(i))); err != nil {
			log.ErrLog(ErrToString(CanNotFindBlock) + "err" + err.Error() + "block" + string(i))
			continue
		} else {
			totalLen := blockToRead.Transactions().Len()

			if totalLen < 1 {
				log.InfoLog(ErrToString(BlockTxLengthZero))
				continue
			} else {
				log.InfoLog("Scan Block Success : " + ToString(blockToRead.Number()))

				go c.saveBlock(blockToRead)                              // Tx가 존재하는 블록을 저장
				go c.saveTx(blockToRead, totalLen, blockToRead.Header()) // 블럭에 있는 tx 리스트 들을 모두 저장

			}

		}

	}

}

func (c *Chain) saveTx(blockToRead *types.Block, totalLen int, blockHeader *types.Header) {
	client := c.repo.Node.GetClient()

	var writeModel []mongo.WriteModel

	for j := 0; j < totalLen; j++ {
		tx := blockToRead.Transactions()[j]

		if re, err := client.TransactionReceipt(Context(), tx.Hash()); err != nil {
			log.ErrLog(err.Error())
		} else {
			ct := MakeCustomTx(tx, uint64(j), uint64(totalLen), blockHeader, re, &c.signer)

			if v, err := ToJson(ct); err != nil {
				log.ErrLog(fmt.Sprintf("failed to json tx : %s", hexutil.Encode(ct.Hash[:])))
				continue
			} else {
				writeModel = append(
					writeModel,
					mongo.NewUpdateOneModel().SetUpsert(true).
						SetFilter(bson.M{"hash": hexutil.Encode(ct.Hash[:])}).
						SetUpdate(bson.M{"$set": v}),
				)
			}
		}
	}

	if len(writeModel) != 0 {
		if err := c.repo.DB.BulkSaveTx(writeModel); err != nil {
			log.ErrLog(err.Error())
		}
	}

}

func (c *Chain) saveBlock(b *types.Block) {
	if err := c.repo.DB.SaveBlock(MakeCustomBlockType(b, c.chainID.Int64())); err != nil {
		log.ErrLog(err.Error())
	}
}

func (c *Chain) latestBlock() (uint64, error) {
	return c.repo.Node.GetLatestBlock()
}

func (c *Chain) getBlockByNumber(number *big.Int) (*types.Block, error) {
	return c.repo.Node.GetBlockByNumber(number)
}

func (c *Chain) getChainId() (*big.Int, error) {
	return c.repo.Node.GetChainID()
}
