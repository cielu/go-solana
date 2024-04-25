package token

import (
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/bincode"
	"github.com/cielu/go-solana/types"
)

type Instruction uint8

const (
	InstructionInitializeMint Instruction = iota
	InstructionInitializeAccount
	InstructionInitializeMultisig
	InstructionTransfer
	InstructionApprove
	InstructionRevoke
	InstructionSetAuthority
	InstructionMintTo
	InstructionBurn
	InstructionCloseAccount
	InstructionFreezeAccount
	InstructionThawAccount
	InstructionTransferChecked
	InstructionApproveChecked
	InstructionMintToChecked
	InstructionBurnChecked
	InstructionInitializeAccount2
	InstructionSyncNative
	InstructionInitializeAccount3
	InstructionInitializeMultisig2
	InstructionInitializeMint2
)

type TransferParam struct {
	From    common.Address
	To      common.Address
	Auth    common.Address
	Signers []common.Address
	Amount  uint64
}

type TransferCheckedParam struct {
	From     common.Address
	To       common.Address
	Mint     common.Address
	Auth     common.Address
	Signers  []common.Address
	Amount   uint64
	Decimals uint8
}

func TransferChecked(param TransferCheckedParam) types.Instruction {
	data, err := bincode.SerializeData(struct {
		Instruction Instruction
		Amount      uint64
		Decimals    uint8
	}{
		Instruction: InstructionTransferChecked,
		Amount:      param.Amount,
		Decimals:    param.Decimals,
	})
	if err != nil {
		panic(err)
	}

	accounts := make([]types.AccountMeta, 0, 4+len(param.Signers))
	accounts = append(accounts, types.AccountMeta{PubKey: param.From, IsSigner: false, IsWritable: true})
	accounts = append(accounts, types.AccountMeta{PubKey: param.Mint, IsSigner: false, IsWritable: false})
	accounts = append(accounts, types.AccountMeta{PubKey: param.To, IsSigner: false, IsWritable: true})
	accounts = append(accounts, types.AccountMeta{PubKey: param.Auth, IsSigner: len(param.Signers) == 0, IsWritable: false})
	for _, signerPubkey := range param.Signers {
		accounts = append(accounts, types.AccountMeta{PubKey: signerPubkey, IsSigner: true, IsWritable: false})
	}

	return types.Instruction{
		ProgramID: common.TokenProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

func Transfer(param TransferParam) types.Instruction {
	data, err := bincode.SerializeData(struct {
		Instruction Instruction
		Amount      uint64
	}{
		Instruction: InstructionTransfer,
		Amount:      param.Amount,
	})
	if err != nil {
		panic(err)
	}

	accounts := make([]types.AccountMeta, 0, 3+len(param.Signers))
	accounts = append(accounts, types.AccountMeta{PubKey: param.From, IsSigner: false, IsWritable: true})
	accounts = append(accounts, types.AccountMeta{PubKey: param.To, IsSigner: false, IsWritable: true})
	accounts = append(accounts, types.AccountMeta{PubKey: param.Auth, IsSigner: len(param.Signers) == 0, IsWritable: false})
	for _, signerPubkey := range param.Signers {
		accounts = append(accounts, types.AccountMeta{PubKey: signerPubkey, IsSigner: true, IsWritable: false})
	}
	return types.Instruction{
		ProgramID: common.TokenProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}
