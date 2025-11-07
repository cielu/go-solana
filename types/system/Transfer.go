package system

import (
	"encoding/binary"

	"github.com/cielu/go-solana"
	"github.com/cielu/go-solana/pkg/encodbin"
)

type Transfer struct {
	// Number of lamports to transfer to the new account
	Lamports *uint64

	// [0] = [WRITE, SIGNER] FundingAccount
	// ··········· Funding account
	//
	// [1] = [WRITE] RecipientAccount
	// ··········· Recipient account
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers                 []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
}

func NewTransferInstructionBuilder() *Transfer {
	nd := &Transfer{
		AccountMetaSlice: make([]*solana.AccountMeta, 2),
	}
	return nd
}

func (trans Transfer) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, trans.AccountMetaSlice...)
	accounts = append(accounts, trans.Signers...)
	return
}

func (trans *Transfer) SetLamports(lamports uint64) *Transfer {
	trans.Lamports = &lamports
	return trans
}

// Funding account
func (trans *Transfer) SetFundingAccount(fundingAccount solana.PublicKey) *Transfer {
	trans.AccountMetaSlice[0] = solana.Meta(fundingAccount).WRITE().SIGNER()
	return trans
}

// Recipient account
func (trans *Transfer) SetRecipientAccount(recipientAccount solana.PublicKey) *Transfer {
	trans.AccountMetaSlice[1] = solana.Meta(recipientAccount).WRITE()
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
	fundingAccount solana.PublicKey,
	recipientAccount solana.PublicKey,
	lamports uint64) *Transfer {
	return NewTransferInstructionBuilder().
		SetLamports(lamports).
		SetFundingAccount(fundingAccount).
		SetRecipientAccount(recipientAccount)
}
