package native

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
	Signers     []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

func NewTransferInstructionBuilder() *Transfer {
	nd := &Transfer{
		AccountMeta: make([]*base.AccountMeta, 2),
	}
	return nd
}

func (trans Transfer) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, trans.AccountMeta...)
	accounts = append(accounts, trans.Signers...)
	return
}

func (trans *Transfer) SetLamports(lamports uint64) *Transfer {
	trans.Lamports = &lamports
	return trans
}

// Funding account
func (trans *Transfer) SetFundingAccount(fundingAccount common.Address) *Transfer {
	trans.AccountMeta[0] = base.Meta(fundingAccount).WRITE().SIGNER()
	return trans
}

// Recipient account
func (trans *Transfer) SetRecipientAccount(recipientAccount common.Address) *Transfer {
	trans.AccountMeta[1] = base.Meta(recipientAccount).WRITE()
	return trans
}

func (trans Transfer) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   trans,
		TypeID: encodbin.TypeIDFromUint32(Instruction_Transfer, binary.LittleEndian),
	}}
}

func (trans Transfer) MarshalWithEncoder(encoder encodbin.Encoder) error {
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*trans.Lamports)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewTransferInstruction(
	// Parameters:
	// Accounts:
	fundingAccount common.Address,
	recipientAccount common.Address,
	lamports uint64) *Transfer {
	return NewTransferInstructionBuilder().
		SetLamports(lamports).
		SetFundingAccount(fundingAccount).
		SetRecipientAccount(recipientAccount)
}
