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
func (sc *Client) AccountSubscribe(ctx context.Context, ch chan<- *types.SubAccountInfo, account common.Address, cfg ...types.RpcCommitmentWithEncodingCfg) (core.Subscription, error) {
	// SolSubscribe
	sub, err := sc.c.Subscribe(ctx, "account", ch, account, getRpcCfg(cfg))
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// BlockSubscribe Subscribe to receive notification anytime a new block is confirmed or finalized.
func (sc *Client) BlockSubscribe(ctx context.Context, ch chan<- *types.BlockInfoNotify, filter interface{}, cfg ...types.RpcGetBlockContextCfg) (core.Subscription, error) {
	// SolSubscribe
	switch filter.(type) {
	case string:
		filter = "all"
	case types.MentionsAccountProgramParam:
	default:
		return nil, errors.New("invalid args. Require: [string|types.MentionsAccountProgramParam]")
	}
	sub, err := sc.c.Subscribe(ctx, "block", ch, filter, getRpcCfg(cfg))
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// LogsSubscribe Subscribe to transaction logging
func (sc *Client) LogsSubscribe(ctx context.Context, ch chan<- *types.LogsInfoNotify, mentions interface{}, cfg ...types.RpcCommitmentCfg) (core.Subscription, error) {
	// SolSubscribe
	switch mentions.(type) {
	case string:
		if mentions != "all" && mentions != "allWithVotes" {
			mentions = "all"
		}
	case types.MentionsParam:
	default:
		return nil, errors.New("invalid args. Require: [string|types.MentionsParam]")
	}
	sub, err := sc.c.Subscribe(ctx, "logs", ch, mentions, getRpcCfg(cfg))
	if err != nil {
		return nil, err
	}
	return sub, nil
}
