// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

type EnumRpcCommitment string

// Rpc commitment context config
const (
	// RpcCommitmentFinalized the node will query the most recent block confirmed
	// by supermajority of the cluster as having reached maximum lockout,
	// meaning the cluster has recognized this block as finalized
	RpcCommitmentFinalized EnumRpcCommitment = "finalized"

	// RpcCommitmentConfirmed the node will query the most recent block
	// that has been voted on by supermajority of the cluster.
	RpcCommitmentConfirmed EnumRpcCommitment = "confirmed"

	// RpcCommitmentProcessed the node will query its most recent block.
	// Note that the block may still be skipped by the cluster.
	RpcCommitmentProcessed EnumRpcCommitment = "processed"
)

type EnumEncoding string

// base58 base64 base64+zstd jsonParsed
const (
	EncodingBase58     EnumEncoding = "base58"
	EncodingBase64     EnumEncoding = "base64"
	EncodingBase64Zstd EnumEncoding = "base64+zstd"
	EncodingJson       EnumEncoding = "json"
	EncodingJsonParsed EnumEncoding = "jsonParsed"
)

type EnumTxDetailLevel string

// level of transaction detail to return
const (
	TxDetailLevelNone       EnumTxDetailLevel = "none"
	TxDetailLevelFull       EnumTxDetailLevel = "full"
	TxDetailLevelAccounts   EnumTxDetailLevel = "accounts"
	TxDetailLevelSignatures EnumTxDetailLevel = "signatures"
)

type EnumCirculateFilter string

// filter results by account type
const (
	FilterCirculating    EnumCirculateFilter = "circulating"
	FilterNonCirculating EnumCirculateFilter = "nonCirculating"
)
