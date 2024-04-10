// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.


package solclient

import (
	"context"
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



