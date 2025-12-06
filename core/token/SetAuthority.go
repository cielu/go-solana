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
	"github.com/cielu/go-solana/library"
	"github.com/cielu/go-solana/pkg/encodbin"
)

// SetAuthority Sets a new authority of a mint or account.
type SetAuthority struct {
	// The type of authority to update.
	AuthorityType *AuthorityType

	// The new authority.
	NewAuthority *solana.PublicKey `bin:"optional"`

	// [0] = [WRITE] subject
	// ··········· The mint or account to change the authority of.
	//
	// [1] = [] authority
	// ··········· The current authority of the mint or account.
	//
	// [2...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (sAut *SetAuthority) SetAccounts(accounts []*solana.AccountMeta) error {
	sAut.Accounts, sAut.Signers = library.SliceSplitFrom(accounts, 2)
	return nil
}

func (sAut SetAuthority) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, sAut.Accounts...)
	accounts = append(accounts, sAut.Signers...)
	return
}

// NewSetAuthorityInstructionBuilder creates a new `SetAuthority` instruction builder.
func NewSetAuthorityInstructionBuilder() *SetAuthority {
	nd := &SetAuthority{
		Accounts: make([]*solana.AccountMeta, 2),
		Signers:  make([]*solana.AccountMeta, 0),
	}
	return nd
}

// SetAuthorityType sets the "authority_type" parameter.
// The type of authority to update.
func (sAut *SetAuthority) SetAuthorityType(authority_type AuthorityType) *SetAuthority {
	sAut.AuthorityType = &authority_type
	return sAut
}

// SetNewAuthority sets the "new_authority" parameter.
// The new authority.
func (sAut *SetAuthority) SetNewAuthority(new_authority solana.PublicKey) *SetAuthority {
	sAut.NewAuthority = &new_authority
	return sAut
}

// SetSubjectAccount sets the "subject" account.
// The mint or account to change the authority of.
func (sAut *SetAuthority) SetSubjectAccount(subject solana.PublicKey) *SetAuthority {
	sAut.Accounts[0] = solana.Meta(subject).WRITE()
	return sAut
}

// GetSubjectAccount gets the "subject" account.
// The mint or account to change the authority of.
func (sAut *SetAuthority) GetSubjectAccount() *solana.AccountMeta {
	return sAut.Accounts[0]
}

// SetAuthorityAccount sets the "authority" account.
// The current authority of the mint or account.
func (sAut *SetAuthority) SetAuthorityAccount(authority solana.PublicKey, multisigSigners ...solana.PublicKey) *SetAuthority {
	sAut.Accounts[1] = solana.Meta(authority)
	if len(multisigSigners) == 0 {
		sAut.Accounts[1].SIGNER()
	}
	for _, signer := range multisigSigners {
		sAut.Signers = append(sAut.Signers, solana.Meta(signer).SIGNER())
	}
	return sAut
}

// GetAuthorityAccount gets the "authority" account.
// The current authority of the mint or account.
func (sAut *SetAuthority) GetAuthorityAccount() *solana.AccountMeta {
	return sAut.Accounts[1]
}

func (sAut SetAuthority) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   sAut,
		TypeID: encodbin.TypeIDFromUint8(Instruction_SetAuthority),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (sAut SetAuthority) ValidateAndBuild() (*Instruction, error) {
	if err := sAut.Validate(); err != nil {
		return nil, err
	}
	return sAut.Build(), nil
}

func (sAut *SetAuthority) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if sAut.AuthorityType == nil {
			return errors.New("AuthorityType parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if sAut.Accounts[0] == nil {
			return errors.New("accounts.Subject is not set")
		}
		if sAut.Accounts[1] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if !sAut.Accounts[1].IsSigner && len(sAut.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(sAut.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(sAut.Signers))
		}
	}
	return nil
}

func (sAut SetAuthority) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `AuthorityType` param:
	err = encoder.Encode(sAut.AuthorityType)
	if err != nil {
		return err
	}
	// Serialize `NewAuthority` param (optional):
	{
		if sAut.NewAuthority == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(sAut.NewAuthority)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// NewSetAuthorityInstruction declares a new SetAuthority instruction with the provided parameters and accounts.
func NewSetAuthorityInstruction(
	// Parameters:
	authority_type AuthorityType,
	new_authority solana.PublicKey,
	// Accounts:
	subject solana.PublicKey,
	authority solana.PublicKey,
	multisigSigners []solana.PublicKey,
) *SetAuthority {
	return NewSetAuthorityInstructionBuilder().
		SetAuthorityType(authority_type).
		SetNewAuthority(new_authority).
		SetSubjectAccount(subject).
		SetAuthorityAccount(authority, multisigSigners...)
}
