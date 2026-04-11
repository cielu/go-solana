// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

import (
	"encoding/json"
	"math/big"
)

// ContextSlot is the `context` object returned alongside most `*WithCtx`
// response bodies. It reports the slot the response was derived from.
type ContextSlot struct {
	Slot       uint64 `json:"slot"`
	ApiVersion string `json:"apiVersion,omitempty"`
}

// LastBlock describes a blockhash together with the last block height at
// which it remains valid. Returned by getLatestBlockhash and used inside
// SimulationInfo.ReplacementBlockhash.
type LastBlock struct {
	// a Hash as base-58 encoded string
	Blockhash Hash `json:"blockhash"`
	//  last block height at which the blockhash will be valid
	LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`
}

// LoadedAddresses is the set of accounts loaded from on-chain address lookup
// tables for a versioned transaction.
type LoadedAddresses struct {
	ReadOnly []PublicKey `json:"readonly,omitempty"`
	Writable []PublicKey `json:"writable,omitempty"`
}

// TxStatus is the legacy JSON shape used on some RPCs to report the outcome
// of a transaction. Use TransactionMeta.Err for new code.
type TxStatus struct {
	Ok  interface{}     `json:"Ok"`
	Err json.RawMessage `json:"Err"`
}

// U64ValueWithCtx is a context-wrapped optional uint64 used by several RPC
// methods that return a single scalar (getFeeForMessage,
// getStakeMinimumDelegation, etc).
type U64ValueWithCtx struct {
	Context ContextSlot `json:"context"`
	Value   *uint64     `json:"value,omitempty"`
}

// BlockReward is an entry in the `rewards` array of both block responses
// and transaction metadata.
type BlockReward struct {
	Commission  *uint8    `json:"commission"`
	Lamports    *big.Int  `json:"lamports"`
	PostBalance uint64    `json:"postBalance"`
	RewardType  string    `json:"rewardType"`
	Pubkey      PublicKey `json:"pubkey"`
}
