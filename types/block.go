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

type CTx struct {
	Hash          common.Hash            `json:"hash"`
	BlockNumber   *big.Int               `json:"blockNumber"`
	Type          uint8                  `json:"type"`
	Index         uint64                 `json:"index"`
	Size          uint64                 `json:"size"`
	From          common.Address         `json:"from"`
	To            *common.Address        `json:"to"`
	Nonce         uint64                 `json:"nonce"`
	GasPrice      *big.Int               `json:"gasPrice"`
	UsedGas       uint64                 `json:"usedGas"`
	Amount        string                 `json:"amount"`
	Status        uint64                 `json:"status"`
	TxTime        uint64                 `json:"txTime"`
	TotalTxAmount uint64                 `json:"totalTxAmount"`
	Fee           *big.Int               `json:"fee"`
	Header        *ethTypes.Header       `json:"header"`
	Signer        *ethTypes.EIP155Signer `json:"signer"`
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

func MakeCustomTx(
	tx *ethTypes.Transaction,
	index,
	totalTxAmount uint64,
	blockHeader *ethTypes.Header,
	receipt *ethTypes.Receipt,
	signer *ethTypes.EIP155Signer,
) *CTx {
	t := &CTx{
		Hash:          tx.Hash(),
		BlockNumber:   blockHeader.Number,
		Type:          tx.Type(),
		Index:         index,
		Size:          tx.Size(),
		Nonce:         tx.Nonce(),
		GasPrice:      tx.GasPrice(),
		Amount:        tx.Value().String(),
		UsedGas:       receipt.GasUsed,
		Status:        receipt.Status,
		TxTime:        blockHeader.Time,
		TotalTxAmount: totalTxAmount,
		Header:        blockHeader,
		Signer:        signer,
		Fee:           new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipt.GasUsed))),
		To:            tx.To(),
	}

	t.From, _ = ethTypes.Sender(signer, tx)

	return t
}
