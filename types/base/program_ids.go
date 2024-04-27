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

import "github.com/cielu/go-solana/common"

var (
	// Create new accounts, allocate account data, assign accounts to owning programs,
	// transfer lamports from System Program owned accounts and pay transacation fees.
	SystemProgramID = common.StrToAddress("11111111111111111111111111111111")

	// Add configuration data to the chain and the list of public keys that are permitted to modify it.
	ConfigProgramID = common.StrToAddress("Config1111111111111111111111111111111111111")

	// Create and manage accounts representing stake and rewards for delegations to validators.
	StakeProgramID = common.StrToAddress("Stake11111111111111111111111111111111111111")

	// Create and manage accounts that track validator voting state and rewards.
	VoteProgramID = common.StrToAddress("Vote111111111111111111111111111111111111111")

	BPFLoaderDeprecatedProgramID = common.StrToAddress("BPFLoader1111111111111111111111111111111111")
	// Deploys, upgrades, and executes programs on the chain.
	BPFLoaderProgramID            = common.StrToAddress("BPFLoader2111111111111111111111111111111111")
	BPFLoaderUpgradeableProgramID = common.StrToAddress("BPFLoaderUpgradeab1e11111111111111111111111")

	// Verify secp256k1 public key recovery operations (ecrecover).
	Secp256k1ProgramID = common.StrToAddress("KeccakSecp256k11111111111111111111111111111")

	FeatureProgramID = common.StrToAddress("Feature111111111111111111111111111111111111")

	ComputeBudget = common.StrToAddress("ComputeBudget111111111111111111111111111111")
)

// SPL:
var (
	// A Token program on the Solana blockchain.
	// This program defines a common implementation for Fungible and Non Fungible tokens.
	TokenProgramID = common.StrToAddress("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

	Token2022ProgramID = common.StrToAddress("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb")

	// A Uniswap-like exchange for the Token program on the Solana blockchain,
	// implementing multiple automated market maker (AMM) curves.
	TokenSwapProgramID = common.StrToAddress("SwaPpA9LAaLfeLi3a68M4DjnLqgtticKg6CnyNwgAC8")
	TokenSwapFeeOwner  = common.StrToAddress("HfoTxFR1Tm6kGmWgYWD6J7YHVy1UwqSULUGVLXkJqaKN")

	// A lending protocol for the Token program on the Solana blockchain inspired by Aave and Compound.
	TokenLendingProgramID = common.StrToAddress("LendZqTs8gn5CTSJU1jWKhKuVpjJGom45nnwPb2AMTi")

	// This program defines the convention and provides the mechanism for mapping
	// the user's wallet address to the associated token accounts they hold.
	SPLAssociatedTokenAccountProgramID = common.StrToAddress("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")

	// The Memo program is a simple program that validates a string of UTF-8 encoded characters
	// and verifies that any accounts provided are signers of the transaction.
	// The program also logs the memo, as well as any verified signer addresses,
	// to the transaction log, so that anyone can easily observe memos
	// and know they were approved by zero or more addresses
	// by inspecting the transaction log from a trusted provider.
	MemoProgramID = common.StrToAddress("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")
)

var (
	// The Mint for native SOL Token accounts
	SolMint = common.StrToAddress("So11111111111111111111111111111111111111112")

	SolMint2022 = common.StrToAddress("9pan9bMn5HatX4EJdBwg9VgCa7Uz5HL8N1m5D3NdXejP")

	WrappedSol = SolMint
)

var (
	TokenMetadataProgramID = common.StrToAddress("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
)

var (
	// The Clock sysvar contains data on cluster time,
	// including the current slot, epoch, and estimated wall-clock Unix timestamp.
	// It is updated every slot.
	SysVarClockPubkey = common.StrToAddress("SysvarC1ock11111111111111111111111111111111")

	// The EpochSchedule sysvar contains epoch scheduling constants that are set in genesis,
	// and enables calculating the number of slots in a given epoch,
	// the epoch for a given slot, etc.
	// (Note: the epoch schedule is distinct from the leader schedule)
	SysVarEpochSchedulePubkey = common.StrToAddress("SysvarEpochSchedu1e111111111111111111111111")

	// The Fees sysvar contains the fee calculator for the current slot.
	// It is updated every slot, based on the fee-rate governor.
	SysVarFeesPubkey = common.StrToAddress("SysvarFees111111111111111111111111111111111")

	// The Instructions sysvar contains the serialized instructions in a Message while that Message is being processed.
	// This allows program instructions to reference other instructions in the same transaction.
	SysVarInstructionsPubkey = common.StrToAddress("Sysvar1nstructions1111111111111111111111111")

	// The RecentBlockhashes sysvar contains the active recent blockhashes as well as their associated fee calculators.
	// It is updated every slot.
	// Entries are ordered by descending block height,
	// so the first entry holds the most recent block hash,
	// and the last entry holds an old block hash.
	SysVarRecentBlockHashesPubkey = common.StrToAddress("SysvarRecentB1ockHashes11111111111111111111")

	// The Rent sysvar contains the rental rate.
	// Currently, the rate is static and set in genesis.
	// The Rent burn percentage is modified by manual feature activation.
	SysVarRentPubkey = common.StrToAddress("SysvarRent111111111111111111111111111111111")

	//
	SysVarRewardsPubkey = common.StrToAddress("SysvarRewards111111111111111111111111111111")

	// The SlotHashes sysvar contains the most recent hashes of the slot's parent banks.
	// It is updated every slot.
	SysVarSlotHashesPubkey = common.StrToAddress("SysvarS1otHashes111111111111111111111111111")

	// The SlotHistory sysvar contains a bitvector of slots present over the last epoch. It is updated every slot.
	SysVarSlotHistoryPubkey = common.StrToAddress("SysvarS1otHistory11111111111111111111111111")

	// The StakeHistory sysvar contains the history of cluster-wide stake activations and de-activations per epoch.
	// It is updated at the start of every epoch.
	SysVarStakeHistoryPubkey = common.StrToAddress("SysvarStakeHistory1111111111111111111111111")
)

var (
	MetaplexCandyMachineV2ProgramID = common.StrToAddress("cndy3Z4yapfJBmL3ShUp5exZKqR3z33thTzeNMm2gRZ")
	MetaplexTokenMetadataProgramID  = common.StrToAddress("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
)
