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

// SyncNative Given a wrapped / native token account (a token account containing SOL)
// updates its amount field based on the account's underlying `lamports`.
// This is useful if a non-wrapped SOL account uses `system_instruction::transfer`
// to move lamports to a wrapped token account, and needs to have its token
// `amount` field updated.
type SyncNative struct {

	// [0] = [WRITE] tokenAccount
	// ··········· The native token account to sync with its underlying lamports.
	AccountMeta []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

// NewSyncNativeInstructionBuilder creates a new `SyncNative` instruction builder.
func NewSyncNativeInstructionBuilder() *SyncNative {
	nd := &SyncNative{
		AccountMeta: make([]*base.AccountMeta, 1),
	}
	return nd
}

func (sync SyncNative) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, sync.AccountMeta...)
	return
}

// SetTokenAccount sets the "tokenAccount" account.
// The native token account to sync with its underlying lamports.
func (sNative *SyncNative) SetTokenAccount(tokenAccount common.Address) *SyncNative {
	sNative.AccountMeta[0] = base.Meta(tokenAccount).WRITE()
	return sNative
}

// GetTokenAccount gets the "tokenAccount" account.
// The native token account to sync with its underlying lamports.
func (sNative *SyncNative) GetTokenAccount() *base.AccountMeta {
	return sNative.AccountMeta[0]
}

func (sNative SyncNative) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   sNative,
		TypeID: encodbin.TypeIDFromUint8(Instruction_SyncNative),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (sNative SyncNative) ValidateAndBuild() (*Instruction, error) {
	if err := sNative.Validate(); err != nil {
		return nil, err
	}
	return sNative.Build(), nil
}

func (sNative *SyncNative) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if sNative.AccountMeta[0] == nil {
			return errors.New("accounts.TokenAccount is not set")
		}
	}
	return nil
}

func (sNative SyncNative) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	return nil
}

// NewSyncNativeInstruction declares a new SyncNative instruction with the provided parameters and accounts.
func NewSyncNativeInstruction(
	// Accounts:
	tokenAccount common.Address) *SyncNative {
	return NewSyncNativeInstructionBuilder().
		SetTokenAccount(tokenAccount)
}
