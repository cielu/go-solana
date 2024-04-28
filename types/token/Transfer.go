package token

import (
	"encoding/binary"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

type Transfer struct {
	// Number of lamports to transfer to the new account
	Lamports *uint64

	// [0] = [WRITE, SIGNER] FundingAccount
	// ··········· Funding account
	//
	// [1] = [WRITE] RecipientAccount
	// ··········· Recipient account
	AccountMeta []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

func NewTransferInstructionBuilder() *Transfer {
	nd := &Transfer{
		AccountMeta: make([]*base.AccountMeta, 2),
	}
	return nd
}

func (inst *Transfer) SetLamports(lamports uint64) *Transfer {
	inst.Lamports = &lamports
	return inst
}

// Funding account
func (inst *Transfer) SetFundingAccount(fundingAccount common.Address) *Transfer {
	inst.AccountMeta[0] = base.Meta(fundingAccount).WRITE().SIGNER()
	return inst
}

// Recipient account
func (inst *Transfer) SetRecipientAccount(recipientAccount common.Address) *Transfer {
	inst.AccountMeta[1] = base.Meta(recipientAccount).WRITE()
	return inst
}

func (inst Transfer) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   inst,
		TypeID: encodbin.TypeIDFromUint32(Instruction_Transfer, binary.LittleEndian),
	}}
}

func (inst Transfer) MarshalWithEncoder(encoder encodbin.Encoder) error {
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*inst.Lamports)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewTransferInstruction(
	// Parameters:
	lamports uint64,
	// Accounts:
	fundingAccount common.Address,
	recipientAccount common.Address) *Transfer {
	return NewTransferInstructionBuilder().
		SetLamports(lamports).
		SetFundingAccount(fundingAccount).
		SetRecipientAccount(recipientAccount)
}
