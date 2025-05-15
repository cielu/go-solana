package types

import (
	"github.com/cielu/go-solana/common"
)

type RawTransaction struct {
	instructions []Instruction
	blockHash    common.Hash
	payer        common.Address
	signers      []string
}

func NewRawTransaction(blockHash common.Hash, payer string, inst []Instruction, signers []string) *RawTransaction {
	return &RawTransaction{
		instructions: inst,
		blockHash:    blockHash,
		payer:        common.StrToAddress(payer),
		signers:      signers,
	}
}
