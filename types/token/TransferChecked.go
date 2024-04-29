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
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *TransferChecked) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(4)
	return nil
}

func (slice TransferChecked) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
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
func (inst *TransferChecked) SetAmount(amount uint64) *TransferChecked {
	inst.Amount = &amount
	return inst
}

// SetDecimals sets the "decimals" parameter.
// Expected number of base 10 digits to the right of the decimal place.
func (inst *TransferChecked) SetDecimals(decimals uint8) *TransferChecked {
	inst.Decimals = &decimals
	return inst
}

// SetSourceAccount sets the "source" account.
// The source account.
func (inst *TransferChecked) SetSourceAccount(source common.Address) *TransferChecked {
	inst.Accounts[0] = base.Meta(source).WRITE()
	return inst
}

// GetSourceAccount gets the "source" account.
// The source account.
func (inst *TransferChecked) GetSourceAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (inst *TransferChecked) SetMintAccount(mint common.Address) *TransferChecked {
	inst.Accounts[1] = base.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (inst *TransferChecked) GetMintAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetDestinationAccount sets the "destination" account.
// The destination account.
func (inst *TransferChecked) SetDestinationAccount(destination common.Address) *TransferChecked {
	inst.Accounts[2] = base.Meta(destination).WRITE()
	return inst
}

// GetDestinationAccount gets the "destination" account.
// The destination account.
func (inst *TransferChecked) GetDestinationAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

// SetOwnerAccount sets the "owner" account.
// The source account's owner/delegate.
func (inst *TransferChecked) SetOwnerAccount(owner common.Address, multisigSigners ...common.Address) *TransferChecked {
	inst.Accounts[3] = base.Meta(owner)
	if len(multisigSigners) == 0 {
		inst.Accounts[3].SIGNER()
	}
	for _, signer := range multisigSigners {
		inst.Signers = append(inst.Signers, base.Meta(signer).SIGNER())
	}
	return inst
}

// GetOwnerAccount gets the "owner" account.
// The source account's owner/delegate.
func (inst *TransferChecked) GetOwnerAccount() *base.AccountMeta {
	return inst.Accounts[3]
}

func (inst TransferChecked) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   inst,
		TypeID: encodbin.TypeIDFromUint8(Instruction_Transfer),
	}}
}

func (obj TransferChecked) MarshalWithEncoder(encoder encodbin.Encoder) (err error) {
	// Serialize `Amount` param:

	err = encoder.Encode(obj.Amount)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(obj.Decimals)
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
