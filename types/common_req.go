// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package types

// RpcCommitmentCfg rpc config of commitment
type RpcCommitmentCfg struct {
	Commitment EnumRpcCommitment `json:"commitment,omitempty"`
}

// RpcCommitmentAndMinSlotCfg commitment & min slot
type RpcCommitmentAndMinSlotCfg struct {
	Commitment     EnumRpcCommitment `json:"commitment,omitempty"`
	MinContextSlot uint64            `json:"minContextSlot,omitempty"`
}

// RpcGetBlockContextCfg commitment & min slot
type RpcGetBlockContextCfg struct {
	Rewards                        bool              `json:"rewards,omitempty"`
	Commitment                     EnumRpcCommitment `json:"commitment,omitempty"`
	TransactionDetails             EnumTxDetailLevel `json:"transactionDetails,omitempty"`
	MaxSupportedTransactionVersion uint64            `json:"maxSupportedTransactionVersion,omitempty"`
}
