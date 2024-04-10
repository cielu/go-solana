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
		accountNotify = make(chan types.AccountInfoWithCtx)
	)

	account := common.Base58ToAddress("6v3nv8BUJKpXvnBnD4ZvpDiG3u847ALLYyo1NACn2zmV")
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

func TestClient_BlockSubscribe(t *testing.T) {
	var (
		c               = newClient()
		ctx             = context.Background()
		blockInfoNotify = make(chan types.BlockInfoWithCtx)
		filter          = ""
	)

	//
	sub, err := c.BlockSubscribe(ctx, blockInfoNotify, filter)

	if err != nil {
		t.Error("BlockSubscribe Failed: %w", err)
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
		case blockInfo := <-blockInfoNotify:
			// analyse transaction from hash by querying the client
			fmt.Printf("%#v\n", blockInfo)
			fmt.Println("Slot", blockInfo.Context.Slot)
			fmt.Println("Err", blockInfo.BlockInfo.Err)
			fmt.Println("Height", blockInfo.BlockInfo.BlockHeight)
		}
	}
}

func TestClient_LogsSubscribe(t *testing.T) {
	var (
		c              = newClient()
		ctx            = context.Background()
		logsInfoNotify = make(chan types.LogsInfoWithCtx)
		filter         = "all"
	)

	//
	sub, err := c.LogsSubscribe(ctx, logsInfoNotify, filter)
	if err != nil {
		t.Error("LogsSubscribe Failed: %w", err)
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
		case accountInfo := <-logsInfoNotify:
			// analyse transaction from hash by querying the client
			fmt.Println(accountInfo)
		}
	}
}

func TestClient_ProgramSubscribe(t *testing.T) {
	var (
		c             = newClient()
		ctx           = context.Background()
		programNotify = make(chan types.ProgramInfoWithCtx)
	)
	address := common.Base58ToAddress("3p7U58GR11MnfRuWCBufj9AW3Y7P1x848CWgtECpNQpt")
	sub, err := c.ProgramSubscribe(ctx, programNotify, address)
	if err != nil {
		t.Error("ProgramSubscribe Failed: %w", err)
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
		case accountInfo := <-programNotify:
			// analyse transaction from hash by querying the client
			fmt.Println(accountInfo)
		}
	}
}
