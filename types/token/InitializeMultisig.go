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
	"github.com/cielu/go-solana/types/base"
)

// InitializeMultisig Initializes a multisignature account with N provided signers.
//
// Multisignature accounts can used in place of any single owner/delegate
// accounts in any token instruction that require an owner/delegate to be
// present.  The variant field represents the number of signers (M)
// required to validate this multisignature account.
//
// The `InitializeMultisig` instruction requires no signers and MUST be
// included within the same Transaction as the system program's
// `CreateAccount` instruction that creates the account being initialized.
// Otherwise another party can acquire ownership of the uninitialized
// account.
type InitializeMultisig struct {
	// The number of signers (M) required to validate this multisignature
	// account.
	M *uint8

	// [0] = [WRITE] account
	// ··········· The multisignature account to initialize.
	//
	// [1] = [] $(SysVarRentPubkey)
	// ··········· Rent sysvar.
	//
	// [2...] = [SIGNER] signers
	// ··········· ..2+N The signer accounts, must equal to N where 1 <= N <=11
	Accounts []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (initMs *InitializeMultisig) SetAccounts(accounts []*solana.AccountMeta) error {
	initMs.Accounts, initMs.Signers = core.SliceSplitFrom(accounts, 2)
	return nil
}

func (initMs InitializeMultisig) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, initMs.Accounts...)
	accounts = append(accounts, initMs.Signers...)
	return
}

// NewInitializeMultisigInstructionBuilder creates a new `InitializeMultisig` instruction builder.
func NewInitializeMultisigInstructionBuilder() *InitializeMultisig {
	nd := &InitializeMultisig{
		Accounts: make([]*solana.AccountMeta, 2),
		Signers:  make([]*solana.AccountMeta, 0),
	}
	nd.Accounts[1] = solana.Meta(base.SysVarRentPubkey)
	return nd
}

// SetM sets the "m" parameter.
// The number of signers (M) required to validate this multisignature
// account.
func (initMs *InitializeMultisig) SetM(m uint8) *InitializeMultisig {
	initMs.M = &m
	return initMs
}

// SetAccount sets the "account" account.
// The multisignature account to initialize.
func (initMs *InitializeMultisig) SetAccount(account solana.PublicKey) *InitializeMultisig {
	initMs.Accounts[0] = solana.Meta(account).WRITE()
	return initMs
}

// GetAccount gets the "account" account.
// The multisignature account to initialize.
func (initMs *InitializeMultisig) GetAccount() *solana.AccountMeta {
	return initMs.Accounts[0]
}

// SetSysVarRentPubkeyAccount sets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (initMs *InitializeMultisig) SetSysVarRentPubkeyAccount(SysVarRentPubkey solana.PublicKey) *InitializeMultisig {
	initMs.Accounts[1] = solana.Meta(SysVarRentPubkey)
	return initMs
}

// GetSysVarRentPubkeyAccount gets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (initMs *InitializeMultisig) GetSysVarRentPubkeyAccount() *solana.AccountMeta {
	return initMs.Accounts[1]
}

// AddSigners adds the "signers" accounts.
// ..2+N The signer accounts, must equal to N where 1 <= N <=11
func (initMs *InitializeMultisig) AddSigners(signers ...solana.PublicKey) *InitializeMultisig {
	for _, signer := range signers {
		initMs.Signers = append(initMs.Signers, solana.Meta(signer).SIGNER())
	}
	return initMs
}

func (initMs InitializeMultisig) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   initMs,
		TypeID: encodbin.TypeIDFromUint8(Instruction_InitializeMultisig),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (initMs InitializeMultisig) ValidateAndBuild() (*Instruction, error) {
	if err := initMs.Validate(); err != nil {
		return nil, err
	}
	return initMs.Build(), nil
}

func (initMs *InitializeMultisig) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if initMs.M == nil {
			return errors.New("M parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if initMs.Accounts[0] == nil {
			return fmt.Errorf("accounts.Account is not set")
		}
		if initMs.Accounts[1] == nil {
			return fmt.Errorf("accounts.SysVarRentPubkey is not set")
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

func (initMs InitializeMultisig) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `M` param:
	err = encoder.Encode(initMs.M)
	if err != nil {
		return err
	}
	return nil
}

// NewInitializeMultisigInstruction declares a new InitializeMultisig instruction with the provided parameters and accounts.
func NewInitializeMultisigInstruction(
	// Parameters:
	m uint8,
	// Accounts:
	account solana.PublicKey,
	SysVarRentPubkey solana.PublicKey,
	signers []solana.PublicKey,
) *InitializeMultisig {
	return NewInitializeMultisigInstructionBuilder().
		SetM(m).
		SetAccount(account).
		SetSysVarRentPubkeyAccount(SysVarRentPubkey).
		AddSigners(signers...)
}
