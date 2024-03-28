// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package types

type EnumRpcCommitment string

// Rpc commitment context config
const (
	// RpcCommitmentFinalized the node will query the most recent block confirmed
	// by supermajority of the cluster as having reached maximum lockout,
	// meaning the cluster has recognized this block as finalized
	RpcCommitmentFinalized = "finalized"

	// RpcCommitmentConfirmed the node will query the most recent block
	// that has been voted on by supermajority of the cluster.
	RpcCommitmentConfirmed = "confirmed"

	// RpcCommitmentProcessed the node will query its most recent block.
	// Deprecated: Note that the block may still be skipped by the cluster.
	RpcCommitmentProcessed = "processed"
)

type EnumTxDetailLevel string

// level of transaction detail to return
const (
	TxDetailLevelNone       = "none"
	TxDetailLevelFull       = "full"
	TxDetailLevelAccounts   = "accounts"
	TxDetailLevelSignatures = "signatures"
)
