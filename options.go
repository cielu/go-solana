// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

// This file declares functional-option types for RPC methods whose original
// signatures accepted a variadic `args ...interface{}` list. The old form was
// order-sensitive, type-unsafe, and silently dropped unknown arguments. The
// options below replace that with compile-time type safety, IDE discoverability,
// and self-documenting call sites.
//
// Naming convention: every option function is prefixed with the method it
// applies to (WithBlocks*, WithReward*, WithSchedule*, WithRentExempt*,
// WithFees*). This avoids collisions between the several methods that all
// accept a commitment level.

package solana

// -----------------------------------------------------------------------------
// GetBlocks
// -----------------------------------------------------------------------------

// BlocksOption configures a GetBlocks call.
type BlocksOption func(*blocksConfig)

type blocksConfig struct {
	endSlot *uint64
	cfg     *RpcCommitmentCfg
}

// WithBlocksEndSlot sets the inclusive end slot for GetBlocks. The resulting
// range [startSlot, endSlot] must be at most 500000 slots wide.
func WithBlocksEndSlot(endSlot uint64) BlocksOption {
	return func(c *blocksConfig) { c.endSlot = &endSlot }
}

// WithBlocksCommitment sets the commitment level for GetBlocks.
func WithBlocksCommitment(commitment EnumRpcCommitment) BlocksOption {
	return func(c *blocksConfig) {
		c.cfg = &RpcCommitmentCfg{Commitment: commitment}
	}
}

// -----------------------------------------------------------------------------
// GetInflationReward
// -----------------------------------------------------------------------------

// RewardOption configures a GetInflationReward call.
type RewardOption func(*rewardConfig)

type rewardConfig struct {
	accounts []PublicKey
	cfg      *RpcCommitmentCfg
}

// WithRewardAccounts appends one or more accounts to the inflation-reward
// query. It may be called multiple times.
func WithRewardAccounts(accounts ...PublicKey) RewardOption {
	return func(c *rewardConfig) { c.accounts = append(c.accounts, accounts...) }
}

// WithRewardCommitment sets the commitment level for GetInflationReward.
func WithRewardCommitment(commitment EnumRpcCommitment) RewardOption {
	return func(c *rewardConfig) {
		c.cfg = &RpcCommitmentCfg{Commitment: commitment}
	}
}

// -----------------------------------------------------------------------------
// GetLeaderSchedule
// -----------------------------------------------------------------------------

// ScheduleOption configures a GetLeaderSchedule call.
type ScheduleOption func(*scheduleConfig)

type scheduleConfig struct {
	slot *uint64
	cfg  *RpcCommitmentWithIdentity
}

// WithScheduleSlot fetches the leader schedule for the epoch containing this slot.
func WithScheduleSlot(slot uint64) ScheduleOption {
	return func(c *scheduleConfig) { c.slot = &slot }
}

// WithScheduleCommitment sets the commitment level for GetLeaderSchedule.
func WithScheduleCommitment(commitment EnumRpcCommitment) ScheduleOption {
	return func(c *scheduleConfig) {
		if c.cfg == nil {
			c.cfg = &RpcCommitmentWithIdentity{}
		}
		c.cfg.Commitment = commitment
	}
}

// WithScheduleIdentity filters the leader schedule to the given validator identity.
func WithScheduleIdentity(identity PublicKey) ScheduleOption {
	return func(c *scheduleConfig) {
		if c.cfg == nil {
			c.cfg = &RpcCommitmentWithIdentity{}
		}
		c.cfg.Identity = &identity
	}
}

// -----------------------------------------------------------------------------
// GetMinimumBalanceForRentExemption
// -----------------------------------------------------------------------------

// RentExemptOption configures a GetMinimumBalanceForRentExemption call.
type RentExemptOption func(*rentExemptConfig)

type rentExemptConfig struct {
	dataLen *uint64
	cfg     *RpcCommitmentCfg
}

// WithRentExemptDataLen sets the account data length used for the rent-exemption
// calculation. Passing 0 is equivalent to omitting the option.
func WithRentExemptDataLen(dataLen uint64) RentExemptOption {
	return func(c *rentExemptConfig) { c.dataLen = &dataLen }
}

// WithRentExemptCommitment sets the commitment level.
func WithRentExemptCommitment(commitment EnumRpcCommitment) RentExemptOption {
	return func(c *rentExemptConfig) {
		c.cfg = &RpcCommitmentCfg{Commitment: commitment}
	}
}

// -----------------------------------------------------------------------------
// GetRecentPrioritizationFees
// -----------------------------------------------------------------------------

// FeesOption configures a GetRecentPrioritizationFees call.
type FeesOption func(*feesConfig)

type feesConfig struct {
	accounts []PublicKey
}

// WithFeesAccounts appends one or more accounts to scope the prioritization-fee
// query. It may be called multiple times. At most 128 accounts total are allowed
// by the RPC node.
func WithFeesAccounts(accounts ...PublicKey) FeesOption {
	return func(c *feesConfig) { c.accounts = append(c.accounts, accounts...) }
}
