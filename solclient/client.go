// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solclient

import (
	"context"
	"errors"
	"go-solana/common"
	"go-solana/rpc"
	"go-solana/types"
	"math/big"
)

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	c *rpc.Client
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

// Close closes the underlying RPC connection.
func (sc *Client) Close() {
	sc.c.Close()
}

// Client gets the underlying RPC client.
func (sc *Client) Client() *rpc.Client {
	return sc.c
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
func (sc *Client) GetAccountInfo(ctx context.Context, account common.Address, cfg ...types.RpcCommitmentCfg) (res types.AccountInfoWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getAccountInfo", account, getRpcCfg(cfg))
	return
}

// GetBalance Returns the lamport balance of the account of provided Pubkey
func (sc *Client) GetBalance(ctx context.Context, account common.Address, cfg ...types.RpcCommitmentAndMinSlotCfg) (balance types.BalanceWithCtx, err error) {
	err = sc.c.CallContext(ctx, &balance, "getBalance", account, getRpcCfg(cfg))
	return
}

// GetBlock Returns identity and transaction information about a confirmed block in the ledger
func (sc *Client) GetBlock(ctx context.Context, blockNum uint64, cfg ...types.RpcGetBlockContextCfg) (blockInfo types.BlockInfo, err error) {
	err = sc.c.CallContext(ctx, &blockInfo, "getBlock", blockNum, getRpcCfg(cfg))
	return
}

// GetBlockCommitment Returns commitment for particular block
func (sc *Client) GetBlockCommitment(ctx context.Context, blockNum uint64) (blockCmt types.BlockCommitment, err error) {
	err = sc.c.CallContext(ctx, &blockCmt, "getBlockCommitment", blockNum)
	return
}

// GetBlockHeight Returns the current block height of the node
func (sc *Client) GetBlockHeight(ctx context.Context, cfg ...types.RpcCommitmentAndMinSlotCfg) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlockHeight", getRpcCfg(cfg))
	return
}

// GetBlockProduction Returns recent block production information from the current or previous epoch.
func (sc *Client) GetBlockProduction(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlockProduction")
	return
}

// GetBlockTime commitment, array of u64 integers
func (sc *Client) GetBlockTime(ctx context.Context, blockNum uint64) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlockTime", blockNum)
	return
}

// GetBlocks returns array of u64 integers
func (sc *Client) GetBlocks(ctx context.Context, startSlot, endSlot uint64) (res []uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlocks", startSlot, endSlot)
	return
}

// GetBlocksWithLimit returns array of u64 integers
func (sc *Client) GetBlocksWithLimit(ctx context.Context, startSlot, limit uint64) (res []uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlocksWithLimit", startSlot, limit)
	return
}

// GetClusterNodes commitment, array of u64 integers
func (sc *Client) GetClusterNodes(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getClusterNodes")
	return
}

// GetEpochInfo commitment, array of u64 integers
func (sc *Client) GetEpochInfo(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getEpochInfo")
	return
}

// GetEpochSchedule commitment, array of u64 integers
func (sc *Client) GetEpochSchedule(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getEpochSchedule")
	return
}

// GetFeeForMessage Base-64 encoded Message
func (sc *Client) GetFeeForMessage(ctx context.Context, msg string) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getFeeForMessage", msg)
	return
}

// GetFirstAvailableBlock commitment, array of u64 integers
func (sc *Client) GetFirstAvailableBlock(ctx context.Context) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getFirstAvailableBlock")
	return
}

// GetGenesisHash commitment, array of u64 integers
func (sc *Client) GetGenesisHash(ctx context.Context) (res common.Hash, err error) {
	err = sc.c.CallContext(ctx, &res, "getGenesisHash")
	return
}

// GetHealth commitment, array of u64 integers
func (sc *Client) GetHealth(ctx context.Context) (res string, err error) {
	err = sc.c.CallContext(ctx, &res, "getHealth")
	return
}

// GetHighestSnapshotSlot Base-64 encoded Message
func (sc *Client) GetHighestSnapshotSlot(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getHighestSnapshotSlot")
	return
}

// GetIdentity commitment, array of u64 integers
func (sc *Client) GetIdentity(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getIdentity")
	return
}

// GetInflationGovernor commitment, array of u64 integers
func (sc *Client) GetInflationGovernor(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getInflationGovernor")
	return
}

// GetInflationRate Base-64 encoded Message
func (sc *Client) GetInflationRate(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getInflationRate")
	return
}

// GetInflationReward Base-64 encoded Message
func (sc *Client) GetInflationReward(ctx context.Context, accounts []common.Address) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getInflationReward", accounts)
	return
}

// GetLargestAccounts Base-64 encoded Message
func (sc *Client) GetLargestAccounts(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getLargestAccounts")
	return
}

// GetLatestBlockhash Returns the latest blockhash
func (sc *Client) GetLatestBlockhash(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getLatestBlockhash")
	return
}

// GetLeaderSchedule Returns the leader schedule for an epoch
func (sc *Client) GetLeaderSchedule(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getLeaderSchedule")
	return
}

// GetMaxRetransmitSlot Get the max slot seen from retransmit stage.
func (sc *Client) GetMaxRetransmitSlot(ctx context.Context) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getMaxRetransmitSlot")
	return
}

// GetMaxShredInsertSlot Get the max slot seen from after shred insert.
func (sc *Client) GetMaxShredInsertSlot(ctx context.Context) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getMaxShredInsertSlot")
	return
}

// GetMinimumBalanceForRentExemption Returns minimum balance required to make account rent exempt.
func (sc *Client) GetMinimumBalanceForRentExemption(ctx context.Context, accDataLen ...uint64) (res uint64, err error) {
	// accountDataLength > 0
	if len(accDataLen) == 1 {
		err = sc.c.CallContext(ctx, &res, "getMinimumBalanceForRentExemption", accDataLen[0])
	} else {
		err = sc.c.CallContext(ctx, &res, "getMinimumBalanceForRentExemption")
	}
	return
}

// GetMultipleAccounts Returns the account information for a list of Pubkeys.
func (sc *Client) GetMultipleAccounts(ctx context.Context, accounts []common.Address) (res map[string]interface{}, err error) {
	// require accounts len <= 100
	if len(accounts) > 100 {
		return res, errors.New("accounts maximum is 100)")
	}
	err = sc.c.CallContext(ctx, &res, "getMultipleAccounts", accounts)
	return
}

// GetProgramAccounts Returns all accounts owned by the provided program Pubkey
func (sc *Client) GetProgramAccounts(ctx context.Context, program common.Address) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getProgramAccounts", program)
	return
}

// GetRecentPerformanceSamples Returns a list of recent performance samples, in reverse slot order.
// Performance samples are taken every 60 seconds and include the number of transactions and slots that occur in a given time window.
func (sc *Client) GetRecentPerformanceSamples(ctx context.Context, limit ...uint64) (res map[string]interface{}, err error) {
	// has limit
	if len(limit) == 1 {
		err = sc.c.CallContext(ctx, &res, "getRecentPerformanceSamples", limit[0])
	} else {
		err = sc.c.CallContext(ctx, &res, "getRecentPerformanceSamples")
	}
	return
}

// GetRecentPrioritizationFees Returns a list of prioritization fees from recent blocks.
func (sc *Client) GetRecentPrioritizationFees(ctx context.Context, accounts []common.Address) (res map[string]interface{}, err error) {
	// require accounts len <= 100
	if len(accounts) > 128 {
		return res, errors.New("accounts maximum is 128)")
	}
	err = sc.c.CallContext(ctx, &res, "getRecentPrioritizationFees", accounts)
	return
}

// GetSignatureStatuses Returns the statuses of a list of signatures. Each signature must be a txid, the first signature of a transaction.
func (sc *Client) GetSignatureStatuses(ctx context.Context, signatures []common.Signature) (res map[string]interface{}, err error) {
	// require accounts len <= 100
	if len(signatures) > 256 {
		return res, errors.New("signatures maximum is 256)")
	}
	err = sc.c.CallContext(ctx, &res, "getSignatureStatuses", signatures)
	return
}

// GetSignaturesForAddress Returns signatures for confirmed transactions that include the given address in their accountKeys list.
// Returns signatures backwards in time from the provided signature or most recent confirmed block
func (sc *Client) GetSignaturesForAddress(ctx context.Context, account common.Address) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getSignaturesForAddress", account)
	return
}

// GetSlot Returns the slot that has reached the given or default commitment level
// https://solana.com/docs/rpc#configuring-state-commitment
func (sc *Client) GetSlot(ctx context.Context) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getSlot")
	return
}

// GetSlotLeader Returns the current slot leader
func (sc *Client) GetSlotLeader(ctx context.Context) (res common.Address, err error) {
	err = sc.c.CallContext(ctx, &res, "getSlotLeader")
	return
}

// GetSlotLeaders Returns the slot leaders for a given slot range
func (sc *Client) GetSlotLeaders(ctx context.Context, startSlot, limit uint64) (res []common.Address, err error) {
	// require limit < 5000
	if limit > 5000 {
		return res, errors.New("limit maximum is 5000)")
	}
	err = sc.c.CallContext(ctx, &res, "getSlotLeaders", startSlot, limit)
	return
}

// GetStakeActivation Returns epoch activation information for a stake account
func (sc *Client) GetStakeActivation(ctx context.Context, account common.Address) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getStakeActivation", account)
	return
}

// GetStakeMinimumDelegation Returns the stake minimum delegation, in lamports.
func (sc *Client) GetStakeMinimumDelegation(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getStakeMinimumDelegation")
	return
}

// GetSupply Returns information about the current supply.
func (sc *Client) GetSupply(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getSupply")
	return
}

// GetTokenAccountBalance Returns the token balance of an SPL Token account.
func (sc *Client) GetTokenAccountBalance(ctx context.Context, account common.Address) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenAccountBalance", account)
	return
}

// GetTokenAccountsByDelegate Returns all SPL Token accounts by approved Delegate.
func (sc *Client) GetTokenAccountsByDelegate(ctx context.Context, account common.Address) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenAccountsByDelegate", account)
	return
}

// GetTokenAccountsByOwner Returns all SPL Token accounts by token owner.
func (sc *Client) GetTokenAccountsByOwner(ctx context.Context, account common.Address) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenAccountsByOwner", account)
	return
}

// GetTokenLargestAccounts Returns the 20 largest accounts of a particular SPL Token type.
func (sc *Client) GetTokenLargestAccounts(ctx context.Context, account common.Address) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenLargestAccounts", account)
	return
}

// GetTokenSupply Returns the total supply of an SPL Token type.
func (sc *Client) GetTokenSupply(ctx context.Context, account common.Address) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenSupply", account)
	return
}

// GetTransaction Returns transaction details for a confirmed transaction
func (sc *Client) GetTransaction(ctx context.Context, account common.Address) (res types.BlockTransaction, err error) {
	err = sc.c.CallContext(ctx, &res, "getTransaction", account)
	return
}

// GetTransactionCount Returns the current Transaction count from the ledger
func (sc *Client) GetTransactionCount(ctx context.Context) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getTransactionCount")
	return
}

// GetVersion Returns the current Solana version running on the node
func (sc *Client) GetVersion(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getVersion")
	return
}

// GetVoteAccounts Returns the account info and associated stake for all the voting accounts in the current bank.
func (sc *Client) GetVoteAccounts(ctx context.Context) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "getVoteAccounts")
	return
}

// IsBlockhashValid Returns whether a blockhash is still valid or not
func (sc *Client) IsBlockhashValid(ctx context.Context, hash common.Hash) (res types.BlockTransaction, err error) {
	err = sc.c.CallContext(ctx, &res, "isBlockhashValid", hash)
	return
}

// MinimumLedgerSlot Returns the lowest slot that the node has information about in its ledger.
func (sc *Client) MinimumLedgerSlot(ctx context.Context) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "minimumLedgerSlot")
	return
}

// RequestAirdrop Requests an airdrop of lamports to a Pubkey
func (sc *Client) RequestAirdrop(ctx context.Context, address common.Address, lamports *big.Int) (res common.Signature, err error) {
	err = sc.c.CallContext(ctx, &res, "requestAirdrop", address, lamports)
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
func (sc *Client) SendTransaction(ctx context.Context, signedTx common.Base58Data) (res common.Signature, err error) {
	err = sc.c.CallContext(ctx, &res, "sendTransaction", signedTx)
	return
}

// SimulateTransaction Simulate sending a transaction
func (sc *Client) SimulateTransaction(ctx context.Context, signedTx common.Base58Data) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "simulateTransaction", signedTx)
	return
}
