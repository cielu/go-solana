// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

import "math/big"

// AccountInfo is the raw on-chain account payload returned by getAccountInfo,
// getMultipleAccounts, getProgramAccounts and simulateTransaction.
type AccountInfo struct {
	// data associated with the account, either as encoded binary data or JSON format {<program>: <state>} - depending on encoding parameter
	Data SolData `json:"data"`
	// base-58 encoded Pubkey of the program this account has been assigned to
	Owner PublicKey `json:"owner"`
	// number of lamports assigned to this account, as an u64
	Lamports *big.Int `json:"lamports"`
	// the epoch at which this account will next owe rent, as u64
	RentEpoch *big.Int `json:"rentEpoch"`
	// boolean indicating if the account contains a program (and is strictly read-only)
	Executable bool `json:"executable"`
	// the data size of the account
	Space uint64 `json:"space,omitempty"`
	//  the data size of the account
	Size uint64 `json:"size,omitempty"`
}

// AccountInfoWithCtx is the context-wrapped form returned by getAccountInfo.
type AccountInfoWithCtx struct {
	Context     ContextSlot  `json:"context"`
	AccountInfo *AccountInfo `json:"value,omitempty"`
}

// AccountsInfoWithCtx is the context-wrapped list returned by
// getMultipleAccounts.
type AccountsInfoWithCtx struct {
	Context  ContextSlot    `json:"context"`
	Accounts []*AccountInfo `json:"value,omitempty"`
}

// BalanceWithCtx wraps the lamport balance returned by getBalance.
type BalanceWithCtx struct {
	Context ContextSlot `json:"context"`
	Balance *big.Int    `json:"value"`
}

// AccountWithLamport is an entry in the list returned by getLargestAccounts.
type AccountWithLamport struct {
	// base-58 encoded address of the account
	Address PublicKey `json:"address"`
	// number of lamports in the account, as a u64
	Lamports *big.Int `json:"lamports"`
}

// ProgramAccount is an entry in the list returned by getProgramAccounts.
type ProgramAccount struct {
	Account AccountInfo `json:"account"`
	PubKey  PublicKey   `json:"pubKey"`
}

// StakeActivation is the response shape of getStakeActivation.
type StakeActivation struct {
	// the stake account's activation state, either: active, inactive, activating, or deactivating
	State string `json:"state"`
	// stake active during the epoch
	Active uint64 `json:"active"`
	// stake inactive during the epoch
	Inactive uint64 `json:"inactive"`
}
