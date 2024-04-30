// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package token

import (
	"errors"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

// InitializeAccount2 Like InitializeAccount, but the owner pubkey is passed via instruction data
// rather than the accounts list. This variant may be preferable when using
// Cross Program Invocation from an instruction that does not need the owner's
// `AccountInfo` otherwise.
type InitializeAccount2 struct {
	// The new account's owner/multisignature.
	Owner *common.Address

	// [0] = [WRITE] account
	// ··········· The account to initialize.
	//
	// [1] = [] mint
	// ··········· The mint this account will be associated with.
	//
	// [2] = [] $(SysVarRentPubkey)
	// ··········· Rent sysvar.
	AccountMeta []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

// NewInitializeAccount2InstructionBuilder creates a new `InitializeAccount2` instruction builder.
func NewInitializeAccount2InstructionBuilder() *InitializeAccount2 {
	nd := &InitializeAccount2{
		AccountMeta: make([]*base.AccountMeta, 3),
	}
	nd.AccountMeta[2] = base.Meta(base.SysVarRentPubkey)
	return nd
}

// SetOwner sets the "owner" parameter.
// The new account's owner/multisignature.
func (initAcc2 *InitializeAccount2) SetOwner(owner common.Address) *InitializeAccount2 {
	initAcc2.Owner = &owner
	return initAcc2
}

// SetAccount sets the "account" account.
// The account to initialize.
func (initAcc2 *InitializeAccount2) SetAccount(account common.Address) *InitializeAccount2 {
	initAcc2.AccountMeta[0] = base.Meta(account).WRITE()
	return initAcc2
}

// GetAccount gets the "account" account.
// The account to initialize.
func (initAcc2 *InitializeAccount2) GetAccount() *base.AccountMeta {
	return initAcc2.AccountMeta[0]
}

// SetMintAccount sets the "mint" account.
// The mint this account will be associated with.
func (initAcc2 *InitializeAccount2) SetMintAccount(mint common.Address) *InitializeAccount2 {
	initAcc2.AccountMeta[1] = base.Meta(mint)
	return initAcc2
}

// GetMintAccount gets the "mint" account.
// The mint this account will be associated with.
func (initAcc2 *InitializeAccount2) GetMintAccount() *base.AccountMeta {
	return initAcc2.AccountMeta[1]
}

// SetSysVarRentPubkeyAccount sets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (initAcc2 *InitializeAccount2) SetSysVarRentPubkeyAccount(SysVarRentPubkey common.Address) *InitializeAccount2 {
	initAcc2.AccountMeta[2] = base.Meta(SysVarRentPubkey)
	return initAcc2
}

// GetSysVarRentPubkeyAccount gets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (initAcc2 *InitializeAccount2) GetSysVarRentPubkeyAccount() *base.AccountMeta {
	return initAcc2.AccountMeta[2]
}

func (initAcc2 InitializeAccount2) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   initAcc2,
		TypeID: encodbin.TypeIDFromUint8(Instruction_InitializeAccount2),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (initAcc2 InitializeAccount2) ValidateAndBuild() (*Instruction, error) {
	if err := initAcc2.Validate(); err != nil {
		return nil, err
	}
	return initAcc2.Build(), nil
}

func (initAcc2 *InitializeAccount2) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if initAcc2.Owner == nil {
			return errors.New("Owner parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if initAcc2.AccountMeta[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if initAcc2.AccountMeta[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if initAcc2.AccountMeta[2] == nil {
			return errors.New("accounts.SysVarRentPubkey is not set")
		}
	}
	return nil
}

func (initAcc2 InitializeAccount2) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Owner` param:
	err = encoder.Encode(initAcc2.Owner)
	if err != nil {
		return err
	}
	return nil
}

// NewInitializeAccount2Instruction declares a new InitializeAccount2 instruction with the provided parameters and accounts.
func NewInitializeAccount2Instruction(
	// Parameters:
	owner common.Address,
	// Accounts:
	account common.Address,
	mint common.Address) *InitializeAccount2 {
	return NewInitializeAccount2InstructionBuilder().
		SetOwner(owner).
		SetAccount(account).
		SetMintAccount(mint).
		SetSysVarRentPubkeyAccount(base.SysVarRentPubkey)
}
