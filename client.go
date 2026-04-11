// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cielu/go-solana/rpc"
)

// Client defines typed wrappers for the solana RPC API.
type Client struct {
	c *rpc.Client
}

// call invokes the named RPC method and wraps any returned error with the
// method name for easier debugging. All typed wrappers funnel through this.
func (sc *Client) call(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	if err := sc.c.CallContext(ctx, result, method, args...); err != nil {
		return fmt.Errorf("solana rpc %s: %w", method, err)
	}
	return nil
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*Client, error) {
	return DialContext(context.Background(), rawurl)
}

// DialContext connects a client to the given URL with context.
func DialContext(ctx context.Context, rawurl string) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

// NewClient creates a client that uses the given RPC client.
func NewClient(c *rpc.Client) *Client {
	return &Client{c}
}

// SetDebug set solClient debug
func (sc *Client) SetDebug(isDebug bool) {
	sc.c.IsDebug = isDebug
}

// Close closes the underlying RPC connection.
func (sc *Client) Close() {
	sc.c.Close()
}

// --------------------------------------------------------
// --------------------------------------------------------
// ------------------Blockchain Access---------------------
// ------------------Blockchain Access---------------------
// --------------------------------------------------------
// --------------------------------------------------------

func getRpcCfg[T any](cfg []T) *T {
	// never set rpc ctx cfg
	if len(cfg) == 0 {
		return nil
	}
	return &cfg[0]
}

// GetAccountInfo Returns all information associated with the account of provided Pubkey
func (sc *Client) GetAccountInfo(ctx context.Context, account PublicKey, cfg ...RpcAccountInfoCfg) (res AccountInfoWithCtx, err error) {
	err = sc.call(ctx, &res, "getAccountInfo", account, getRpcCfg(cfg))
	return
}

// GetBalance Returns the lamport balance of the account of provided Pubkey
func (sc *Client) GetBalance(ctx context.Context, account PublicKey, cfg ...RpcCommitmentWithMinSlotCfg) (balance BalanceWithCtx, err error) {
	err = sc.call(ctx, &balance, "getBalance", account, getRpcCfg(cfg))
	return
}

// GetBlock Returns identity and transaction information about a confirmed block in the ledger
func (sc *Client) GetBlock(ctx context.Context, blockNum uint64, cfg ...RpcGetBlockContextCfg) (blockInfo BlockInfo, err error) {
	c := getRpcCfg(cfg)
	if c == nil {
		c = &RpcGetBlockContextCfg{}
	}
	// set Encoding as base64
	// if c.Encoding == "" {
	// 	c.Encoding = types.EncodingBase64
	// }
	err = sc.call(ctx, &blockInfo, "getBlock", blockNum, c)
	return
}

// GetBlockCommitment Returns commitment for particular block
func (sc *Client) GetBlockCommitment(ctx context.Context, blockNum uint64) (blockCmt BlockCommitment, err error) {
	err = sc.call(ctx, &blockCmt, "getBlockCommitment", blockNum)
	return
}

// GetBlockHeight Returns the current block height of the node
func (sc *Client) GetBlockHeight(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res uint64, err error) {
	err = sc.call(ctx, &res, "getBlockHeight", getRpcCfg(cfg))
	return
}

// GetBlockProduction Returns recent block production information from the current or previous epoch.
func (sc *Client) GetBlockProduction(ctx context.Context, cfg ...RpcGetBlockProduction) (res BlockProductionWithCtx, err error) {
	err = sc.call(ctx, &res, "getBlockProduction", getRpcCfg(cfg))
	return
}

// GetBlockTime Returns the estimated production time of a block.
func (sc *Client) GetBlockTime(ctx context.Context, blockNum uint64) (res int64, err error) {
	err = sc.call(ctx, &res, "getBlockTime", blockNum)
	return
}

// GetBlocks Returns a list of confirmed blocks between two slots.
// The optional slot range [startSlot, endSlot] must be at most 500000 slots wide.
// Options: WithBlocksEndSlot, WithBlocksCommitment.
func (sc *Client) GetBlocks(ctx context.Context, startSlot uint64, opts ...BlocksOption) ([]uint64, error) {
	c := &blocksConfig{}
	for _, opt := range opts {
		opt(c)
	}
	if c.endSlot != nil {
		if *c.endSlot < startSlot {
			return nil, fmt.Errorf("getBlocks: endSlot (%d) must be >= startSlot (%d)", *c.endSlot, startSlot)
		}
		if *c.endSlot-startSlot >= 500000 {
			return nil, fmt.Errorf("getBlocks: slot range must be <= 500000 (got %d)", *c.endSlot-startSlot)
		}
	}
	var res []uint64
	if err := sc.call(ctx, &res, "getBlocks", startSlot, c.endSlot, c.cfg); err != nil {
		return nil, err
	}
	return res, nil
}

// GetBlocksWithLimit Returns a list of confirmed blocks starting at the given slot
func (sc *Client) GetBlocksWithLimit(ctx context.Context, startSlot, limit uint64, cfg ...RpcCommitmentCfg) (res []uint64, err error) {
	err = sc.call(ctx, &res, "getBlocksWithLimit", startSlot, limit, getRpcCfg(cfg))
	return
}

// GetClusterNodes Returns information about all the nodes participating in the cluster
func (sc *Client) GetClusterNodes(ctx context.Context) (res []ClusterInformation, err error) {
	err = sc.call(ctx, &res, "getClusterNodes")
	return
}

// GetEpochInfo Returns information about the current epoch
func (sc *Client) GetEpochInfo(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res EpochInformation, err error) {
	err = sc.call(ctx, &res, "getEpochInfo", getRpcCfg(cfg))
	return
}

// GetEpochSchedule Returns information about the current epoch
func (sc *Client) GetEpochSchedule(ctx context.Context) (res EpochSchedule, err error) {
	err = sc.call(ctx, &res, "getEpochSchedule")
	return
}

// GetFeeForMessage Get the fee the network will charge for a particular Message
func (sc *Client) GetFeeForMessage(ctx context.Context, msg string, cfg ...RpcCommitmentWithMinSlotCfg) (res U64ValueWithCtx, err error) {
	err = sc.call(ctx, &res, "getFeeForMessage", msg, getRpcCfg(cfg))
	return
}

// GetFirstAvailableBlock Returns the slot of the lowest confirmed block that has not been purged from the ledger
func (sc *Client) GetFirstAvailableBlock(ctx context.Context) (res uint64, err error) {
	err = sc.call(ctx, &res, "getFirstAvailableBlock")
	return
}

// GetGenesisHash Returns the genesis hash
func (sc *Client) GetGenesisHash(ctx context.Context) (res Hash, err error) {
	err = sc.call(ctx, &res, "getGenesisHash")
	return
}

// GetHealth Returns the current health of the node.
// A healthy node is one that is within HEALTH_CHECK_SLOT_DISTANCE slots of the latest cluster confirmed slot.
func (sc *Client) GetHealth(ctx context.Context) (res string, err error) {
	err = sc.call(ctx, &res, "getHealth")
	return
}

// GetHighestSnapshotSlot Returns the highest slot information that the node has snapshots for.
// This will find the highest full snapshot slot, and the highest incremental snapshot slot based on the full snapshot slot, if there is one.
func (sc *Client) GetHighestSnapshotSlot(ctx context.Context) (res HighestSnapshotSlot, err error) {
	err = sc.call(ctx, &res, "getHighestSnapshotSlot")
	return
}

// GetIdentity Returns the identity pubkey for the current node
func (sc *Client) GetIdentity(ctx context.Context) (res Identity, err error) {
	err = sc.call(ctx, &res, "getIdentity")
	return
}

// GetInflationGovernor Returns the current inflation governor
func (sc *Client) GetInflationGovernor(ctx context.Context, cfg ...RpcCommitmentCfg) (res InflationGovernor, err error) {
	err = sc.call(ctx, &res, "getInflationGovernor", getRpcCfg(cfg))
	return
}

// GetInflationRate Returns the specific inflation values for the current epoch
func (sc *Client) GetInflationRate(ctx context.Context) (res InflationRate, err error) {
	err = sc.call(ctx, &res, "getInflationRate")
	return
}

// GetInflationReward Returns the inflation / staking reward for a list of
// PublicKeys for an epoch.
// Options: WithRewardAccounts, WithRewardCommitment.
func (sc *Client) GetInflationReward(ctx context.Context, opts ...RewardOption) ([]InflationReward, error) {
	c := &rewardConfig{}
	for _, opt := range opts {
		opt(c)
	}
	var res []InflationReward
	if err := sc.call(ctx, &res, "getInflationReward", c.accounts, c.cfg); err != nil {
		return nil, err
	}
	return res, nil
}

// GetLargestAccounts Returns the 20 largest accounts, by lamport balance (results may be cached up to two hours)
func (sc *Client) GetLargestAccounts(ctx context.Context, cfg ...RpcCommitmentWithFilter) (res AccountWithLamport, err error) {
	err = sc.call(ctx, &res, "getLargestAccounts", getRpcCfg(cfg))
	return
}

// GetLatestBlockhash Returns the latest blockhash
func (sc *Client) GetLatestBlockhash(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res LastBlockWithCtx, err error) {
	err = sc.call(ctx, &res, "getLatestBlockhash", getRpcCfg(cfg))
	return
}

// GetLeaderSchedule Returns the leader schedule for an epoch. By default this
// is the current epoch; pass WithScheduleSlot to select the epoch that contains
// a specific slot.
// Options: WithScheduleSlot, WithScheduleCommitment, WithScheduleIdentity.
func (sc *Client) GetLeaderSchedule(ctx context.Context, opts ...ScheduleOption) (map[string][]uint64, error) {
	c := &scheduleConfig{}
	for _, opt := range opts {
		opt(c)
	}
	var res map[string][]uint64
	if err := sc.call(ctx, &res, "getLeaderSchedule", c.slot, c.cfg); err != nil {
		return nil, err
	}
	return res, nil
}

// GetMaxRetransmitSlot Get the max slot seen from retransmit stage.
func (sc *Client) GetMaxRetransmitSlot(ctx context.Context) (res uint64, err error) {
	err = sc.call(ctx, &res, "getMaxRetransmitSlot")
	return
}

// GetMaxShredInsertSlot Get the max slot seen from after shred insert.
func (sc *Client) GetMaxShredInsertSlot(ctx context.Context) (res uint64, err error) {
	err = sc.call(ctx, &res, "getMaxShredInsertSlot")
	return
}

// GetMinimumBalanceForRentExemption Returns the minimum lamport balance required
// to make an account of the given data size rent-exempt.
// Options: WithRentExemptDataLen, WithRentExemptCommitment.
func (sc *Client) GetMinimumBalanceForRentExemption(ctx context.Context, opts ...RentExemptOption) (uint64, error) {
	c := &rentExemptConfig{}
	for _, opt := range opts {
		opt(c)
	}
	var res uint64
	if err := sc.call(ctx, &res, "getMinimumBalanceForRentExemption", c.dataLen, c.cfg); err != nil {
		return 0, err
	}
	return res, nil
}

// GetMultipleAccounts Returns the account information for a list of Pubkeys.
// Up to 100 accounts per call.
func (sc *Client) GetMultipleAccounts(ctx context.Context, accounts []PublicKey, cfg ...RpcAccountInfoCfg) (res AccountsInfoWithCtx, err error) {
	if len(accounts) > 100 {
		return res, fmt.Errorf("getMultipleAccounts: at most 100 accounts allowed (got %d)", len(accounts))
	}
	err = sc.call(ctx, &res, "getMultipleAccounts", accounts, getRpcCfg(cfg))
	return
}

// GetProgramAccounts Returns all accounts owned by the provided program Pubkey
func (sc *Client) GetProgramAccounts(ctx context.Context, program PublicKey, cfg ...RpcCombinedCfg) (res []ProgramAccount, err error) {
	err = sc.call(ctx, &res, "getProgramAccounts", program, getRpcCfg(cfg))
	return
}

// GetRecentPerformanceSamples Returns a list of recent performance samples, in reverse slot order.
// Performance samples are taken every 60 seconds and include the number of transactions and slots that occur in a given time window.
// Optional arg: limit (max 720).
func (sc *Client) GetRecentPerformanceSamples(ctx context.Context, args ...uint64) (res []RpcPerfSample, err error) {
	var limit *uint64
	if len(args) > 0 {
		if args[0] > 720 {
			return res, fmt.Errorf("getRecentPerformanceSamples: limit must be <= 720 (got %d)", args[0])
		}
		limit = &args[0]
	}
	err = sc.call(ctx, &res, "getRecentPerformanceSamples", limit)
	return
}

// GetRecentPrioritizationFees Returns a list of prioritization fees from recent
// blocks, optionally scoped to a set of accounts.
// Options: WithFeesAccounts (at most 128 accounts total).
func (sc *Client) GetRecentPrioritizationFees(ctx context.Context, opts ...FeesOption) ([]RpcPrioritizationFee, error) {
	c := &feesConfig{}
	for _, opt := range opts {
		opt(c)
	}
	if len(c.accounts) > 128 {
		return nil, fmt.Errorf("getRecentPrioritizationFees: at most 128 accounts allowed (got %d)", len(c.accounts))
	}
	var res []RpcPrioritizationFee
	if err := sc.call(ctx, &res, "getRecentPrioritizationFees", c.accounts); err != nil {
		return nil, err
	}
	return res, nil
}

// GetSignatureStatuses Returns the statuses of a list of signatures. Each signature must be a txid, the first signature of a transaction.
// Up to 256 signatures per call.
func (sc *Client) GetSignatureStatuses(ctx context.Context, signatures []Signature, cfg ...RpcSearchTxHistoryCfg) (res SignatureStatusWithCtx, err error) {
	if len(signatures) > 256 {
		return res, fmt.Errorf("getSignatureStatuses: at most 256 signatures allowed (got %d)", len(signatures))
	}
	err = sc.call(ctx, &res, "getSignatureStatuses", signatures, getRpcCfg(cfg))
	return
}

// GetSignaturesForAddress Returns signatures for confirmed transactions that include the given PublicKey in their accountKeys list.
// Returns signatures backwards in time from the provided signature or most recent confirmed block
func (sc *Client) GetSignaturesForAddress(ctx context.Context, account PublicKey, cfg ...RpcSignaturesForAddressCfg) (res []SignatureInfo, err error) {
	err = sc.call(ctx, &res, "getSignaturesForAddress", account, getRpcCfg(cfg))
	return
}

// GetSlot Returns the slot that has reached the given or default commitment level
// https://solana.com/docs/rpc#configuring-state-commitment
func (sc *Client) GetSlot(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res uint64, err error) {
	err = sc.call(ctx, &res, "getSlot", getRpcCfg(cfg))
	return
}

// GetSlotLeader Returns the current slot leader
func (sc *Client) GetSlotLeader(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res PublicKey, err error) {
	err = sc.call(ctx, &res, "getSlotLeader", getRpcCfg(cfg))
	return
}

// GetSlotLeaders Returns the slot leaders for a given slot range.
// args: startSlot[, limit]; limit must be <= 5000.
func (sc *Client) GetSlotLeaders(ctx context.Context, args ...uint64) (res []PublicKey, err error) {
	var startSlot, limit *uint64
	switch len(args) {
	case 1:
		startSlot = &args[0]
	case 2:
		startSlot = &args[0]
		limit = &args[1]
		if *limit > 5000 {
			return res, fmt.Errorf("getSlotLeaders: limit must be <= 5000 (got %d)", *limit)
		}
	default:
		return res, fmt.Errorf("getSlotLeaders: expected 1 or 2 args (startSlot[, limit]), got %d", len(args))
	}
	err = sc.call(ctx, &res, "getSlotLeaders", startSlot, limit)
	return
}

// GetStakeActivation Returns epoch activation information for a stake account
func (sc *Client) GetStakeActivation(ctx context.Context, account PublicKey, cfg ...RpcCommitmentWithMinSlotCfg) (res StakeActivation, err error) {
	err = sc.call(ctx, &res, "getStakeActivation", account, getRpcCfg(cfg))
	return
}

// GetStakeMinimumDelegation Returns the stake minimum delegation, in lamports.
func (sc *Client) GetStakeMinimumDelegation(ctx context.Context, cfg ...RpcCommitmentCfg) (res U64ValueWithCtx, err error) {
	err = sc.call(ctx, &res, "getStakeMinimumDelegation", getRpcCfg(cfg))
	return
}

// GetSupply Returns information about the current supply.
func (sc *Client) GetSupply(ctx context.Context, cfg ...RpcSupplyCfg) (res SupplyWithCtx, err error) {
	err = sc.call(ctx, &res, "getSupply", getRpcCfg(cfg))
	return
}

// GetTokenAccountBalance Returns the token balance of an SPL Token account.
func (sc *Client) GetTokenAccountBalance(ctx context.Context, account PublicKey, cfg ...RpcCommitmentCfg) (res TokenAccountWithCtx, err error) {
	err = sc.call(ctx, &res, "getTokenAccountBalance", account, getRpcCfg(cfg))
	return
}

// GetTokenAccountsByDelegate Returns all SPL Token accounts by approved Delegate.
func (sc *Client) GetTokenAccountsByDelegate(ctx context.Context, delegate PublicKey, mintProg RpcMintWithProgramID, cfg ...RpcAccountInfoCfg) (res TokenAccountsWithCtx, err error) {
	// `params` should have at least 2 argument(s)
	err = sc.call(ctx, &res, "getTokenAccountsByDelegate", delegate, mintProg, getRpcCfg(cfg))
	return
}

// GetTokenAccountsByOwner Returns all SPL Token accounts by token owner.
func (sc *Client) GetTokenAccountsByOwner(ctx context.Context, owner PublicKey, program RpcMintWithProgramID, cfg ...RpcAccountInfoCfg) (res TokenAccountsWithCtx, err error) {
	// use base64
	tmpCfg := getRpcCfg(cfg)
	// isNull
	if tmpCfg == nil {
		tmpCfg = &RpcAccountInfoCfg{}
	}
	// set encoding to base64
	if tmpCfg.Encoding == "" {
		tmpCfg.Encoding = EncodingBase64
	}
	err = sc.call(ctx, &res, "getTokenAccountsByOwner", owner, program, tmpCfg)
	return
}

// GetTokenLargestAccounts Returns the 20 largest accounts of a particular SPL Token type.
func (sc *Client) GetTokenLargestAccounts(ctx context.Context, splToken PublicKey, cfg ...RpcCommitmentCfg) (res TokenLargestHolders, err error) {
	err = sc.call(ctx, &res, "getTokenLargestAccounts", splToken, getRpcCfg(cfg))
	return
}

// GetTokenSupply Returns the total supply of an SPL Token type.
func (sc *Client) GetTokenSupply(ctx context.Context, splToken PublicKey, cfg ...RpcCommitmentCfg) (res TokenAccountWithCtx, err error) {
	err = sc.call(ctx, &res, "getTokenSupply", splToken, getRpcCfg(cfg))
	return
}

// GetTransaction Returns transaction details for a confirmed transaction
func (sc *Client) GetTransaction(ctx context.Context, signature Signature, cfg ...RpcGetTransactionCfg) (res TransactionInfo, err error) {
	err = sc.call(ctx, &res, "getTransaction", signature, getRpcCfg(cfg))
	return
}

// GetTransactionCount Returns the current Transaction count from the ledger
func (sc *Client) GetTransactionCount(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res uint64, err error) {
	err = sc.call(ctx, &res, "getTransactionCount", getRpcCfg(cfg))
	return
}

// GetVersion Returns the current Solana version running on the node
func (sc *Client) GetVersion(ctx context.Context) (res SolVersion, err error) {
	err = sc.call(ctx, &res, "getVersion")
	return
}

// GetVoteAccounts Returns the account info and associated stake for all the voting accounts in the current bank.
func (sc *Client) GetVoteAccounts(ctx context.Context, cfg ...RpcVoteAccountCfg) (res RpcVoteAccounts, err error) {
	err = sc.call(ctx, &res, "getVoteAccounts", getRpcCfg(cfg))
	return
}

// IsBlockHashValid Returns whether a blockHash is still valid or not
func (sc *Client) IsBlockHashValid(ctx context.Context, hash Hash, cfg ...RpcCommitmentWithMinSlotCfg) (res bool, err error) {
	err = sc.call(ctx, &res, "isBlockhashValid", hash, getRpcCfg(cfg))
	return
}

// MinimumLedgerSlot Returns the lowest slot that the node has information about in its ledger.
func (sc *Client) MinimumLedgerSlot(ctx context.Context) (res uint64, err error) {
	err = sc.call(ctx, &res, "minimumLedgerSlot")
	return
}

// RequestAirdrop Requests an airdrop of lamports to a Pubkey
func (sc *Client) RequestAirdrop(ctx context.Context, address PublicKey, lamport *big.Int) (res Signature, err error) {
	err = sc.call(ctx, &res, "requestAirdrop", address, lamport)
	return
}

// SendTransaction  Submits a signed transaction to the cluster for processing.
// This method does not alter the transaction in any way; it relays the transaction created by clients to the node as-is.
// If the node's rpc service receives the transaction, this method immediately succeeds, without waiting for any confirmations. A successful response from this method does not guarantee the transaction is processed or confirmed by the cluster.
// While the rpc service will reasonably retry to submit it, the transaction could be rejected if transaction's recent_blockhash expires before it lands.
// Use getSignatureStatuses to ensure a transaction is processed and confirmed.
// Before submitting, the following preflight checks are performed:
// The transaction signatures are verified
// The transaction is simulated against the bank slot specified by the preflight commitment. On failure an error will be returned. Preflight checks may be disabled if desired. It is recommended to specify the same commitment and preflight commitment to avoid confusing behavior.
// The returned signature is the first signature in the transaction, which is used to identify the transaction (transaction id). This identifier can be easily extracted from the transaction data before submission.
func (sc *Client) SendTransaction(ctx context.Context, signedTx Base58Data, cfg ...RpcSendTxCfg) (res Signature, err error) {
	err = sc.call(ctx, &res, "sendTransaction", signedTx, getRpcCfg(cfg))
	return
}

// SimulateTransaction Simulate sending a transaction
func (sc *Client) SimulateTransaction(ctx context.Context, signedTx Base58Data) (res SimulateResult, err error) {
	err = sc.call(ctx, &res, "simulateTransaction", signedTx)
	return
}
