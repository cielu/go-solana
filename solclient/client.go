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
func (sc *Client) GetAccountInfo(ctx context.Context, account common.Address, cfg ...types.RpcAccountInfoCfg) (res types.AccountInfoWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getAccountInfo", account, getRpcCfg(cfg))
	return
}

// GetBalance Returns the lamport balance of the account of provided Pubkey
func (sc *Client) GetBalance(ctx context.Context, account common.Address, cfg ...types.RpcCommitmentWithMinSlotCfg) (balance types.BalanceWithCtx, err error) {
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
func (sc *Client) GetBlockHeight(ctx context.Context, cfg ...types.RpcCommitmentWithMinSlotCfg) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlockHeight", getRpcCfg(cfg))
	return
}

// GetBlockProduction Returns recent block production information from the current or previous epoch.
func (sc *Client) GetBlockProduction(ctx context.Context, cfg ...types.RpcGetBlockProduction) (res types.BlockProductionWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlockProduction", getRpcCfg(cfg))
	return
}

// GetBlockTime Returns the estimated production time of a block.
func (sc *Client) GetBlockTime(ctx context.Context, blockNum uint64) (res int64, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlockTime", blockNum)
	return
}

// GetBlocks Returns a list of confirmed blocks between two slots
func (sc *Client) GetBlocks(ctx context.Context, startSlot uint64, args ...interface{}) (res []uint64, err error) {
	var (
		tmpSlot uint64
		endSlot *uint64
		cfg *types.RpcCommitmentCfg
	)
	for _, arg := range args {
		// set endSlot & cfg
		switch v := arg.(type) {
		case int:
			tmpSlot = uint64(v)
		case uint64:
			tmpSlot = v
		case types.RpcCommitmentCfg:
			cfg = &v
		default:
			return res, errors.New("invalid args. Require: [uint64|types.RpcCommitmentCfg]")
		}
	}
	// setTmpSlot
	if tmpSlot > startSlot && tmpSlot - startSlot < 500000 {
		endSlot = &tmpSlot
	}
	err = sc.c.CallContext(ctx, &res, "getBlocks", startSlot, endSlot, cfg)
	return
}

// GetBlocksWithLimit Returns a list of confirmed blocks starting at the given slot
func (sc *Client) GetBlocksWithLimit(ctx context.Context, startSlot, limit uint64, cfg ...types.RpcCommitmentCfg) (res []uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlocksWithLimit", startSlot, limit, getRpcCfg(cfg))
	return
}

// GetClusterNodes Returns information about all the nodes participating in the cluster
func (sc *Client) GetClusterNodes(ctx context.Context) (res []types.ClusterInformation, err error) {
	err = sc.c.CallContext(ctx, &res, "getClusterNodes")
	return
}

// GetEpochInfo Returns information about the current epoch
func (sc *Client) GetEpochInfo(ctx context.Context, cfg ...types.RpcCommitmentWithMinSlotCfg) (res types.EpochInformation, err error) {
	err = sc.c.CallContext(ctx, &res, "getEpochInfo", getRpcCfg(cfg))
	return
}

// GetEpochSchedule Returns information about the current epoch
func (sc *Client) GetEpochSchedule(ctx context.Context) (res types.EpochSchedule, err error) {
	err = sc.c.CallContext(ctx, &res, "getEpochSchedule")
	return
}

// GetFeeForMessage Get the fee the network will charge for a particular Message
func (sc *Client) GetFeeForMessage(ctx context.Context, msg string, cfg ...types.RpcCommitmentWithMinSlotCfg) (res types.U64ValueWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getFeeForMessage", msg, getRpcCfg(cfg))
	return
}

// GetFirstAvailableBlock Returns the slot of the lowest confirmed block that has not been purged from the ledger
func (sc *Client) GetFirstAvailableBlock(ctx context.Context) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getFirstAvailableBlock")
	return
}

// GetGenesisHash Returns the genesis hash
func (sc *Client) GetGenesisHash(ctx context.Context) (res common.Hash, err error) {
	err = sc.c.CallContext(ctx, &res, "getGenesisHash")
	return
}

// GetHealth Returns the current health of the node.
// A healthy node is one that is within HEALTH_CHECK_SLOT_DISTANCE slots of the latest cluster confirmed slot.
func (sc *Client) GetHealth(ctx context.Context) (res string, err error) {
	err = sc.c.CallContext(ctx, &res, "getHealth")
	return
}

// GetHighestSnapshotSlot Returns the highest slot information that the node has snapshots for.
// This will find the highest full snapshot slot, and the highest incremental snapshot slot based on the full snapshot slot, if there is one.
func (sc *Client) GetHighestSnapshotSlot(ctx context.Context) (res types.HighestSnapshotSlot, err error) {
	err = sc.c.CallContext(ctx, &res, "getHighestSnapshotSlot")
	return
}

// GetIdentity Returns the identity pubkey for the current node
func (sc *Client) GetIdentity(ctx context.Context) (res types.Identity, err error) {
	err = sc.c.CallContext(ctx, &res, "getIdentity")
	return
}

// GetInflationGovernor Returns the current inflation governor
func (sc *Client) GetInflationGovernor(ctx context.Context, cfg ...types.RpcCommitmentCfg) (res types.InflationGovernor, err error) {
	err = sc.c.CallContext(ctx, &res, "getInflationGovernor", getRpcCfg(cfg))
	return
}

// GetInflationRate Returns the specific inflation values for the current epoch
func (sc *Client) GetInflationRate(ctx context.Context) (res types.InflationRate, err error) {
	err = sc.c.CallContext(ctx, &res, "getInflationRate")
	return
}

// GetInflationReward Returns the inflation / staking reward for a list of addresses for an epoch
func (sc *Client) GetInflationReward(ctx context.Context, args ...interface{}) (res map[string]interface{}, err error) {
	var (
		accounts []common.Address
		cfg *types.RpcCommitmentCfg
	)
	for _, arg := range args {
		// set endSlot & cfg
		switch v := arg.(type) {
		case common.Address:
			accounts = append(accounts, v)
		case []common.Address:
			accounts = v
		case types.RpcCommitmentCfg:
			cfg = &v
		default:
			return res, errors.New("invalid args. Require: [common.Address|[]common.Address|types.RpcCommitmentCfg]")
		}
	}
	err = sc.c.CallContext(ctx, &res, "getInflationReward", accounts, cfg)
	return
}

// GetLargestAccounts Returns the 20 largest accounts, by lamport balance (results may be cached up to two hours)
func (sc *Client) GetLargestAccounts(ctx context.Context, cfg ...types.RpcCommitmentWithFilter) (res types.AccountWithLamport, err error) {
	err = sc.c.CallContext(ctx, &res, "getLargestAccounts", getRpcCfg(cfg))
	return
}

// GetLatestBlockhash Returns the latest blockhash
func (sc *Client) GetLatestBlockhash(ctx context.Context, cfg ...types.RpcCommitmentWithMinSlotCfg) (res types.LastBlockWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getLatestBlockhash", getRpcCfg(cfg))
	return
}

// GetLeaderSchedule Returns the leader schedule for an epoch --> args [slot:u64]
func (sc *Client) GetLeaderSchedule(ctx context.Context, args ...interface{}) (res map[string][]uint64, err error) {
	// Fetch the leader schedule for the epoch that corresponds to the provided slot.
	var (
		slot *uint64
		tmpSlot uint64
		cfg *types.RpcCommitmentWithIdentity
	)
	// args
	for _, arg := range args {
		// set endSlot & cfg
		switch v := arg.(type) {
		case types.RpcCommitmentWithIdentity:
			cfg = &v
		case int:
			tmpSlot = uint64(v)
		case uint64:
			tmpSlot = v
		default:
			return res, errors.New("invalid args. Require: [uint64|types.RpcCommitmentWithIdentityCtx]")
		}
	}
	// slot
	if tmpSlot > 0 {
		slot = &tmpSlot
	}
	err = sc.c.CallContext(ctx, &res, "getLeaderSchedule", slot, cfg)
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
func (sc *Client) GetMinimumBalanceForRentExemption(ctx context.Context, args ...interface{}) (res uint64, err error) {
	// the Account's data length
	var (
		accLen *uint64
		tmpLen uint64
		cfg *types.RpcCommitmentCfg
	)
	// args
	for _, arg := range args {
		// set endSlot & cfg
		switch v := arg.(type) {
		case types.RpcCommitmentCfg:
			cfg = &v
		case int:
			tmpLen = uint64(v)
		case uint64:
			tmpLen = v
		default:
			return res, errors.New("invalid args. Require: [uint64|types.RpcCommitmentCfg]")
		}
	}
	// slot
	if tmpLen > 0 {
		accLen = &tmpLen
	}
	err = sc.c.CallContext(ctx, &res, "getMinimumBalanceForRentExemption", accLen, cfg)
	return
}

// GetMultipleAccounts Returns the account information for a list of Pubkeys.
func (sc *Client) GetMultipleAccounts(ctx context.Context, accounts []common.Address, cfg ...types.RpcAccountInfoCfg) (res types.AccountsInfoWithCtx, err error) {
	// require accounts len <= 100
	if len(accounts) > 100 {
		return res, errors.New("accounts maximum is 100)")
	}
	err = sc.c.CallContext(ctx, &res, "getMultipleAccounts", accounts, getRpcCfg(cfg))
	return
}

// GetProgramAccounts Returns all accounts owned by the provided program Pubkey
func (sc *Client) GetProgramAccounts(ctx context.Context, program common.Address, cfg ...types.RpcCombinedCfg) (res []types.ProgramAccount, err error) {
	err = sc.c.CallContext(ctx, &res, "getProgramAccounts", program, getRpcCfg(cfg))
	return
}

// GetRecentPerformanceSamples Returns a list of recent performance samples, in reverse slot order.
// Performance samples are taken every 60 seconds and include the number of transactions and slots that occur in a given time window.
func (sc *Client) GetRecentPerformanceSamples(ctx context.Context, args ...uint64) (res []types.RpcPerfSample, err error) {
	// has limit
	var limit *uint64
	// has arg
	if len(args) > 0 && args[0] <= 720 {
		limit = &args[0]
	}
	err = sc.c.CallContext(ctx, &res, "getRecentPerformanceSamples", limit)
	return
}

// GetRecentPrioritizationFees Returns a list of prioritization fees from recent blocks.
func (sc *Client) GetRecentPrioritizationFees(ctx context.Context, args ...interface{}) (res []types.RpcPrioritizationFee, err error) {
	// var []common.Address
	var accounts []common.Address
	// args
	for _, arg := range args {
		// require accounts len <= 100
		switch v := arg.(type) {
		case common.Address:
			accounts = append(accounts, v)
		case []common.Address:
			accounts = v
		}
	}
	if len(accounts) > 128 {
		return res, errors.New("accounts maximum is 128)")
	}
	err = sc.c.CallContext(ctx, &res, "getRecentPrioritizationFees", accounts)
	return
}

// GetSignatureStatuses Returns the statuses of a list of signatures. Each signature must be a txid, the first signature of a transaction.
func (sc *Client) GetSignatureStatuses(ctx context.Context, signatures []common.Signature, cfg ...types.RpcSearchTxHistoryCfg) (res map[string]interface{}, err error) {
	// require accounts len <= 100
	if len(signatures) > 256 {
		return res, errors.New("signatures maximum is 256)")
	}
	err = sc.c.CallContext(ctx, &res, "getSignatureStatuses", signatures, getRpcCfg(cfg))
	return
}

// GetSignaturesForAddress Returns signatures for confirmed transactions that include the given address in their accountKeys list.
// Returns signatures backwards in time from the provided signature or most recent confirmed block
func (sc *Client) GetSignaturesForAddress(ctx context.Context, account common.Address, cfg ...types.RpcSignaturesForAddressCfg) (res []types.SignatureInfo, err error) {
	err = sc.c.CallContext(ctx, &res, "getSignaturesForAddress", account, getRpcCfg(cfg))
	return
}

// GetSlot Returns the slot that has reached the given or default commitment level
// https://solana.com/docs/rpc#configuring-state-commitment
func (sc *Client) GetSlot(ctx context.Context, cfg ...types.RpcCommitmentWithMinSlotCfg) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getSlot", getRpcCfg(cfg))
	return
}

// GetSlotLeader Returns the current slot leader
func (sc *Client) GetSlotLeader(ctx context.Context, cfg ...types.RpcCommitmentWithMinSlotCfg) (res common.Address, err error) {
	err = sc.c.CallContext(ctx, &res, "getSlotLeader", getRpcCfg(cfg))
	return
}

// GetSlotLeaders Returns the slot leaders for a given slot range
func (sc *Client) GetSlotLeaders(ctx context.Context, args ...uint64) (res []common.Address, err error) {
	// startSlot, limit uint64, cfg...
	// var []common.Address
	var startSlot, limit *uint64
	switch len(args) {
	case 1:
		startSlot = &args[0]
	case 2:
		startSlot = &args[0]
		limit = &args[1]
		// require limit < 5000
		if *limit > 5000 {
			err = errors.New("limit maximum is 5000)")
			return
		}
	default:
		err = errors.New("invalid args length")
		return
	}
	err = sc.c.CallContext(ctx, &res, "getSlotLeaders", startSlot, limit)
	return
}

// GetStakeActivation Returns epoch activation information for a stake account
func (sc *Client) GetStakeActivation(ctx context.Context, account common.Address, cfg ...types.RpcCommitmentWithMinSlotCfg) (res types.StakeActivation, err error) {
	err = sc.c.CallContext(ctx, &res, "getStakeActivation", account, getRpcCfg(cfg))
	return
}

// GetStakeMinimumDelegation Returns the stake minimum delegation, in lamports.
func (sc *Client) GetStakeMinimumDelegation(ctx context.Context, cfg ...types.RpcCommitmentCfg) (res types.U64ValueWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getStakeMinimumDelegation", getRpcCfg(cfg))
	return
}

// GetSupply Returns information about the current supply.
func (sc *Client) GetSupply(ctx context.Context, cfg ...types.RpcSupplyCfg) (res types.SupplyWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getSupply", getRpcCfg(cfg))
	return
}

// GetTokenAccountBalance Returns the token balance of an SPL Token account.
func (sc *Client) GetTokenAccountBalance(ctx context.Context, account common.Address, cfg ...types.RpcCommitmentCfg) (res types.TokenAccountWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenAccountBalance", account, getRpcCfg(cfg))
	return
}

// GetTokenAccountsByDelegate Returns all SPL Token accounts by approved Delegate.
func (sc *Client) GetTokenAccountsByDelegate(ctx context.Context, account common.Address, mintProg types.RpcMintWithProgramID, cfg ...types.RpcAccountInfoCfg) (res types.TokenAccountsWithCtx, err error) {
	// `params` should have at least 2 argument(s)
	err = sc.c.CallContext(ctx, &res, "getTokenAccountsByDelegate", account, mintProg, getRpcCfg(cfg))
	return
}

// GetTokenAccountsByOwner Returns all SPL Token accounts by token owner.
func (sc *Client) GetTokenAccountsByOwner(ctx context.Context, account common.Address, mintProg types.RpcMintWithProgramID, cfg ...types.RpcAccountInfoCfg) (res types.TokenAccountsWithCtx, err error) {
	// use base64
	tmpCfg := getRpcCfg(cfg)
	// isNull
	if tmpCfg == nil {
		tmpCfg = &types.RpcAccountInfoCfg{}
	}
	// set encoding to base64
	if tmpCfg.Encoding == "" {
		tmpCfg.Encoding = types.EncodingBase64
	}
	err = sc.c.CallContext(ctx, &res, "getTokenAccountsByOwner", account, mintProg, tmpCfg)
	return
}

// GetTokenLargestAccounts Returns the 20 largest accounts of a particular SPL Token type.
func (sc *Client) GetTokenLargestAccounts(ctx context.Context, account common.Address, cfg ...types.RpcCommitmentCfg) (res types.TokenLargestHolders, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenLargestAccounts", account, getRpcCfg(cfg))
	return
}

// GetTokenSupply Returns the total supply of an SPL Token type.
func (sc *Client) GetTokenSupply(ctx context.Context, account common.Address, cfg ...types.RpcCommitmentCfg) (res types.TokenAccountWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenSupply", account, getRpcCfg(cfg))
	return
}

// GetTransaction Returns transaction details for a confirmed transaction
func (sc *Client) GetTransaction(ctx context.Context, account common.Signature, cfg ...types.RpcGetTransactionCfg) (res *types.BlockTransaction, err error) {
	err = sc.c.CallContext(ctx, &res, "getTransaction", account, getRpcCfg(cfg))
	return
}

// GetTransactionCount Returns the current Transaction count from the ledger
func (sc *Client) GetTransactionCount(ctx context.Context, cfg ...types.RpcCommitmentWithMinSlotCfg) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getTransactionCount", getRpcCfg(cfg))
	return
}

// GetVersion Returns the current Solana version running on the node
func (sc *Client) GetVersion(ctx context.Context) (res types.SolVersion, err error) {
	err = sc.c.CallContext(ctx, &res, "getVersion")
	return
}

// GetVoteAccounts Returns the account info and associated stake for all the voting accounts in the current bank.
func (sc *Client) GetVoteAccounts(ctx context.Context, cfg ...types.RpcVoteAccountCfg) (res types.RpcVoteAccounts, err error) {
	err = sc.c.CallContext(ctx, &res, "getVoteAccounts", getRpcCfg(cfg))
	return
}

// IsBlockHashValid Returns whether a blockHash is still valid or not
func (sc *Client) IsBlockHashValid(ctx context.Context, hash common.Hash, cfg ...types.RpcCommitmentWithMinSlotCfg) (res bool, err error) {
	err = sc.c.CallContext(ctx, &res, "isBlockhashValid", hash, getRpcCfg(cfg))
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
func (sc *Client) SendTransaction(ctx context.Context, signedTx common.SolData, cfg ...types.RpcSendTxCfg) (res common.Signature, err error) {
	err = sc.c.CallContext(ctx, &res, "sendTransaction", signedTx, getRpcCfg(cfg))
	return
}

// SimulateTransaction Simulate sending a transaction
func (sc *Client) SimulateTransaction(ctx context.Context, signedTx common.SolData) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "simulateTransaction", signedTx)
	return
}
