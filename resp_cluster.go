// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

// ClusterInformation is an element of the getClusterNodes response.
type ClusterInformation struct {
	// Node public key, as base-58 encoded string
	PubKey PublicKey `json:"pubKey"`
	// Gossip network address for the node
	Gossip string `json:"gossip,omitempty"`
	// TPU network address for the node
	Tpu string `json:"tpu,omitempty"`
	// JSON RPC network address for the node, or null if the JSON RPC service is not enabled
	Rpc string `json:"rpc,omitempty"`
	// The software version of the node, or null if the version information is not available
	Version string `json:"version,omitempty"`
	// The unique identifier of the node's feature set
	FeatureSet *uint32 `json:"featureSet,omitempty"`
	// The shred version the node has been configured to use
	ShredVersion *uint16 `json:"shredVersion,omitempty"`
}

// EpochInformation is the response shape of getEpochInfo.
type EpochInformation struct {
	// the current slot
	AbsoluteSlot uint64 `json:"absoluteSlot"`
	// the current block height
	BlockHeight uint64 `json:"blockHeight"`
	// the current epoch
	Epoch uint64 `json:"epoch"`
	// the current slot relative to the start of the current epoch
	SlotIndex uint64 `json:"slotIndex"`
	// the number of slots in this epoch
	SlotsInEpoch uint64 `json:"slotsInEpoch"`
	// total number of transactions processed without error since genesis
	TransactionCount *uint64 `json:"transactionCount"`
}

// EpochSchedule is the response shape of getEpochSchedule.
type EpochSchedule struct {
	// the maximum number of slots in each epoch
	SlotsPerEpoch uint64 `json:"slotsPerEpoch"`
	// the number of slots before beginning of an epoch to calculate a leader schedule for that epoch
	LeaderScheduleSlotOffset uint64 `json:"leaderScheduleSlotOffset"`
	// whether epochs start short and grow
	Warmup bool `json:"warmup"`
	// first normal-length epoch, log2(slotsPerEpoch) - log2(MINIMUM_SLOTS_PER_EPOCH)
	FirstNormalEpoch uint64 `json:"firstNormalEpoch"`
	// MINIMUM_SLOTS_PER_EPOCH * (2.pow(firstNormalEpoch) - 1)
	FirstNormalSlot uint64 `json:"firstNormalSlot"`
}

// Identity is the response shape of getIdentity.
type Identity struct {
	Identity PublicKey `json:"identity"`
}

// InflationGovernor is the response shape of getInflationGovernor.
type InflationGovernor struct {
	// the initial inflation percentage from time 0
	Initial float64 `json:"initial"`
	// terminal inflation percentage
	Terminal float64 `json:"terminal"`
	// rate per year at which inflation is lowered. (Rate reduction is derived using the target slot time in genesis config)
	Taper float64 `json:"taper"`
	// percentage of total inflation allocated to the foundation
	Foundation float64 `json:"foundation"`
	// duration of foundation pool inflation in years
	FoundationTerm float64 `json:"foundationTerm"`
}

// InflationRate is the response shape of getInflationRate.
type InflationRate struct {
	// total inflation
	Total float64 `json:"total"`
	// inflation allocated to validators
	Validator float64 `json:"validator"`
	// inflation allocated to the foundation
	Foundation float64 `json:"foundation"`
	// epoch for which these values are valid
	Epoch uint64 `json:"epoch"`
}

// InflationReward is an element of the getInflationReward response.
type InflationReward struct {
	// epoch for which reward occured
	Epoch uint64 `json:"epoch"`
	// the slot in which the rewards are effective
	EffectiveSlot uint64 `json:"effectiveSlot"`
	// reward amount in lamports
	Amount uint64 `json:"amount"`
	// post balance of the account in lamports
	PostBalance uint64 `json:"postBalance"`
	// vote account commission when the reward was credited
	Commission *uint8 `json:"commission,omitempty"`
}

// RpcPerfSample is an element of the getRecentPerformanceSamples response.
type RpcPerfSample struct {
	// Slot in which sample was taken at
	Slot uint64 `json:"slot"`
	// Number of transactions processed during the sample period
	NumTransactions uint64 `json:"numTransactions"`
	// Number of slots completed during the sample period
	NumSlots uint64 `json:"numSlots"`
	// Number of seconds in a sample window
	SamplePeriodSecs uint16 `json:"samplePeriodSecs"`
	// Number of non-vote transactions processed during the sample period.
	NumNonVoteTransaction uint64 `json:"numNonVoteTransaction"`
}

// RpcPrioritizationFee is an element of the getRecentPrioritizationFees
// response.
type RpcPrioritizationFee struct {
	// slot in which the fee was observed
	Slot uint64 `json:"slot"`
	// the per-compute-unit fee paid by at least one successfully landed transaction, specified in increments of micro-lamports (0.000001 lamports)
	PrioritizationFee uint64 `json:"prioritizationFee"`
}

// SupplyInfo is the value portion of the getSupply response.
type SupplyInfo struct {
	// Total supply in lamports
	Total uint64 `json:"total"`
	// Circulating supply in lamports
	Circulating uint64 `json:"circulating"`
	// Non-circulating supply in lamports
	NonCirculating uint64 `json:"nonCirculating"`
	// an array of account addresses of non-circulating accounts, as strings. If excludeNonCirculatingAccountsList is enabled, the returned array will be empty.
	NonCirculatingAccounts []PublicKey `json:"nonCirculatingAccounts"`
}

// SupplyWithCtx is the context-wrapped response shape of getSupply.
type SupplyWithCtx struct {
	Context ContextSlot `json:"context"`
	Supply  SupplyInfo  `json:"value"`
}

// SolVersion is the response shape of getVersion.
type SolVersion struct {
	// software version of solana-core as a string
	SolanaCore string `json:"solana-core"`
	// unique identifier of the current software's feature set as a u32
	FeatureSet uint32 `json:"feature-set"`
}

// VoteAccount is an element of RpcVoteAccounts.
type VoteAccount struct {
	// Vote account address, as base-58 encoded string
	VotePubkey PublicKey `json:"votePubkey"`
	// Validator identity, as base-58 encoded string
	NodePubkey PublicKey `json:"nodePubkey"`
	// the stake, in lamports, delegated to this vote account and active in this epoch
	ActivatedStake uint64 `json:"activatedStake"`
	// bool, whether the vote account is staked for this epoch
	EpochVoteAccount bool `json:"epochVoteAccount"`
	// percentage (0-100) of rewards payout owed to the vote account
	Commission uint8 `json:"commission"`
	// Most recent slot voted on by this vote account
	LastVote uint64 `json:"lastVote"`
	// Latest history of earned credits for up to five epochs, as an array of arrays containing: [epoch, credits, previousCredits].
	EpochCredits [][]uint64 `json:"epochCredits"`
	// Current root slot for this vote
	RootSlot uint64 `json:"rootSlot"`
}

// RpcVoteAccounts is the response shape of getVoteAccounts.
type RpcVoteAccounts struct {
	Current    []VoteAccount `json:"current"`
	Delinquent []VoteAccount `json:"delinquent"`
}
