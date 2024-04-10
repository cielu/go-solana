// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package types

import "github.com/cielu/go-solana/common"

// RpcCommitmentCfg rpc config of commitment
type RpcCommitmentCfg struct {
	Commitment EnumRpcCommitment `json:"commitment,omitempty"`
}

// RpcCommitmentWithEncodingCfg rpc config of commitment
type RpcCommitmentWithEncodingCfg struct {
	Commitment EnumRpcCommitment `json:"commitment,omitempty"`
	Encoding   EnumEncoding      `json:"encoding,omitempty"`
}

// RpcCommitmentWithMinSlotCfg commitment & min slot
type RpcCommitmentWithMinSlotCfg struct {
	Commitment     EnumRpcCommitment `json:"commitment,omitempty"`
	MinContextSlot *uint64           `json:"minContextSlot,omitempty"`
}

// RpcGetBlockContextCfg commitment & min slot
type RpcGetBlockContextCfg struct {
	Commitment                     EnumRpcCommitment `json:"commitment,omitempty"`
	Encoding                       EnumEncoding      `json:"encoding,omitempty"`
	TransactionDetails             EnumTxDetailLevel `json:"transactionDetails,omitempty"`
	MaxSupportedTransactionVersion *uint64           `json:"maxSupportedTransactionVersion,omitempty"`
	Rewards                        *bool             `json:"rewards,omitempty"`
}

type RpcCommitmentWithIdentity struct {
	Commitment EnumRpcCommitment `json:"commitment,omitempty"`
	Identity   *common.Address   `json:"identity,omitempty"`
}

// RpcGetBlockProduction getBlock production
type RpcGetBlockProduction struct {
	Commitment EnumRpcCommitment `json:"commitment,omitempty"`
	Identity   *common.Address   `json:"identity,omitempty"`
	Range      SlotRange         `json:"range,omitempty"`
}

// RpcCommitmentWithFilter commitment with filter
type RpcCommitmentWithFilter struct {
	Commitment EnumRpcCommitment   `json:"commitment,omitempty"`
	Filter     EnumCirculateFilter `json:"filter,omitempty"`
}

// RpcAccountInfoCfg Get multiple data
type RpcAccountInfoCfg struct {
	Encoding       EnumEncoding      `json:"encoding,omitempty"`
	Commitment     EnumRpcCommitment `json:"commitment,omitempty"`
	MinContextSlot *uint64           `json:"minContextSlot,omitempty"`
	DataSlice      *DataSlice        `json:"dataSlice,omitempty"`
}

type RpcCombinedCfg struct {
	WithContext    bool                     `json:"withContext,omitempty"`
	Encoding       EnumEncoding             `json:"encoding,omitempty"`
	Commitment     EnumRpcCommitment        `json:"commitment,omitempty"`
	MinContextSlot *uint64                  `json:"minContextSlot,omitempty"`
	DataSlice      *DataSlice               `json:"dataSlice,omitempty"`
	Filter         []map[string]interface{} `json:"filter,omitempty"`
}

type RpcSearchTxHistoryCfg struct {
	// if true - a Solana node will search its ledger cache for any signatures not found in the recent status cache
	SearchTransactionHistory bool `json:"searchTransactionHistory,omitempty"`
}

type RpcSignaturesForAddressCfg struct {
	Commitment     EnumRpcCommitment `json:"commitment,omitempty"`
	MinContextSlot *uint64           `json:"minContextSlot,omitempty"`
	// Limit: maximum transaction signatures to return (between 1 and 1,000).
	Limit *uint `json:"limit,omitempty"`
	// start searching backwards from this transaction signature.
	// If not provided the search starts from the top of the highest max confirmed block.
	Before string `json:"before,omitempty"`
	// search until this transaction signature, if found before limit reached
	Util string `json:"util,omitempty"`
}

type RpcSupplyCfg struct {
	Commitment EnumRpcCommitment `json:"commitment,omitempty"`
	// exclude non circulating accounts list from response
	ExcludeNonCirculatingAccountsList *bool `json:"excludeNonCirculatingAccountsList,omitempty"`
}

type RpcMintWithProgramID struct {
	Mint      *common.Address `json:"mint,omitempty"`
	ProgramId *common.Address `json:"programId,omitempty"`
}

// RpcGetTransactionCfg commitment & min slot
type RpcGetTransactionCfg struct {
	Commitment EnumRpcCommitment `json:"commitment,omitempty"`
	Encoding   EnumEncoding      `json:"encoding,omitempty"`

	// MaxSupportedTransactionVersion Set the max transaction version to return in responses.
	// If the requested transaction is a higher version, an error will be returned.
	// If this parameter is omitted, only legacy transactions will be returned, and any versioned transaction will prompt the error.
	MaxSupportedTransactionVersion *uint64 `json:"maxSupportedTransactionVersion,omitempty"`
}

type RpcVoteAccountCfg struct {
	// commitment
	Commitment EnumRpcCommitment `json:"commitment,omitempty"`
	// Only return results for this validator vote address (base-58 encoded)
	VotePubkey *common.Address `json:"votePubkey,omitempty"` //  optional
	// Do not filter out delinquent validators with no stake
	KeepUnstakedDelinquents *bool `json:"keepUnstakedDelinquents,omitempty"` //  optional
	// Specify the number of slots behind the tip that a validator must fall to be considered delinquent. NOTE: For the sake of consistency between ecosystem products, it is not recommended that this argument be specified.
	DelinquentSlotDistance *uint64 `json:"delinquentSlotDistance,omitempty"` // optional
}

// RpcSendTxCfg struct
type RpcSendTxCfg struct {
	// Encoding used for the transaction data.
	// Default: base58
	// Values: base58 (slow, DEPRECATED), or base64.
	Encoding EnumEncoding `json:"encoding,omitempty"`
	// Default: false
	// when true, skip the preflight transaction checks
	SkipPreflight *bool `json:"skipPreflight,omitempty"`
	// Default: finalized
	// Commitment level to use for preflight.
	PreflightCommitment string `json:"preflightCommitment,omitempty"`
	// Maximum number of times for the RPC node to retry sending the transaction to the leader.
	// If this parameter not provided, the RPC node will retry the transaction until it is finalized or until the blockhash expires.
	MaxRetries *uint64 `json:"maxRetries,omitempty"`
	// set the minimum slot at which to perform preflight transaction checks
	MinContextSlot *uint64 `json:"minContextSlot,omitempty"`
}

type MentionsAccountProgramCfg struct {
	MentionsAccountOrProgram string `json:"MentionsAccountOrProgram,omitempty"`
}

type MentionsCfg struct {
	Mentions []string `json:"mentions,omitempty"`
}
