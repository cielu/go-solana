// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

import (
	"context"
	"errors"
	"math/big"

	"github.com/cielu/go-solana/rpc"
)

// Client defines typed wrappers for the solana RPC API.
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
	err = sc.c.CallContext(ctx, &res, "getAccountInfo", account, getRpcCfg(cfg))
	return
}

// GetBalance Returns the lamport balance of the account of provided Pubkey
func (sc *Client) GetBalance(ctx context.Context, account PublicKey, cfg ...RpcCommitmentWithMinSlotCfg) (balance BalanceWithCtx, err error) {
	err = sc.c.CallContext(ctx, &balance, "getBalance", account, getRpcCfg(cfg))
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
	err = sc.c.CallContext(ctx, &blockInfo, "getBlock", blockNum, c)
	return
}

// GetBlockCommitment Returns commitment for particular block
func (sc *Client) GetBlockCommitment(ctx context.Context, blockNum uint64) (blockCmt BlockCommitment, err error) {
	err = sc.c.CallContext(ctx, &blockCmt, "getBlockCommitment", blockNum)
	return
}

// GetBlockHeight Returns the current block height of the node
func (sc *Client) GetBlockHeight(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlockHeight", getRpcCfg(cfg))
	return
}

// GetBlockProduction Returns recent block production information from the current or previous epoch.
func (sc *Client) GetBlockProduction(ctx context.Context, cfg ...RpcGetBlockProduction) (res BlockProductionWithCtx, err error) {
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
		cfg     *RpcCommitmentCfg
	)
	for _, arg := range args {
		// set endSlot & cfg
		switch v := arg.(type) {
		case int:
			tmpSlot = uint64(v)
		case uint64:
			tmpSlot = v
		case RpcCommitmentCfg:
			cfg = &v
		default:
			return res, errors.New("invalid args. Require: [uint64|types.RpcCommitmentCfg]")
		}
	}
	// setTmpSlot
	if tmpSlot > startSlot && tmpSlot-startSlot < 500000 {
		endSlot = &tmpSlot
	}
	err = sc.c.CallContext(ctx, &res, "getBlocks", startSlot, endSlot, cfg)
	return
}

// GetBlocksWithLimit Returns a list of confirmed blocks starting at the given slot
func (sc *Client) GetBlocksWithLimit(ctx context.Context, startSlot, limit uint64, cfg ...RpcCommitmentCfg) (res []uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getBlocksWithLimit", startSlot, limit, getRpcCfg(cfg))
	return
}

// GetClusterNodes Returns information about all the nodes participating in the cluster
func (sc *Client) GetClusterNodes(ctx context.Context) (res []ClusterInformation, err error) {
	err = sc.c.CallContext(ctx, &res, "getClusterNodes")
	return
}

// GetEpochInfo Returns information about the current epoch
func (sc *Client) GetEpochInfo(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res EpochInformation, err error) {
	err = sc.c.CallContext(ctx, &res, "getEpochInfo", getRpcCfg(cfg))
	return
}

// GetEpochSchedule Returns information about the current epoch
func (sc *Client) GetEpochSchedule(ctx context.Context) (res EpochSchedule, err error) {
	err = sc.c.CallContext(ctx, &res, "getEpochSchedule")
	return
}

// GetFeeForMessage Get the fee the network will charge for a particular Message
func (sc *Client) GetFeeForMessage(ctx context.Context, msg string, cfg ...RpcCommitmentWithMinSlotCfg) (res U64ValueWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getFeeForMessage", msg, getRpcCfg(cfg))
	return
}

// GetFirstAvailableBlock Returns the slot of the lowest confirmed block that has not been purged from the ledger
func (sc *Client) GetFirstAvailableBlock(ctx context.Context) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getFirstAvailableBlock")
	return
}

// GetGenesisHash Returns the genesis hash
func (sc *Client) GetGenesisHash(ctx context.Context) (res Hash, err error) {
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
func (sc *Client) GetHighestSnapshotSlot(ctx context.Context) (res HighestSnapshotSlot, err error) {
	err = sc.c.CallContext(ctx, &res, "getHighestSnapshotSlot")
	return
}

// GetIdentity Returns the identity pubkey for the current node
func (sc *Client) GetIdentity(ctx context.Context) (res Identity, err error) {
	err = sc.c.CallContext(ctx, &res, "getIdentity")
	return
}

// GetInflationGovernor Returns the current inflation governor
func (sc *Client) GetInflationGovernor(ctx context.Context, cfg ...RpcCommitmentCfg) (res InflationGovernor, err error) {
	err = sc.c.CallContext(ctx, &res, "getInflationGovernor", getRpcCfg(cfg))
	return
}

// GetInflationRate Returns the specific inflation values for the current epoch
func (sc *Client) GetInflationRate(ctx context.Context) (res InflationRate, err error) {
	err = sc.c.CallContext(ctx, &res, "getInflationRate")
	return
}

// GetInflationReward Returns the inflation / staking reward for a list of PublicKeyes for an epoch
func (sc *Client) GetInflationReward(ctx context.Context, args ...interface{}) (res []InflationReward, err error) {
	var (
		accounts []PublicKey
		cfg      *RpcCommitmentCfg
	)
	for _, arg := range args {
		// set endSlot & cfg
		switch v := arg.(type) {
		case PublicKey:
			accounts = append(accounts, v)
		case []PublicKey:
			accounts = v
		case RpcCommitmentCfg:
			cfg = &v
		default:
			return res, errors.New("invalid args. Require: [PublicKey|[]PublicKey|RpcCommitmentCfg]")
		}
	}
	err = sc.c.CallContext(ctx, &res, "getInflationReward", accounts, cfg)
	return
}

// GetLargestAccounts Returns the 20 largest accounts, by lamport balance (results may be cached up to two hours)
func (sc *Client) GetLargestAccounts(ctx context.Context, cfg ...RpcCommitmentWithFilter) (res AccountWithLamport, err error) {
	err = sc.c.CallContext(ctx, &res, "getLargestAccounts", getRpcCfg(cfg))
	return
}

// GetLatestBlockhash Returns the latest blockhash
func (sc *Client) GetLatestBlockhash(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res LastBlockWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getLatestBlockhash", getRpcCfg(cfg))
	return
}

// GetLeaderSchedule Returns the leader schedule for an epoch --> args [slot:u64]
func (sc *Client) GetLeaderSchedule(ctx context.Context, args ...interface{}) (res map[string][]uint64, err error) {
	// Fetch the leader schedule for the epoch that corresponds to the provided slot.
	var (
		slot    *uint64
		tmpSlot uint64
		cfg     *RpcCommitmentWithIdentity
	)
	// args
	for _, arg := range args {
		// set endSlot & cfg
		switch v := arg.(type) {
		case RpcCommitmentWithIdentity:
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
		cfg    *RpcCommitmentCfg
	)
	// args
	for _, arg := range args {
		// set endSlot & cfg
		switch v := arg.(type) {
		case RpcCommitmentCfg:
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
func (sc *Client) GetMultipleAccounts(ctx context.Context, accounts []PublicKey, cfg ...RpcAccountInfoCfg) (res AccountsInfoWithCtx, err error) {
	// require accounts len <= 100
	if len(accounts) > 100 {
		return res, errors.New("accounts maximum is 100)")
	}
	err = sc.c.CallContext(ctx, &res, "getMultipleAccounts", accounts, getRpcCfg(cfg))
	return
}

// GetProgramAccounts Returns all accounts owned by the provided program Pubkey
func (sc *Client) GetProgramAccounts(ctx context.Context, program PublicKey, cfg ...RpcCombinedCfg) (res []ProgramAccount, err error) {
	err = sc.c.CallContext(ctx, &res, "getProgramAccounts", program, getRpcCfg(cfg))
	return
}

// GetRecentPerformanceSamples Returns a list of recent performance samples, in reverse slot order.
// Performance samples are taken every 60 seconds and include the number of transactions and slots that occur in a given time window.
func (sc *Client) GetRecentPerformanceSamples(ctx context.Context, args ...uint64) (res []RpcPerfSample, err error) {
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
func (sc *Client) GetRecentPrioritizationFees(ctx context.Context, args ...interface{}) (res []RpcPrioritizationFee, err error) {
	// var []common.PublicKey
	var accounts []PublicKey
	// args
	for _, arg := range args {
		// require accounts len <= 100
		switch v := arg.(type) {
		case PublicKey:
			accounts = append(accounts, v)
		case []PublicKey:
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
func (sc *Client) GetSignatureStatuses(ctx context.Context, signatures []Signature, cfg ...RpcSearchTxHistoryCfg) (res SignatureStatusWithCtx, err error) {
	// require signatures len <= 256
	if len(signatures) > 256 {
		return res, errors.New("signatures maximum is 256)")
	}
	err = sc.c.CallContext(ctx, &res, "getSignatureStatuses", signatures, getRpcCfg(cfg))
	return
}

// GetSignaturesForAddress Returns signatures for confirmed transactions that include the given PublicKey in their accountKeys list.
// Returns signatures backwards in time from the provided signature or most recent confirmed block
func (sc *Client) GetSignaturesForAddress(ctx context.Context, account PublicKey, cfg ...RpcSignaturesForAddressCfg) (res []SignatureInfo, err error) {
	err = sc.c.CallContext(ctx, &res, "getSignaturesForAddress", account, getRpcCfg(cfg))
	return
}

// GetSlot Returns the slot that has reached the given or default commitment level
// https://solana.com/docs/rpc#configuring-state-commitment
func (sc *Client) GetSlot(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getSlot", getRpcCfg(cfg))
	return
}

// GetSlotLeader Returns the current slot leader
func (sc *Client) GetSlotLeader(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res PublicKey, err error) {
	err = sc.c.CallContext(ctx, &res, "getSlotLeader", getRpcCfg(cfg))
	return
}

// GetSlotLeaders Returns the slot leaders for a given slot range
func (sc *Client) GetSlotLeaders(ctx context.Context, args ...uint64) (res []PublicKey, err error) {
	// startSlot, limit uint64, cfg...
	// var []common.PublicKey
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
func (sc *Client) GetStakeActivation(ctx context.Context, account PublicKey, cfg ...RpcCommitmentWithMinSlotCfg) (res StakeActivation, err error) {
	err = sc.c.CallContext(ctx, &res, "getStakeActivation", account, getRpcCfg(cfg))
	return
}

// GetStakeMinimumDelegation Returns the stake minimum delegation, in lamports.
func (sc *Client) GetStakeMinimumDelegation(ctx context.Context, cfg ...RpcCommitmentCfg) (res U64ValueWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getStakeMinimumDelegation", getRpcCfg(cfg))
	return
}

// GetSupply Returns information about the current supply.
func (sc *Client) GetSupply(ctx context.Context, cfg ...RpcSupplyCfg) (res SupplyWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getSupply", getRpcCfg(cfg))
	return
}

// GetTokenAccountBalance Returns the token balance of an SPL Token account.
func (sc *Client) GetTokenAccountBalance(ctx context.Context, account PublicKey, cfg ...RpcCommitmentCfg) (res TokenAccountWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenAccountBalance", account, getRpcCfg(cfg))
	return
}

// GetTokenAccountsByDelegate Returns all SPL Token accounts by approved Delegate.
func (sc *Client) GetTokenAccountsByDelegate(ctx context.Context, delegate PublicKey, mintProg RpcMintWithProgramID, cfg ...RpcAccountInfoCfg) (res TokenAccountsWithCtx, err error) {
	// `params` should have at least 2 argument(s)
	err = sc.c.CallContext(ctx, &res, "getTokenAccountsByDelegate", delegate, mintProg, getRpcCfg(cfg))
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
	err = sc.c.CallContext(ctx, &res, "getTokenAccountsByOwner", owner, program, tmpCfg)
	return
}

// GetTokenLargestAccounts Returns the 20 largest accounts of a particular SPL Token type.
func (sc *Client) GetTokenLargestAccounts(ctx context.Context, splToken PublicKey, cfg ...RpcCommitmentCfg) (res TokenLargestHolders, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenLargestAccounts", splToken, getRpcCfg(cfg))
	return
}

// GetTokenSupply Returns the total supply of an SPL Token type.
func (sc *Client) GetTokenSupply(ctx context.Context, splToken PublicKey, cfg ...RpcCommitmentCfg) (res TokenAccountWithCtx, err error) {
	err = sc.c.CallContext(ctx, &res, "getTokenSupply", splToken, getRpcCfg(cfg))
	return
}

// GetTransaction Returns transaction details for a confirmed transaction
func (sc *Client) GetTransaction(ctx context.Context, signature Signature, cfg ...RpcGetTransactionCfg) (res TransactionInfo, err error) {
	err = sc.c.CallContext(ctx, &res, "getTransaction", signature, getRpcCfg(cfg))
	return
}

// GetTransactionCount Returns the current Transaction count from the ledger
func (sc *Client) GetTransactionCount(ctx context.Context, cfg ...RpcCommitmentWithMinSlotCfg) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "getTransactionCount", getRpcCfg(cfg))
	return
}

// GetVersion Returns the current Solana version running on the node
func (sc *Client) GetVersion(ctx context.Context) (res SolVersion, err error) {
	err = sc.c.CallContext(ctx, &res, "getVersion")
	return
}

// GetVoteAccounts Returns the account info and associated stake for all the voting accounts in the current bank.
func (sc *Client) GetVoteAccounts(ctx context.Context, cfg ...RpcVoteAccountCfg) (res RpcVoteAccounts, err error) {
	err = sc.c.CallContext(ctx, &res, "getVoteAccounts", getRpcCfg(cfg))
	return
}

// IsBlockHashValid Returns whether a blockHash is still valid or not
func (sc *Client) IsBlockHashValid(ctx context.Context, hash Hash, cfg ...RpcCommitmentWithMinSlotCfg) (res bool, err error) {
	err = sc.c.CallContext(ctx, &res, "isBlockhashValid", hash, getRpcCfg(cfg))
	return
}

// MinimumLedgerSlot Returns the lowest slot that the node has information about in its ledger.
func (sc *Client) MinimumLedgerSlot(ctx context.Context) (res uint64, err error) {
	err = sc.c.CallContext(ctx, &res, "minimumLedgerSlot")
	return
}

// RequestAirdrop Requests an airdrop of lamports to a Pubkey
func (sc *Client) RequestAirdrop(ctx context.Context, address PublicKey, lamport *big.Int) (res Signature, err error) {
	err = sc.c.CallContext(ctx, &res, "requestAirdrop", address, lamport)
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
func (sc *Client) SendTransaction(ctx context.Context, signedTx []byte, cfg ...RpcSendTxCfg) (res Signature, err error) {
	err = sc.c.CallContext(ctx, &res, "sendTransaction", BytesToSolData(signedTx), getRpcCfg(cfg))
	return
}

// SimulateTransaction Simulate sending a transaction
func (sc *Client) SimulateTransaction(ctx context.Context, signedTx []byte) (res map[string]interface{}, err error) {
	err = sc.c.CallContext(ctx, &res, "simulateTransaction", BytesToSolData(signedTx))
	return
}
