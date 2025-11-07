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
	"fmt"

	"github.com/cielu/go-solana"
	"github.com/cielu/go-solana/core"
	"github.com/cielu/go-solana/pkg/encodbin"
)

// Burn tokens by removing them from an account.  `Burn` does not support
// accounts associated with the native mint, use `CloseAccount` instead.
type Burn struct {
	// The amount of tokens to burn.
	Amount *uint64

	// [0] = [WRITE] source
	// ··········· The account to burn from.
	//
	// [1] = [WRITE] mint
	// ··········· The token mint.
	//
	// [2] = [] owner
	// ··········· The account's owner/delegate.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (br *Burn) SetAccounts(accounts []*solana.AccountMeta) error {
	br.Accounts, br.Signers = core.SliceSplitFrom(accounts, 3)
	return nil
}

func (br Burn) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, br.Accounts...)
	accounts = append(accounts, br.Signers...)
	return
}

// NewBurnInstructionBuilder creates a new `Burn` instruction builder.
func NewBurnInstructionBuilder() *Burn {
	nd := &Burn{
		Accounts: make([]*solana.AccountMeta, 3),
		Signers:  make([]*solana.AccountMeta, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens to burn.
func (br *Burn) SetAmount(amount uint64) *Burn {
	br.Amount = &amount
	return br
}

// SetSourceAccount sets the "source" account.
// The account to burn from.
func (br *Burn) SetSourceAccount(source solana.PublicKey) *Burn {
	br.Accounts[0] = solana.Meta(source).WRITE()
	return br
}

// GetSourceAccount gets the "source" account.
// The account to burn from.
func (br *Burn) GetSourceAccount() *solana.AccountMeta {
	return br.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (br *Burn) SetMintAccount(mint solana.PublicKey) *Burn {
	br.Accounts[1] = solana.Meta(mint).WRITE()
	return br
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (br *Burn) GetMintAccount() *solana.AccountMeta {
	return br.Accounts[1]
}

// SetOwnerAccount sets the "owner" account.
// The account's owner/delegate.
func (br *Burn) SetOwnerAccount(owner solana.PublicKey, multisigSigners ...solana.PublicKey) *Burn {
	br.Accounts[2] = solana.Meta(owner)
	if len(multisigSigners) == 0 {
		br.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		br.Signers = append(br.Signers, solana.Meta(signer).SIGNER())
	}
	return br
}

// GetOwnerAccount gets the "owner" account.
// The account's owner/delegate.
func (br *Burn) GetOwnerAccount() *solana.AccountMeta {
	return br.Accounts[2]
}

func (br Burn) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   br,
		TypeID: encodbin.TypeIDFromUint8(Instruction_Burn),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (br Burn) ValidateAndBuild() (*Instruction, error) {
	if err := br.Validate(); err != nil {
		return nil, err
	}
	return br.Build(), nil
}

func (br *Burn) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if br.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if br.Accounts[0] == nil {
			return errors.New("accounts.Source is not set")
		}
		if br.Accounts[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if br.Accounts[2] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if !br.Accounts[2].IsSigner && len(br.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(br.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(br.Signers))
		}
	}
	return nil
}

func (br Burn) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(br.Amount)
	if err != nil {
		return err
	}
	return nil
}

// NewBurnInstruction declares a new Burn instruction with the provided parameters and accounts.
func NewBurnInstruction(
	// Parameters:
	amount uint64,
	// Accounts:
	source solana.PublicKey,
	mint solana.PublicKey,
	owner solana.PublicKey,
	multisigSigners []solana.PublicKey,
) *Burn {
	return NewBurnInstructionBuilder().
		SetAmount(amount).
		SetSourceAccount(source).
		SetMintAccount(mint).
		SetOwnerAccount(owner, multisigSigners...)
}
