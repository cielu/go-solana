// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

import (
	"encoding/json"
	"math/big"
)

// InnerInstruction is an inner-program invocation recorded in a transaction's
// metadata.
type InnerInstruction struct {
	// Index of the transaction instruction from which the inner instruction(s) originated
	Index uint16 `json:"index"`

	// Ordered list of inner program instructions that were invoked during a single transaction instruction.
	Instructions []CompiledInstruction `json:"instructions"`
}

// TransactionMeta is the `meta` object returned with a confirmed transaction.
type TransactionMeta struct {
	// TODO if has zero ComputeUnitsConsumed
	ComputeUnitsConsumed *uint64 `json:"computeUnitsConsumed"`
	// Error if transaction failed, null if transaction succeeded.
	// https://github.com/solana-labs/solana/blob/master/sdk/src/transaction.rs#L24
	Err json.RawMessage `json:"err"`

	// Fee this transaction was charged
	Fee uint64 `json:"fee"`

	// Array of *big.Int account balances from before the transaction was processed
	PreBalances []*big.Int `json:"preBalances"`

	// Array of *big.Int account balances after the transaction was processed
	PostBalances []*big.Int `json:"postBalances"`

	// List of inner instructions or omitted if inner instruction recording
	// was not yet enabled during this transaction
	InnerInstructions []InnerInstruction `json:"innerInstructions"`

	// List of token balances from before the transaction was processed
	// or omitted if token balance recording was not yet enabled during this transaction
	PreTokenBalances []TokenBalance `json:"preTokenBalances"`

	// List of token balances from after the transaction was processed
	// or omitted if token balance recording was not yet enabled during this transaction
	PostTokenBalances []TokenBalance `json:"postTokenBalances"`

	// Array of string log messages or omitted if log message
	// recording was not yet enabled during this transaction
	LogMessages []string `json:"logMessages"`

	// Transaction status.
	Status TxStatus `json:"status"`

	Rewards []BlockReward `json:"rewards"`

	LoadedAddresses LoadedAddresses `json:"loadedAddresses"`
}

// TransactionInfo is the response shape of getTransaction and the element type
// of BlockInfo.Transactions.
type TransactionInfo struct {
	// Transaction status metadata object
	Meta *TransactionMeta `json:"meta"`
	// The slot this transaction was processed in.
	Slot uint64 `json:"slot"`

	// Estimated production time, as Unix timestamp (seconds since the Unix epoch)
	// of when the transaction was processed.
	// Nil if not available.
	BlockTime *int64 `json:"blockTime" bin:"optional"`

	// Transaction
	Transaction *Transaction `json:"transaction"`
	// Version Of Transaction
	Version TxVersion `json:"version"`
}

// SimulateResult is the context-wrapped response shape of simulateTransaction.
type SimulateResult struct {
	Context ContextSlot    `json:"context"`
	Value   SimulationInfo `json:"value"`
}

// SimulationInfo is the value portion of a simulateTransaction response.
type SimulationInfo struct {
	Accounts []AccountInfo `json:"accounts"`
	// Error if transaction failed, null if transaction succeeded.
	// https://github.com/solana-labs/solana/blob/master/sdk/src/transaction.rs#L24
	Err json.RawMessage `json:"err"`
	// Fee this transaction was charged
	Fee uint64 `json:"fee"`

	// Array of *big.Int account balances from before the transaction was processed
	PreBalances []*big.Int `json:"preBalances"`

	// Array of *big.Int account balances after the transaction was processed
	PostBalances []*big.Int `json:"postBalances"`

	// List of inner instructions or omitted if inner instruction recording
	// was not yet enabled during this transaction
	InnerInstructions []InnerInstruction `json:"innerInstructions"`

	// List of token balances from before the transaction was processed
	// or omitted if token balance recording was not yet enabled during this transaction
	PreTokenBalances []TokenBalance `json:"preTokenBalances"`

	// List of token balances from after the transaction was processed
	// or omitted if token balance recording was not yet enabled during this transaction
	PostTokenBalances []TokenBalance `json:"postTokenBalances"`

	LoadedAccountsDataSize *uint32         `json:"loadedAccountsDataSize"`
	LoadedAddresses        LoadedAddresses `json:"loadedAddresses"`
	// Array of string log messages or omitted if log message
	// recording was not yet enabled during this transaction
	Logs []string `json:"logs"`

	ReplacementBlockhash *LastBlock      `json:"replacementBlockhash"`
	ReturnData           *SimulateReturn `json:"returnData"`
	UnitsConsumed        *uint64         `json:"unitsConsumed"`
}

// SimulateReturn is the `returnData` portion of a simulation response.
type SimulateReturn struct {
	ProgramID string  `json:"programId"`
	Data      SolData `json:"data"`
}

// SignatureStatus is an entry in the response of getSignatureStatuses.
type SignatureStatus struct {
	// The slot that contains the block with the transaction
	Slot uint64 `json:"slot"`
	// Error if transaction failed, null if transaction succeeded. See TransactionError definitions for more info.
	Err json.RawMessage `json:"err"`
	// estimated production time, as Unix timestamp (seconds since the Unix epoch) of when transaction was processed. null if not available.
	Confirmations *uint64 `json:"confirmations,omitempty"`
	// The transaction's cluster confirmation status; Either processed, confirmed, or finalized.
	ConfirmationStatus string `json:"confirmationStatus,omitempty"`
}

// SignatureStatusWithCtx is the context-wrapped response shape of
// getSignatureStatuses.
type SignatureStatusWithCtx struct {
	Context         ContextSlot       `json:"context"`
	SignatureStatus []SignatureStatus `json:"value,omitempty"`
}

// SignatureInfo is an entry in the response of getSignaturesForAddress.
type SignatureInfo struct {
	// transaction signature as base-58 encoded string
	Signature Signature `json:"signature,omitempty"`
	// The slot that contains the block with the transaction
	Slot uint64 `json:"slot"`
	// Error if transaction failed, null if transaction succeeded. See TransactionError definitions for more info.
	Err json.RawMessage `json:"err"`
	// Memo associated with the transaction, null if no memo is present
	Memo string `json:"memo,omitempty"`
	// estimated production time, as Unix timestamp (seconds since the Unix epoch) of when transaction was processed. null if not available.
	BlockTime int64 `json:"blockTime,omitempty"`
	// The transaction's cluster confirmation status; Either processed, confirmed, or finalized.
	ConfirmationStatus string `json:"confirmationStatus,omitempty"`
}
