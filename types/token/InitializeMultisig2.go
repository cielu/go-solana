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

// InitializeMultisig2 Like InitializeMultisig, but does not require the Rent sysvar to be provided.
type InitializeMultisig2 struct {
	// The number of signers (M) required to validate this multisignature account.
	M *uint8

	// [0] = [WRITE] account
	// ··········· The multisignature account to initialize.
	//
	// [1] = [SIGNER] signers
	// ··········· The signer accounts, must equal to N where 1 <= N <= 11.
	Accounts []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (initMs *InitializeMultisig2) SetAccounts(accounts []*solana.AccountMeta) error {
	initMs.Accounts, initMs.Signers = core.SliceSplitFrom(accounts, 1)
	return nil
}

func (initMs InitializeMultisig2) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, initMs.Accounts...)
	accounts = append(accounts, initMs.Signers...)
	return
}

// NewInitializeMultisig2InstructionBuilder creates a new `InitializeMultisig2` instruction builder.
func NewInitializeMultisig2InstructionBuilder() *InitializeMultisig2 {
	nd := &InitializeMultisig2{
		Accounts: make([]*solana.AccountMeta, 1),
		Signers:  make([]*solana.AccountMeta, 0),
	}
	return nd
}

// SetM sets the "m" parameter.
// The number of signers (M) required to validate this multisignature account.
func (initMs *InitializeMultisig2) SetM(m uint8) *InitializeMultisig2 {
	initMs.M = &m
	return initMs
}

// SetAccount sets the "account" account.
// The multisignature account to initialize.
func (initMs *InitializeMultisig2) SetAccount(account solana.PublicKey) *InitializeMultisig2 {
	initMs.Accounts[0] = solana.Meta(account).WRITE()
	return initMs
}

// GetAccount gets the "account" account.
// The multisignature account to initialize.
func (initMs *InitializeMultisig2) GetAccount() *solana.AccountMeta {
	return initMs.Accounts[0]
}

// AddSigners adds the "signers" accounts.
// The signer accounts, must equal to N where 1 <= N <= 11.
func (initMs *InitializeMultisig2) AddSigners(signers ...solana.PublicKey) *InitializeMultisig2 {
	for _, signer := range signers {
		initMs.Signers = append(initMs.Signers, solana.Meta(signer).SIGNER())
	}
	return initMs
}

func (initMs InitializeMultisig2) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   initMs,
		TypeID: encodbin.TypeIDFromUint8(Instruction_InitializeMultisig2),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (initMs InitializeMultisig2) ValidateAndBuild() (*Instruction, error) {
	if err := initMs.Validate(); err != nil {
		return nil, err
	}
	return initMs.Build(), nil
}

func (initMs *InitializeMultisig2) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if initMs.M == nil {
			return errors.New("M parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if initMs.Accounts[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if len(initMs.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(initMs.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(initMs.Signers))
		}
	}
	return nil
}

func (initMs InitializeMultisig2) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `M` param:
	err = encoder.Encode(initMs.M)
	if err != nil {
		return err
	}
	return nil
}

// NewInitializeMultisig2Instruction declares a new InitializeMultisig2 instruction with the provided parameters and accounts.
func NewInitializeMultisig2Instruction(
	// Parameters:
	m uint8,
	// Accounts:
	account solana.PublicKey,
	signers []solana.PublicKey,
) *InitializeMultisig2 {
	return NewInitializeMultisig2InstructionBuilder().
		SetM(m).
		SetAccount(account).
		AddSigners(signers...)
}
