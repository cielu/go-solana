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

// Like InitializeAccount2, but does not require the Rent sysvar to be provided.
type InitializeAccount3 struct {
	// The new account's owner/multisignature.
	Owner *common.Address

	// [0] = [WRITE] account
	// ··········· The account to initialize.
	//
	// [1] = [] mint
	// ··········· The mint this account will be associated with.
	AccountMeta []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

// NewInitializeAccount3InstructionBuilder creates a new `InitializeAccount3` instruction builder.
func NewInitializeAccount3InstructionBuilder() *InitializeAccount3 {
	nd := &InitializeAccount3{
		AccountMeta: make([]*base.AccountMeta, 2),
	}
	return nd
}

// SetOwner sets the "owner" parameter.
// The new account's owner/multisignature.
func (inst *InitializeAccount3) SetOwner(owner common.Address) *InitializeAccount3 {
	inst.Owner = &owner
	return inst
}

// SetAccount sets the "account" account.
// The account to initialize.
func (inst *InitializeAccount3) SetAccount(account common.Address) *InitializeAccount3 {
	inst.AccountMeta[0] = base.Meta(account).WRITE()
	return inst
}

// GetAccount gets the "account" account.
// The account to initialize.
func (inst *InitializeAccount3) GetAccount() *base.AccountMeta {
	return inst.AccountMeta[0]
}

// SetMintAccount sets the "mint" account.
// The mint this account will be associated with.
func (inst *InitializeAccount3) SetMintAccount(mint common.Address) *InitializeAccount3 {
	inst.AccountMeta[1] = base.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
// The mint this account will be associated with.
func (inst *InitializeAccount3) GetMintAccount() *base.AccountMeta {
	return inst.AccountMeta[1]
}

func (inst InitializeAccount3) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   inst,
		TypeID: encodbin.TypeIDFromUint8(Instruction_InitializeAccount3),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst InitializeAccount3) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *InitializeAccount3) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Owner == nil {
			return errors.New("Owner parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMeta[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if inst.AccountMeta[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
	}
	return nil
}

func (obj InitializeAccount3) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Owner` param:
	err = encoder.Encode(obj.Owner)
	if err != nil {
		return err
	}
	return nil
}

// NewInitializeAccount3Instruction declares a new InitializeAccount3 instruction with the provided parameters and accounts.
func NewInitializeAccount3Instruction(
	// Parameters:
	owner common.Address,
	// Accounts:
	account common.Address,
	mint common.Address) *InitializeAccount3 {
	return NewInitializeAccount3InstructionBuilder().
		SetOwner(owner).
		SetAccount(account).
		SetMintAccount(mint)
}
