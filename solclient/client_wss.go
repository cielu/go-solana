// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solclient

import (
	"context"
	"errors"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/core"
	"github.com/cielu/go-solana/types"
)

// AccountSubscribe Subscribe to an account to receive notifications when the lamports or data for a given account public key changes
func (sc *Client) AccountSubscribe(ctx context.Context, ch chan<- types.AccountInfoWithCtx, account common.Address, cfg ...types.RpcCommitmentWithEncodingCfg) (core.Subscription, error) {
	// SolSubscribe
	sub, err := sc.c.Subscribe(ctx, "account", ch, account, getRpcCfg(cfg))
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// BlockSubscribe Subscribe to receive notification anytime a new block is confirmed or finalized.
func (sc *Client) BlockSubscribe(ctx context.Context, ch chan<- types.BlockInfoWithCtx, filter interface{}, cfg ...types.RpcGetBlockContextCfg) (core.Subscription, error) {
	// SolSubscribe
	switch filter.(type) {
	case string:
		filter = "all"
	case types.MentionsAccountProgramCfg:
	default:
		return nil, errors.New("invalid filter arg. Require: [string|types.MentionsAccountProgramCfg]")
	}
	sub, err := sc.c.Subscribe(ctx, "block", ch, filter, getRpcCfg(cfg))
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// LogsSubscribe Subscribe to transaction logging
func (sc *Client) LogsSubscribe(ctx context.Context, ch chan<- types.LogsInfoWithCtx, mentions interface{}, cfg ...types.RpcCommitmentCfg) (core.Subscription, error) {
	// SolSubscribe
	switch mentions.(type) {
	case string:
		// invalid mentions
		if mentions != "all" && mentions != "allWithVotes" {
			mentions = "all"
		}
	case types.MentionsCfg:
	default:
		return nil, errors.New("invalid mentions. Require: [string|types.MentionsCfg]")
	}
	sub, err := sc.c.Subscribe(ctx, "logs", ch, mentions, getRpcCfg(cfg))
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// ProgramSubscribe to a program to receive notifications when the lamports or data for an account owned by the given program changes
func (sc *Client) ProgramSubscribe(ctx context.Context, ch chan<- types.ProgramInfoWithCtx, address common.Address, cfg ...types.RpcCommitmentCfg) (core.Subscription, error) {
	// SolSubscribe
	sub, err := sc.c.Subscribe(ctx, "program", ch, address, getRpcCfg(cfg))
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// SignatureSubscribe Subscribe to receive a notification when the transaction with the given signature reaches the specified commitment level.
func (sc *Client) SignatureSubscribe(ctx context.Context, ch chan<- types.SignatureInfoWithCtx, signature common.Signature, cfg ...types.RpcCommitmentCfg) (core.Subscription, error) {
	// SolSubscribe
	sub, err := sc.c.Subscribe(ctx, "signature", ch, signature, getRpcCfg(cfg))
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// SlotSubscribe Subscribe to receive notification anytime a slot is processed by the validator
func (sc *Client) SlotSubscribe(ctx context.Context, ch chan<- types.SlotInfoWithCtx) (core.Subscription, error) {
	// SolSubscribe
	sub, err := sc.c.Subscribe(ctx, "slot", ch)
	if err != nil {
		return nil, err
	}
	return sub, nil
}
