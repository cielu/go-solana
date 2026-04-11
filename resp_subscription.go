// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

import "encoding/json"

// LogsValue is the payload of a single logsSubscribe notification.
type LogsValue struct {
	Signature Signature       `json:"signature"`
	Err       json.RawMessage `json:"err"`
	Logs      []string        `json:"logs"`
}

// AccountNotifies is the payload of accountSubscribe (identical shape to
// AccountInfoWithCtx).
type AccountNotifies AccountInfoWithCtx

// BlockNotifies is the payload of blockSubscribe.
type BlockNotifies struct {
	Context   ContextSlot `json:"context"`
	BlockInfo BlockInfo   `json:"value"`
}

// LogsNotifies is the payload of logsSubscribe.
type LogsNotifies struct {
	Context ContextSlot `json:"context"`
	Value   LogsValue   `json:"value"`
}

// ProgramNotifies is the payload of programSubscribe.
type ProgramNotifies struct {
	Context ContextSlot  `json:"context"`
	Value   TokenAccount `json:"value"`
}

// SignatureNotifies is the payload of signatureSubscribe.
type SignatureNotifies struct {
	Context ContextSlot `json:"context"`
	Value   interface{} `json:"value"`
}

// SlotNotifies is the payload of slotSubscribe.
type SlotNotifies struct {
	Parent uint64 `json:"parent"`
	Root   uint64 `json:"root"`
	Slot   uint64 `json:"slot"`
}

// RootNotifies carries the latest rooted slot reported by rootSubscribe.
type RootNotifies uint64

// SlotsUpdatesNotifies is the payload for slotsUpdatesSubscribe. See the
// Solana RPC docs for the full set of possible Type values; commonly
// "firstShredReceived", "completed", "createdBank", "frozen", "dead",
// "optimisticConfirmation", and "root".
type SlotsUpdatesNotifies struct {
	Parent    uint64                `json:"parent"`
	Slot      uint64                `json:"slot"`
	Timestamp int64                 `json:"timestamp"`
	Type      string                `json:"type"`
	Stats     *SlotsUpdatesBankStat `json:"stats,omitempty"`
	Err       string                `json:"err,omitempty"`
}

// SlotsUpdatesBankStat is the optional per-bank statistics attached to
// some SlotsUpdatesNotifies events.
type SlotsUpdatesBankStat struct {
	MaxTransactionsPerEntry   uint64 `json:"maxTransactionsPerEntry"`
	NumFailedTransactions     uint64 `json:"numFailedTransactions"`
	NumSuccessfulTransactions uint64 `json:"numSuccessfulTransactions"`
	NumTransactionEntries     uint64 `json:"numTransactionEntries"`
}
