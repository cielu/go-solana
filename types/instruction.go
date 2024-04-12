package types

import "github.com/cielu/go-solana/common"

type Instruction struct {
	ProgramID common.Address
	Accounts  []AccountMeta
	Data      []byte
}

type AccountMeta struct {
	PubKey     common.Address
	IsSigner   bool
	IsWritable bool
}
