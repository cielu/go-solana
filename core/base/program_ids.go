// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package base

import (
	"github.com/cielu/go-solana"
)

var (
	// SystemProgramID Create new accounts, allocate account data, assign accounts to owning programs,
	// transfer lamports from System Program owned accounts and pay transacation fees.
	SystemProgramID = solana.StrToPublicKey("11111111111111111111111111111111")

	// ConfigProgramID Add configuration data to the chain and the list of public keys that are permitted to modify it.
	ConfigProgramID = solana.StrToPublicKey("Config1111111111111111111111111111111111111")

	// StakeProgramID Create and manage accounts representing stake and rewards for delegations to validators.
	StakeProgramID = solana.StrToPublicKey("Stake11111111111111111111111111111111111111")

	// VoteProgramID Create and manage accounts that track validator voting state and rewards.
	VoteProgramID = solana.StrToPublicKey("Vote111111111111111111111111111111111111111")

	BPFLoaderDeprecatedProgramID = solana.StrToPublicKey("BPFLoader1111111111111111111111111111111111")
	// BPFLoaderProgramID Deploys, upgrades, and executes programs on the chain.
	BPFLoaderProgramID            = solana.StrToPublicKey("BPFLoader2111111111111111111111111111111111")
	BPFLoaderUpgradeableProgramID = solana.StrToPublicKey("BPFLoaderUpgradeab1e11111111111111111111111")

	// Secp256k1ProgramID Verify secp256k1 public key recovery operations (ecrecover).
	Secp256k1ProgramID = solana.StrToPublicKey("KeccakSecp256k11111111111111111111111111111")

	FeatureProgramID = solana.StrToPublicKey("Feature111111111111111111111111111111111111")

	ComputeBudget = solana.StrToPublicKey("ComputeBudget111111111111111111111111111111")

	AssetExecutorProgramID = solana.StrToPublicKey("J7Dai94nSeunCgErhYTRfWkssbhLFUeZsiymX4S6DNrL")

	//
	SPLNameServiceProgramID     = solana.StrToPublicKey("namesLPneVptA9Z5rqUDD9tMTWEJwofgaYwp8cawRkX")
	MetaplexTokenMetaProgramID  = solana.StrToPublicKey("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
	ComputeBudgetProgramID      = solana.StrToPublicKey("ComputeBudget111111111111111111111111111111")
	AddressLookupTableProgramID = solana.StrToPublicKey("AddressLookupTab1e1111111111111111111111111")
)

// SPL:
var (
	// TokenProgramID A Token program on the Solana blockchain.
	// This program defines a common implementation for Fungible and Non Fungible tokens.
	TokenProgramID = solana.StrToPublicKey("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

	Token2022ProgramID = solana.StrToPublicKey("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb")

	// TokenSwapProgramID A Uniswap-like exchange for the Token program on the Solana blockchain,
	// implementing multiple automated market maker (AMM) curves.
	TokenSwapProgramID = solana.StrToPublicKey("SwaPpA9LAaLfeLi3a68M4DjnLqgtticKg6CnyNwgAC8")
	TokenSwapFeeOwner  = solana.StrToPublicKey("HfoTxFR1Tm6kGmWgYWD6J7YHVy1UwqSULUGVLXkJqaKN")

	// TokenLendingProgramID A lending protocol for the Token program on the Solana blockchain inspired by Aave and Compound.
	TokenLendingProgramID = solana.StrToPublicKey("LendZqTs8gn5CTSJU1jWKhKuVpjJGom45nnwPb2AMTi")

	// SPLAssociatedTokenAccountProgramID This program defines the convention and provides the mechanism for mapping
	// the user's wallet address to the associated token accounts they hold.
	SPLAssociatedTokenAccountProgramID = solana.StrToPublicKey("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")

	// MemoProgramID The Memo program is a simple program that validates a string of UTF-8 encoded characters
	// and verifies that any accounts provided are signers of the transaction.
	// The program also logs the memo, as well as any verified signer addresses,
	// to the transaction log, so that anyone can easily observe memos
	// and know they were approved by zero or more addresses
	// by inspecting the transaction log from a trusted provider.
	MemoProgramID = solana.StrToPublicKey("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")
)

var (
	// SolMint The Mint for native SOL Token accounts
	SolMint = solana.StrToPublicKey("So11111111111111111111111111111111111111112")

	SolMint2022 = solana.StrToPublicKey("9pan9bMn5HatX4EJdBwg9VgCa7Uz5HL8N1m5D3NdXejP")

	WrappedSol = SolMint

	TokenMetadataProgramID = solana.StrToPublicKey("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

	MetaplexCandyMachineV2ProgramID = solana.StrToPublicKey("cndy3Z4yapfJBmL3ShUp5exZKqR3z33thTzeNMm2gRZ")
	MetaplexTokenMetadataProgramID  = solana.StrToPublicKey("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
)
