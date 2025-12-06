package token

import (
	"github.com/cielu/go-solana"
	"github.com/cielu/go-solana/library"
	"github.com/cielu/go-solana/pkg/encodbin"
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
	Accounts []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (tc *TransferChecked) SetAccounts(accounts []*solana.AccountMeta) error {
	tc.Accounts, tc.Signers = library.SliceSplitFrom(accounts, 4)
	return nil
}

func (tc TransferChecked) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, tc.Accounts...)
	accounts = append(accounts, tc.Signers...)
	return
}

// NewTransferCheckedInstructionBuilder creates a new `TransferChecked` instruction builder.
func NewTransferCheckedInstructionBuilder() *TransferChecked {
	nd := &TransferChecked{
		Accounts: make([]*solana.AccountMeta, 4),
		Signers:  make([]*solana.AccountMeta, 0),
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
func (tc *TransferChecked) SetSourceAccount(source solana.PublicKey) *TransferChecked {
	tc.Accounts[0] = solana.Meta(source).WRITE()
	return tc
}

// GetSourceAccount gets the "source" account.
// The source account.
func (tc *TransferChecked) GetSourceAccount() *solana.AccountMeta {
	return tc.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (tc *TransferChecked) SetMintAccount(mint solana.PublicKey) *TransferChecked {
	tc.Accounts[1] = solana.Meta(mint)
	return tc
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (tc *TransferChecked) GetMintAccount() *solana.AccountMeta {
	return tc.Accounts[1]
}

// SetDestinationAccount sets the "destination" account.
// The destination account.
func (tc *TransferChecked) SetDestinationAccount(destination solana.PublicKey) *TransferChecked {
	tc.Accounts[2] = solana.Meta(destination).WRITE()
	return tc
}

// GetDestinationAccount gets the "destination" account.
// The destination account.
func (tc *TransferChecked) GetDestinationAccount() *solana.AccountMeta {
	return tc.Accounts[2]
}

// SetOwnerAccount sets the "owner" account.
// The source account's owner/delegate.
func (tc *TransferChecked) SetOwnerAccount(owner solana.PublicKey, multisigSigners ...solana.PublicKey) *TransferChecked {
	tc.Accounts[3] = solana.Meta(owner)
	if len(multisigSigners) == 0 {
		tc.Accounts[3].SIGNER()
	}
	for _, signer := range multisigSigners {
		tc.Signers = append(tc.Signers, solana.Meta(signer).SIGNER())
	}
	return tc
}

// GetOwnerAccount gets the "owner" account.
// The source account's owner/delegate.
func (tc *TransferChecked) GetOwnerAccount() *solana.AccountMeta {
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
	source solana.PublicKey,
	mint solana.PublicKey,
	destination solana.PublicKey,
	owner solana.PublicKey,
	multisigSigners []solana.PublicKey,
) *TransferChecked {
	return NewTransferCheckedInstructionBuilder().
		SetAmount(amount).
		SetDecimals(decimals).
		SetSourceAccount(source).
		SetMintAccount(mint).
		SetDestinationAccount(destination).
		SetOwnerAccount(owner, multisigSigners...)
}
