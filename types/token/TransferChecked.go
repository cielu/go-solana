package token

import (
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

type TransferChecked struct {
	// The amount of tokens to transfer.
	Amount *uint64

	// Expected number of base 10 digits to the right of the decimal place.
	Decimals *uint8

	// [0] = [WRITE] source
	// ··········· The source account.
	//
	// [1] = [] mint
	// ··········· The token mint.
	//
	// [2] = [WRITE] destination
	// ··········· The destination account.
	//
	// [3] = [] owner
	// ··········· The source account's owner/delegate.
	//
	// [4...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*base.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (tc *TransferChecked) SetAccounts(accounts []*base.AccountMeta) error {
	tc.Accounts, tc.Signers = base.AccountMetaSlice(accounts).SplitFrom(4)
	return nil
}

func (tc TransferChecked) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, tc.Accounts...)
	accounts = append(accounts, tc.Signers...)
	return
}

// NewTransferCheckedInstructionBuilder creates a new `TransferChecked` instruction builder.
func NewTransferCheckedInstructionBuilder() *TransferChecked {
	nd := &TransferChecked{
		Accounts: make(base.AccountMetaSlice, 4),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens to transfer.
func (tc *TransferChecked) SetAmount(amount uint64) *TransferChecked {
	tc.Amount = &amount
	return tc
}

// SetDecimals sets the "decimals" parameter.
// Expected number of base 10 digits to the right of the decimal place.
func (tc *TransferChecked) SetDecimals(decimals uint8) *TransferChecked {
	tc.Decimals = &decimals
	return tc
}

// SetSourceAccount sets the "source" account.
// The source account.
func (tc *TransferChecked) SetSourceAccount(source common.Address) *TransferChecked {
	tc.Accounts[0] = base.Meta(source).WRITE()
	return tc
}

// GetSourceAccount gets the "source" account.
// The source account.
func (tc *TransferChecked) GetSourceAccount() *base.AccountMeta {
	return tc.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (tc *TransferChecked) SetMintAccount(mint common.Address) *TransferChecked {
	tc.Accounts[1] = base.Meta(mint)
	return tc
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (tc *TransferChecked) GetMintAccount() *base.AccountMeta {
	return tc.Accounts[1]
}

// SetDestinationAccount sets the "destination" account.
// The destination account.
func (tc *TransferChecked) SetDestinationAccount(destination common.Address) *TransferChecked {
	tc.Accounts[2] = base.Meta(destination).WRITE()
	return tc
}

// GetDestinationAccount gets the "destination" account.
// The destination account.
func (tc *TransferChecked) GetDestinationAccount() *base.AccountMeta {
	return tc.Accounts[2]
}

// SetOwnerAccount sets the "owner" account.
// The source account's owner/delegate.
func (tc *TransferChecked) SetOwnerAccount(owner common.Address, multisigSigners ...common.Address) *TransferChecked {
	tc.Accounts[3] = base.Meta(owner)
	if len(multisigSigners) == 0 {
		tc.Accounts[3].SIGNER()
	}
	for _, signer := range multisigSigners {
		tc.Signers = append(tc.Signers, base.Meta(signer).SIGNER())
	}
	return tc
}

// GetOwnerAccount gets the "owner" account.
// The source account's owner/delegate.
func (tc *TransferChecked) GetOwnerAccount() *base.AccountMeta {
	return tc.Accounts[3]
}

func (tc TransferChecked) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   tc,
		TypeID: encodbin.TypeIDFromUint8(Instruction_TransferChecked),
	}}
}

func (tc TransferChecked) MarshalWithEncoder(encoder encodbin.Encoder) (err error) {
	// Serialize `Amount` param:

	err = encoder.Encode(tc.Amount)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(tc.Decimals)
	if err != nil {
		return err
	}
	return nil
}

// NewTransferCheckedInstruction declares a new TransferChecked instruction with the provided parameters and accounts.
func NewTransferCheckedInstruction(
	// Parameters:
	amount uint64,
	decimals uint8,
	// Accounts:
	source common.Address,
	mint common.Address,
	destination common.Address,
	owner common.Address,
	multisigSigners []common.Address,
) *TransferChecked {
	return NewTransferCheckedInstructionBuilder().
		SetAmount(amount).
		SetDecimals(decimals).
		SetSourceAccount(source).
		SetMintAccount(mint).
		SetDestinationAccount(destination).
		SetOwnerAccount(owner, multisigSigners...)
}
