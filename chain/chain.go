package chain

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"scanner/env"
	"scanner/log"
	"scanner/repo"
	. "scanner/types"
	"sync/atomic"
	"time"
)

type Chain struct {
	env *env.Env
	out chan struct{}

	repo *repo.Repo
}

func NewChain(env *env.Env, repo *repo.Repo, startBlock uint64) {
	c := &Chain{
		env:  env,
		out:  make(chan struct{}),
		repo: repo,
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

	go func() {
		for {
			time.Sleep(3 * time.Second)

			ltBlock, err := c.repo.Node.GetLatestBlock()

			if err != nil {
				log.ErrLog(err.Error())
			} else if ltBlock < st {
				log.InfoLog(ErrToString(BlockNumberInvalid))
			} else {
				go c.readBlock(st, ltBlock)

				atomic.StoreUint64(&st, ltBlock)
				c.out <- struct{}{}
			}
		}
	}()

}

func (c *Chain) readBlock(startBlock uint64, endBlock uint64) {
	for i := startBlock; i <= endBlock; i++ {
		// 현재 블록부터, 끝 블록까지 읽는다.
		if blockToRead, err := c.getBlockByNumber(big.NewInt(int64(i))); err != nil {
			log.ErrLog(ErrToString(CanNotFindBlock) + "err" + err.Error() + "block" + string(i))
		} else {

		}

	}

}

func (c *Chain) latestBlock() (uint64, error) {
	return c.repo.Node.GetLatestBlock()
}

func (c *Chain) getBlockByNumber(number *big.Int) (*types.Block, error) {
	return c.repo.Node.GetBlockByNumber(number)
}
