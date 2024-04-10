package solclient

import (
	"context"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/types"
	"testing"
)

func TestClient_AccountSubscribe(t *testing.T) {
	var (
		c             = newClient()
		ctx           = context.Background()
		accountNotify = make(chan *types.SubAccountInfo)
	)
	account := common.Base58ToAddress("CM78CPUeXjn8o3yroDHxUtKsZZgoy4GPkPPXfouKNH12")
	//
	sub, err := c.AccountSubscribe(ctx, accountNotify, account)
	if err != nil {
		t.Error("AccountSubscribe Failed: %w", err)
	}
	// if error
	if err != nil {
		panic(fmt.Sprintf("EthSubscribe Failed: %s", err.Error()))
	}

	defer sub.Unsubscribe()
	// fmt.Println("Start BotClient Pointer:", bot)
	// handler the subscribed pending hash
	for {
		select {
		case err = <-sub.Err():
			panic(fmt.Sprintf("[SUBSCRIPTION] Fatal error: %s", err.Error()))
		// Code block is executed when a new tx hash is piped to the channel
		case accountInfo := <-accountNotify:
			// analyse transaction from hash by querying the client
			fmt.Println(accountInfo)
		}
	}
}
