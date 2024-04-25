package associated_token_account

import (
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/types"
	"github.com/near/borsh-go"
)

type Instruction borsh.Enum

const (
	InstructionCreate Instruction = iota
)

type CreateParam struct {
	Funder                 common.Address
	Owner                  common.Address
	Mint                   common.Address
	AssociatedTokenAccount common.Address
}

func Create(param CreateParam) types.Instruction {
	data, err := borsh.Serialize(struct {
		Instruction Instruction
	}{
		Instruction: InstructionCreate,
	})
	if err != nil {
		panic(err)
	}

	return types.Instruction{
		ProgramID: common.SPLAssociatedTokenAccountProgramID,
		Accounts: []types.AccountMeta{
			{PubKey: param.Funder, IsSigner: true, IsWritable: true},
			{PubKey: param.AssociatedTokenAccount, IsSigner: false, IsWritable: true},
			{PubKey: param.Owner, IsSigner: false, IsWritable: false},
			{PubKey: param.Mint, IsSigner: false, IsWritable: false},
			{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: false},
			{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
		},
		Data: data,
	}
}
