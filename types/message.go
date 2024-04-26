package types

import (
	"bytes"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/bincode"
	"sort"
	"strconv"
)

const (
	MessageVersionLegacy = "legacy"
	MessageVersionV0     = "v0"
)

type MessageVersion string

type CompiledAddressLookupTable struct {
	AccountKey      common.Address
	WritableIndexes []uint8
	ReadonlyIndexes []uint8
}

type MessageBody struct {
	Version             MessageVersion
	Header              MessageHeader
	Accounts            []common.Address
	RecentBlockHash     common.Hash
	Instructions        []CompiledInstruction
	AddressLookupTables []CompiledAddressLookupTable
}

type AddressLookupTableAccount struct {
	Key       common.Address
	Addresses []common.Address
}

type CompiledKeys struct {
	Payer      *common.Address
	KeyMetaMap map[common.Address]CompiledKeyMeta
}

type CompiledKeyMeta struct {
	IsSigner   bool
	IsWritable bool
	IsInvoked  bool
}

func NewCompiledKeys(instructions []Instruction, payer *common.Address) CompiledKeys {
	m := map[common.Address]CompiledKeyMeta{}

	for _, instruction := range instructions {
		// compile program
		v := m[instruction.ProgramID]
		v.IsInvoked = true
		m[instruction.ProgramID] = v

		// compile accounts
		for i := 0; i < len(instruction.Accounts); i++ {
			account := instruction.Accounts[i]

			v := m[account.PubKey]
			v.IsSigner = v.IsSigner || account.IsSigner
			v.IsWritable = v.IsWritable || account.IsWritable
			m[account.PubKey] = v
		}
	}

	if payer != nil && *payer != (common.Address{}) {
		v := m[*payer]
		v.IsSigner = true
		v.IsWritable = true
		m[*payer] = v
	}

	return CompiledKeys{
		Payer:      payer,
		KeyMetaMap: m,
	}
}

func NewMessage(param NewMessageParam) MessageBody {
	writableSignedAccount := []common.Address{}
	readOnlySignedAccount := []common.Address{}
	writableUnsignedAccount := []common.Address{}
	readOnlyUnsignedAccount := []common.Address{}

	addressLookupTableAccountCount := len(param.AddressLookupTableAccounts)
	addressLookupTableWritable := make([][]common.Address, addressLookupTableAccountCount)
	addressLookupTableWritableIdx := make([][]uint8, addressLookupTableAccountCount)
	addressLookupTableReadonly := make([][]common.Address, addressLookupTableAccountCount)
	addressLookupTableReadonlyIdx := make([][]uint8, addressLookupTableAccountCount)

	addressLookupTableMaps := make([]map[common.Address]uint8, 0, addressLookupTableAccountCount)
	for _, addressLookupTableAccount := range param.AddressLookupTableAccounts {
		m := map[common.Address]uint8{}
		for i, address := range addressLookupTableAccount.Addresses {
			m[address] = uint8(i)
		}
		addressLookupTableMaps = append(addressLookupTableMaps, m)
	}

	compiledKeys := NewCompiledKeys(param.Instructions, &param.FeePayer)
	allKeys := make([]common.Address, 0, len(compiledKeys.KeyMetaMap))
	for key := range compiledKeys.KeyMetaMap {
		allKeys = append(allKeys, key)
	}
	sort.Slice(allKeys, func(i, j int) bool {
		return bytes.Compare(allKeys[i].Bytes(), allKeys[j].Bytes()) < 0
	})

NEXT_ACCOUNT:
	for _, key := range allKeys {
		meta := compiledKeys.KeyMetaMap[key]
		if key == param.FeePayer {
			continue NEXT_ACCOUNT
		}
		if meta.IsSigner {
			if meta.IsWritable {
				writableSignedAccount = append(writableSignedAccount, key)
			} else {
				readOnlySignedAccount = append(readOnlySignedAccount, key)
			}
		} else {
			if meta.IsWritable {
				for n, addressLookupTableMap := range addressLookupTableMaps {
					idx, exist := addressLookupTableMap[key]
					if exist && !meta.IsInvoked {
						addressLookupTableWritable[n] = append(addressLookupTableWritable[n], key)
						addressLookupTableWritableIdx[n] = append(addressLookupTableWritableIdx[n], idx)
						continue NEXT_ACCOUNT
					}
				}
				// if not found in address lookup table
				writableUnsignedAccount = append(writableUnsignedAccount, key)
			} else {
				for n, addressLookupTableMap := range addressLookupTableMaps {
					idx, exist := addressLookupTableMap[key]
					if exist && !meta.IsInvoked {
						addressLookupTableReadonly[n] = append(addressLookupTableReadonly[n], key)
						addressLookupTableReadonlyIdx[n] = append(addressLookupTableReadonlyIdx[n], idx)
						continue NEXT_ACCOUNT
					}
				}
				// if not found in address lookup table
				readOnlyUnsignedAccount = append(readOnlyUnsignedAccount, key)
			}
		}
	}

	// add fee payer
	writableSignedAccount = append([]common.Address{param.FeePayer}, writableSignedAccount...)

	l := 0 +
		len(writableSignedAccount) +
		len(readOnlySignedAccount) +
		len(writableUnsignedAccount) +
		len(readOnlyUnsignedAccount) +
		len(addressLookupTableWritable) +
		len(addressLookupTableReadonly)

	publicKeys := make([]common.Address, 0, l)
	publicKeys = append(publicKeys, writableSignedAccount...)
	publicKeys = append(publicKeys, readOnlySignedAccount...)
	publicKeys = append(publicKeys, writableUnsignedAccount...)
	publicKeys = append(publicKeys, readOnlyUnsignedAccount...)

	compiledAddressLookupTables := []CompiledAddressLookupTable{}
	lookupAddressCount := 0
	for i := 0; i < addressLookupTableAccountCount; i++ {
		publicKeys = append(publicKeys, addressLookupTableWritable[i]...)
		lookupAddressCount += len(addressLookupTableWritable[i])
	}
	for i := 0; i < addressLookupTableAccountCount; i++ {
		publicKeys = append(publicKeys, addressLookupTableReadonly[i]...)
		lookupAddressCount += len(addressLookupTableReadonly[i])

		if len(addressLookupTableWritable[i]) > 0 || len(addressLookupTableReadonly[i]) > 0 {
			compiledAddressLookupTables = append(compiledAddressLookupTables, CompiledAddressLookupTable{
				AccountKey:      param.AddressLookupTableAccounts[i].Key,
				WritableIndexes: addressLookupTableWritableIdx[i],
				ReadonlyIndexes: addressLookupTableReadonlyIdx[i],
			})
		}
	}

	var version MessageVersion = MessageVersionLegacy
	if addressLookupTableAccountCount > 0 {
		version = MessageVersionV0
	}

	publicKeyToIdx := map[common.Address]int{}
	for idx, publicKey := range publicKeys {
		publicKeyToIdx[publicKey] = idx
	}

	compiledInstructions := []CompiledInstruction{}
	for _, instruction := range param.Instructions {
		accountIdx := []uint16{}
		for _, account := range instruction.Accounts {
			accountIdx = append(accountIdx, uint16(publicKeyToIdx[account.PubKey]))
		}
		compiledInstructions = append(compiledInstructions, CompiledInstruction{
			ProgramIDIndex: uint16(publicKeyToIdx[instruction.ProgramID]),
			Accounts:       accountIdx,
			Data:           common.SolData{RawData: instruction.Data, Encoding: ""},
		})
	}

	return MessageBody{
		Version: version,
		Header: MessageHeader{
			NumRequiredSignatures:       uint8(len(writableSignedAccount) + len(readOnlySignedAccount)),
			NumReadonlySignedAccounts:   uint8(len(readOnlySignedAccount)),
			NumReadonlyUnsignedAccounts: uint8(len(readOnlyUnsignedAccount)),
		},
		Accounts:            publicKeys[:len(publicKeys)-lookupAddressCount],
		RecentBlockHash:     param.RecentBlockhash,
		Instructions:        compiledInstructions,
		AddressLookupTables: compiledAddressLookupTables,
	}
}

func (m *MessageBody) Serialize() ([]byte, error) {
	b := []byte{}

	b = append(b, m.Header.NumRequiredSignatures)
	b = append(b, m.Header.NumReadonlySignedAccounts)
	b = append(b, m.Header.NumReadonlyUnsignedAccounts)

	b = append(b, bincode.UintToVarLenBytes(uint64(len(m.Accounts)))...)
	for _, key := range m.Accounts {
		b = append(b, key[:]...)
	}

	blockHash := m.RecentBlockHash.Bytes()

	b = append(b, blockHash...)

	b = append(b, bincode.UintToVarLenBytes(uint64(len(m.Instructions)))...)
	for _, instruction := range m.Instructions {
		b = append(b, byte(instruction.ProgramIDIndex))
		b = append(b, bincode.UintToVarLenBytes(uint64(len(instruction.Accounts)))...)
		for _, accountIdx := range instruction.Accounts {
			b = append(b, byte(accountIdx))
		}

		b = append(b, bincode.UintToVarLenBytes(uint64(len(instruction.Data.RawData)))...)
		b = append(b, instruction.Data.RawData...)
	}

	if len(m.Version) > 0 && m.Version != MessageVersionLegacy {
		versionNum, err := strconv.Atoi(string(m.Version[1:]))
		if err != nil || versionNum > 255 {
			return nil, fmt.Errorf("failed to parse message version")
		}
		if versionNum > 128 {
			return nil, fmt.Errorf("unexpected message version")
		}
		b = append([]byte{byte(versionNum + 128)}, b...)

		validAddressLookupCount := 0
		accountLookupTableSerializedData := []byte{}

		if len(m.Accounts) > 0 {
			for _, addressLookupTable := range m.AddressLookupTables {
				if len(addressLookupTable.WritableIndexes) != 0 || len(addressLookupTable.ReadonlyIndexes) != 0 {
					accountLookupTableSerializedData = append(accountLookupTableSerializedData, addressLookupTable.AccountKey.Bytes()...)
					accountLookupTableSerializedData = append(accountLookupTableSerializedData, bincode.UintToVarLenBytes(uint64(len(addressLookupTable.WritableIndexes)))...)
					accountLookupTableSerializedData = append(accountLookupTableSerializedData, addressLookupTable.WritableIndexes...)
					accountLookupTableSerializedData = append(accountLookupTableSerializedData, bincode.UintToVarLenBytes(uint64(len(addressLookupTable.ReadonlyIndexes)))...)
					accountLookupTableSerializedData = append(accountLookupTableSerializedData, addressLookupTable.ReadonlyIndexes...)
					validAddressLookupCount++
				}
			}

		}

		b = append(b, bincode.UintToVarLenBytes(uint64(validAddressLookupCount))...)
		b = append(b, accountLookupTableSerializedData...)
	}

	return b, nil
}
