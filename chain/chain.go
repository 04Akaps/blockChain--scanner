package chain

import (
	"github.com/ethereum/go-ethereum/core/types"
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

				go c.saveBlock(blockToRead) // Tx가 존재하는 블록을 저장

				for j := 0; j < totalLen; j++ {
					tx := blockToRead.Transactions()[j]
					go c.saveTx(tx, j, totalLen, blockToRead.Header())
				}
			}

		}

	}

}

func (c *Chain) saveTx(tx *types.Transaction, index, totalLen int, blockHeader *types.Header) {
	client := c.repo.Node.GetClient()

	if re, err := client.TransactionReceipt(Context(), tx.Hash()); err != nil {
		log.ErrLog(err.Error())
	} else {
		ct := MakeCustomTx(tx, uint64(index), uint64(totalLen), blockHeader, re, &c.signer)
		if err = c.repo.DB.SaveTx(ct); err != nil {
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
