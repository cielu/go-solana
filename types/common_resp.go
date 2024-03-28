// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package types

import (
	"go-solana/common"
	"math/big"
)

type ContextSlot struct {
	Slot uint64 `json:"slot"`
}

type AccountInfo struct {
	Data       common.Base58Data `json:"data"`
	Owner      common.Address    `json:"owner"`
	Lamports   uint64            `json:"lamports"`
	RentEpoch  uint64            `json:"rentEpoch"`
	Executable bool              `json:"executable"`
}

type AccountInfoWithCtx struct {
	Context     ContextSlot `json:"context"`
	AccountInfo AccountInfo `json:"value"`
}

type BalanceWithCtx struct {
	Context ContextSlot `json:"context"`
	Balance *big.Int    `json:"value"`
}

type BlockReward struct {
	Commission  *uint8         `json:"commission"`
	Lamports    int64          `json:"lamports"`
	PostBalance uint64         `json:"postBalance"`
	RewardType  string         `json:"rewardType"`
	Pubkey      common.Address `json:"pubkey"`
}

type UiTokenAmount struct {
	// Raw amount of tokens as a string, ignoring decimals.
	Amount string `json:"amount"`

	// Number of decimals configured for token's mint.
	Decimals uint8 `json:"decimals"`

	// Token amount as a float, accounting for decimals.
	UiAmount float64 `json:"uiAmount"`

	// Token amount as a string, accounting for decimals.
	UiAmountString string `json:"uiAmountString"`
}

type TokenBalance struct {
	// Index of the account in which the token balance is provided for.
	AccountIndex uint16 `json:"accountIndex"`

	// Pubkey of the token's mint.
	Mint common.Address `json:"mint"`

	// Pubkey of token balance's owner.
	Owner common.Address `json:"owner"`

	// ProgramId
	ProgramId string `json:"programId"`

	UiTokenAmount UiTokenAmount `json:"uiTokenAmount"`
}

type CompiledInstruction struct {
	// StackHeight if empty
	StackHeight *uint16 `json:"stackHeight"`
	// Index into the message.accountKeys array indicating the program account that executes this instruction.
	// NOTE: it is actually an uint8, but using an uint16 because uint8 is treated as a byte everywhere,
	// and that can be an issue.
	ProgramIDIndex uint16 `json:"programIdIndex"`

	// List of ordered indices into the message.accountKeys array indicating which accounts to pass to the program.
	// NOTE: it is actually a []uint8, but using an uint16 because []uint8 is treated as a []byte everywhere,
	// and that can be an issue.
	Accounts []uint16 `json:"accounts"`

	// The program input data encoded in a base-58 string.
	Data common.Base58Data `json:"data"`
}

type InnerInstruction struct {
	// Index of the transaction instruction from which the inner instruction(s) originated
	Index uint16 `json:"index"`

	// Ordered list of inner program instructions that were invoked during a single transaction instruction.
	Instructions []CompiledInstruction `json:"instructions"`
}

type LoadedAddresses struct {
	ReadOnly common.Address `json:"readonly,omitempty"`
	Writable common.Address `json:"writable,omitempty"`
}

type TxStatus struct {
	Ok  interface{} `json:"Ok"`
	Err interface{} `json:"Err"`
}

type TransactionMeta struct {
	// TODO if has zero ComputeUnitsConsumed
	ComputeUnitsConsumed *uint64 `json:"computeUnitsConsumed"`
	// Error if transaction failed, null if transaction succeeded.
	// https://github.com/solana-labs/solana/blob/master/sdk/src/transaction.rs#L24
	Err interface{} `json:"err"`

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

type MessageHeader struct {
	NumRequiredSignatures       uint64 `json:"numRequiredSignatures"`
	NumReadonlySignedAccounts   uint64 `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts uint64 `json:"numReadonlyUnsignedAccounts"`
}

type Message struct {
	AccountKeys     []common.Address      `json:"accountKeys"`
	Header          MessageHeader         `json:"header"`
	RecentBlockHash common.Hash           `json:"recentBlockHash"`
	Instructions    []CompiledInstruction `json:"instructions"`
}

type TransactionContent struct {
	//
	Message Message `json:"message"`

	Signatures []common.Signature `json:"signatures"`
}

type BlockTransaction struct {
	// Transaction status metadata object
	Meta TransactionMeta `json:"meta"`
	// The slot this transaction was processed in.
	Slot uint64 `json:"slot"`

	// Estimated production time, as Unix timestamp (seconds since the Unix epoch)
	// of when the transaction was processed.
	// Nil if not available.
	BlockTime *int64 `json:"blockTime" bin:"optional"`

	// Transaction
	Transaction TransactionContent `json:"transaction"`
}

type BlockInfo struct {
	BlockHeight       uint64             `json:"blockHeight"`
	BlockTime         int64              `json:"blockTime"`
	ParentSlot        uint64             `json:"parentSlot"`
	BlockHash         common.Hash        `json:"blockHash"`
	PreviousBlockhash common.Hash        `json:"previousBlockhash"`
	Rewards           []BlockReward      `json:"rewards"`
	BlockTransaction  []BlockTransaction `json:"transactions"`
}

type BlockCommitment struct {
	// nil if Unknown block, or array of u64 integers
	// logging the amount of cluster stake in lamports
	// that has voted on the block at each depth from 0 to `MAX_LOCKOUT_HISTORY` + 1
	Commitment []uint64 `json:"commitment"`

	// Total active stake, in lamports, of the current epoch.
	TotalStake uint64 `json:"totalStake"`
}
