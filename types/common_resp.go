// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package types

import (
	"encoding/json"
	"github.com/cielu/go-solana/common"
	"math/big"
)

type ContextSlot struct {
	Slot       uint64 `json:"slot"`
	ApiVersion string `json:"apiVersion,omitempty"`
}

type AccountInfo struct {
	// data associated with the account, either as encoded binary data or JSON format {<program>: <state>} - depending on encoding parameter
	Data common.SolData `json:"data"`
	// base-58 encoded Pubkey of the program this account has been assigned to
	Owner common.Address `json:"owner"`
	// number of lamports assigned to this account, as an u64
	Lamports *big.Int `json:"lamports"`
	// the epoch at which this account will next owe rent, as u64
	RentEpoch *big.Int `json:"rentEpoch"`
	// boolean indicating if the account contains a program (and is strictly read-only)
	Executable bool `json:"executable"`
	// the data size of the account
	Space uint64 `json:"space,omitempty"`
	//  the data size of the account
	Size uint64 `json:"size,omitempty"`
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
	Lamports    *big.Int       `json:"lamports"`
	PostBalance uint64         `json:"postBalance"`
	RewardType  string         `json:"rewardType"`
	Pubkey      common.Address `json:"pubkey"`
}

type UiTokenAmount struct {
	// Address account
	Address *common.Address `json:"address,omitempty"`

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
	Data common.SolData `json:"data"`
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
	Ok  interface{}     `json:"Ok"`
	Err json.RawMessage `json:"Err"`
}

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
	Err               json.RawMessage    `json:"err"`
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

type BlockProduction struct {
	ByIdentity map[string][2]uint `json:"byIdentity"`
	Range      SlotRange          `json:"range"`
}

type BlockProductionWithCtx struct {
	Context         ContextSlot     `json:"context"`
	BlockProduction BlockProduction `json:"value"`
}

type ClusterInformation struct {
	// Node public key, as base-58 encoded string
	PubKey common.Address `json:"pubKey"`
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

type U64ValueWithCtx struct {
	Context ContextSlot `json:"context"`
	Value   *uint64     `json:"value,omitempty"`
}

type HighestSnapshotSlot struct {
	Full        uint64  `json:"full"`
	Incremental *uint64 `json:"incremental,omitempty"`
}

type Identity struct {
	Identity common.Address `json:"identity"`
}

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

type AccountWithLamport struct {
	// base-58 encoded address of the account
	Address common.Address `json:"address"`
	// number of lamports in the account, as a u64
	Lamports *big.Int `json:"lamports"`
}

type LastBlock struct {
	// a Hash as base-58 encoded string
	Blockhash common.Hash `json:"blockhash"`
	//  last block height at which the blockhash will be valid
	LastValidBlockHeight uint64 `json:"lastValidBlockHeight"`
}

type LastBlockWithCtx struct {
	Context   ContextSlot `json:"context"`
	LastBlock LastBlock   `json:"value"`
}

type AccountsInfoWithCtx struct {
	Context  ContextSlot   `json:"context"`
	Accounts []AccountInfo `json:"value,omitempty"`
}

// ProgramAccount program account
type ProgramAccount struct {
	Account AccountInfo    `json:"account"`
	PubKey  common.Address `json:"pubKey"`
}

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

type RpcPrioritizationFee struct {
	// slot in which the fee was observed
	Slot uint64 `json:"slot"`
	// the per-compute-unit fee paid by at least one successfully landed transaction, specified in increments of micro-lamports (0.000001 lamports)
	PrioritizationFee uint64 `json:"prioritizationFee"`
}

type SignatureInfo struct {
	// transaction signature as base-58 encoded string
	Signature common.Signature `json:"signature,omitempty"`
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

type StakeActivation struct {
	// the stake account's activation state, either: active, inactive, activating, or deactivating
	State string `json:"state"`
	// stake active during the epoch
	Active uint64 `json:"active"`
	// stake inactive during the epoch
	Inactive uint64 `json:"inactive"`
}

type SupplyInfo struct {
	// Total supply in lamports
	Total uint64 `json:"total"`
	// Circulating supply in lamports
	Circulating uint64 `json:"circulating"`
	// Non-circulating supply in lamports
	NonCirculating uint64 `json:"nonCirculating"`
	// an array of account addresses of non-circulating accounts, as strings. If excludeNonCirculatingAccountsList is enabled, the returned array will be empty.
	NonCirculatingAccounts []common.Address `json:"nonCirculatingAccounts"`
}

type SupplyWithCtx struct {
	Context ContextSlot `json:"context"`
	Supply  SupplyInfo  `json:"value"`
}

type TokenAccountWithCtx struct {
	Context ContextSlot   `json:"context"`
	UiToken UiTokenAmount `json:"value"`
}

type TokenAccount struct {
	Account AccountInfo    `json:"account"`
	Pubkey  common.Address `json:"pubkey,omitempty"`
}

type TokenAccountsWithCtx struct {
	Context  ContextSlot    `json:"context"`
	Accounts []TokenAccount `json:"value"`
}

type TokenLargestHolders struct {
	Context ContextSlot     `json:"context"`
	Holders []UiTokenAmount `json:"value"`
}

type SolVersion struct {
	// software version of solana-core as a string
	SolanaCore string `json:"solana-core"`
	// unique identifier of the current software's feature set as a u32
	FeatureSet uint32 `json:"feature-set"`
}

type VoteAccount struct {
	// Vote account address, as base-58 encoded string
	VotePubkey common.Address `json:"votePubkey"`
	// Validator identity, as base-58 encoded string
	NodePubkey common.Address
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

type RpcVoteAccounts struct {
	Current    []VoteAccount `json:"current"`
	Delinquent []VoteAccount `json:"delinquent"`
}

type BlockInfoWithCtx struct {
	Context   ContextSlot `json:"context"`
	BlockInfo BlockInfo   `json:"value"`
}

type LogsValue struct {
	Signature common.Signature `json:"signature"`
	Err       json.RawMessage  `json:"err"`
	Logs      []string         `json:"logs"`
}

type LogsInfoWithCtx struct {
	Context ContextSlot `json:"context"`
	Value   LogsValue   `json:"value"`
}

type ProgramInfoWithCtx struct {
	Context ContextSlot  `json:"context"`
	Value   TokenAccount `json:"value"`
}

type SignatureInfoWithCtx struct {
	Context ContextSlot `json:"context"`
	Value   string      `json:"value"`
}

type SlotInfoWithCtx struct {
	Parent uint64 `json:"parent"`
	Root   uint64 `json:"root"`
	Slot   uint64 `json:"slot"`
}
