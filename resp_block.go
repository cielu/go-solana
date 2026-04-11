// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

import "encoding/json"

// BlockInfo is the response shape of getBlock.
type BlockInfo struct {
	Err               json.RawMessage   `json:"err"`
	BlockHeight       uint64            `json:"blockHeight"`
	BlockTime         int64             `json:"blockTime"`
	ParentSlot        uint64            `json:"parentSlot"`
	BlockHash         Hash              `json:"blockHash"`
	PreviousBlockhash Hash              `json:"previousBlockhash"`
	Rewards           []BlockReward     `json:"rewards"`
	Transactions      []TransactionInfo `json:"transactions"`
}

// BlockCommitment is the response shape of getBlockCommitment.
type BlockCommitment struct {
	// nil if Unknown block, or array of u64 integers
	// logging the amount of cluster stake in lamports
	// that has voted on the block at each depth from 0 to `MAX_LOCKOUT_HISTORY` + 1
	Commitment []uint64 `json:"commitment"`

	// Total active stake, in lamports, of the current epoch.
	TotalStake uint64 `json:"totalStake"`
}

// BlockProduction is the value portion of the getBlockProduction response.
type BlockProduction struct {
	ByIdentity map[string][2]uint `json:"byIdentity"`
	Range      SlotRange          `json:"range"`
}

// BlockProductionWithCtx is the context-wrapped response shape of
// getBlockProduction.
type BlockProductionWithCtx struct {
	Context         ContextSlot     `json:"context"`
	BlockProduction BlockProduction `json:"value"`
}

// LastBlockWithCtx is the context-wrapped response shape of getLatestBlockhash.
type LastBlockWithCtx struct {
	Context   ContextSlot `json:"context"`
	LastBlock LastBlock   `json:"value"`
}

// HighestSnapshotSlot is the response shape of getHighestSnapshotSlot.
type HighestSnapshotSlot struct {
	Full        uint64  `json:"full"`
	Incremental *uint64 `json:"incremental,omitempty"`
}
