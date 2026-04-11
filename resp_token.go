// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

// UiTokenAmount represents a token balance in both raw and UI-friendly form.
// Returned inline by several token-related RPCs.
type UiTokenAmount struct {
	// Address account
	Address *PublicKey `json:"address,omitempty"`

	// Raw amount of tokens as a string, ignoring decimals.
	Amount string `json:"amount"`

	// Number of decimals configured for token's mint.
	Decimals uint8 `json:"decimals"`

	// Token amount as a float, accounting for decimals.
	UiAmount float64 `json:"uiAmount"`

	// Token amount as a string, accounting for decimals.
	UiAmountString string `json:"uiAmountString"`
}

// TokenBalance is an entry in TransactionMeta.Pre/PostTokenBalances.
type TokenBalance struct {
	// Index of the account in which the token balance is provided for.
	AccountIndex uint16 `json:"accountIndex"`

	// Pubkey of the token's mint.
	Mint PublicKey `json:"mint"`

	// Pubkey of token balance's owner.
	Owner PublicKey `json:"owner"`

	// ProgramId
	ProgramId string `json:"programId"`

	UiTokenAmount UiTokenAmount `json:"uiTokenAmount"`
}

// TokenAccountWithCtx is the context-wrapped response shape of
// getTokenAccountBalance and getTokenSupply.
type TokenAccountWithCtx struct {
	Context ContextSlot   `json:"context"`
	UiToken UiTokenAmount `json:"value"`
}

// TokenAccount is an element of the getTokenAccountsByOwner / ByDelegate
// response and is also used as the payload of programSubscribe.
type TokenAccount struct {
	Account AccountInfo `json:"account"`
	Pubkey  PublicKey   `json:"pubkey,omitempty"`
}

// TokenAccountsWithCtx is the context-wrapped list of token accounts.
type TokenAccountsWithCtx struct {
	Context  ContextSlot    `json:"context"`
	Accounts []TokenAccount `json:"value"`
}

// TokenLargestHolders is the response shape of getTokenLargestAccounts.
type TokenLargestHolders struct {
	Context ContextSlot     `json:"context"`
	Holders []UiTokenAmount `json:"value"`
}
