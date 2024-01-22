package types

import (
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type CBlock struct {
	Hash       common.Hash `json:"hash"`
	Time       uint64      `json:"time"`
	Number     uint64      `json:"number"`
	ParentHash common.Hash `json:"parentHash"`
	Size       uint64      `json:"blockSize"`
	UsedGas    uint64      `json:"usedGas"`
	LimitGas   uint64      `json:"limitGas"`
	//Miner      common.Address `json:"miner"` // TODO add Miner
	TotalTxs int         `json:"totalTxs"`
	BaseFee  *big.Int    `json:"baseFee"`
	Root     common.Hash `json:"root"`

	BurnFee *big.Int `json:"burnFee"`
	ChainID int64    `json:"chainID"`
	Diff    *big.Int `json:"diff"`
}

func MakeCustomBlockType(b *ethTypes.Block, chainID int64) *CBlock {
	newBlock := &CBlock{
		Hash:       b.Hash(),
		Number:     b.NumberU64(),
		Time:       b.Time(),
		ParentHash: b.ParentHash(),
		UsedGas:    b.GasUsed(),
		Size:       b.Size(),
		TotalTxs:   b.Transactions().Len(),
		ChainID:    chainID,
		BaseFee:    b.BaseFee(),
		Root:       b.Root(),
		LimitGas:   b.GasLimit(),
		Diff:       b.Difficulty(),
	}

	newBlock.BurnFee = big.NewInt(1).Mul(b.BaseFee(), big.NewInt(int64(b.GasUsed())))

	return newBlock
}
