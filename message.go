package solana

import (
	"encoding/base64"
	"fmt"

	"github.com/cielu/go-solana/library"
	"github.com/cielu/go-solana/pkg/encodbin"
)

type MessageVersion int

const (
	MessageVersionLegacy MessageVersion = 0 // default
	MessageVersionV0     MessageVersion = 1 // v0
)

type MessageAddressTableLookupSlice []MessageAddressTableLookup

// NumLookups returns the number of accounts in all the MessageAddressTableLookupSlice
func (lookups MessageAddressTableLookupSlice) NumLookups() int {
	count := 0
	for _, lookup := range lookups {
		// TODO: check if this is correct.
		count += len(lookup.ReadonlyIndexes)
		count += len(lookup.WritableIndexes)
	}
	return count
}

// NumWritableLookups returns the number of writable accounts
// across all the lookups (all the address tables).
func (lookups MessageAddressTableLookupSlice) NumWritableLookups() int {
	count := 0
	for _, lookup := range lookups {
		count += len(lookup.WritableIndexes)
	}
	return count
}

// GetTableIDs returns the list of all address table IDs.
func (lookups MessageAddressTableLookupSlice) GetTableIDs() []PublicKey {
	if lookups == nil {
		return nil
	}
	ids := make([]PublicKey, 0)
	// lookups
	for _, lookup := range lookups {
		// append unique key
		ids = library.UniqueAppend(ids, lookup.AccountKey)
	}
	return ids
}

type MessageAddressTableLookup struct {
	AccountKey      PublicKey `json:"accountKey"` // The account key of the address table.
	WritableIndexes []uint8   `json:"writableIndexes"`
	ReadonlyIndexes []uint8   `json:"readonlyIndexes"`
}

type Message struct {
	version MessageVersion
	// List of base-58 encoded public keys used by the transaction,
	// including by the instructions and for signatures.
	// The first `message.header.numRequiredSignatures` public keys must sign the transaction.
	AccountKeys []PublicKey `json:"accountKeys"`
	// Details the account types and signatures required by the transaction.
	Header MessageHeader `json:"header"`
	// A base-58 encoded hash of a recent block in the ledger used to
	// prevent transaction duplication and to give transactions lifetimes.
	RecentBlockhash Hash `json:"recentBlockhash"`
	// List of program instructions that will be executed in sequence
	// and committed in one atomic transaction if all succeed.
	Instructions []CompiledInstruction `json:"instructions"`
	// List of address table lookups used to load additional accounts for this transaction.
	AddressTableLookups MessageAddressTableLookupSlice `json:"addressTableLookups"`
	// The actual tables that contain the list of account pubkeys.
	// NOTE: you need to fetch these from the chain, and then call `SetAddressTables`
	// before you use this transaction -- otherwise, you will get a panic.
	addressTables map[PublicKey][]PublicKey
	// if true, the lookups have been resolved, and the `AccountKeys` slice contains all the accounts (static + dynamic).
	resolved bool
}

// SetAddressTables sets the actual address tables used by this message.
// Use `mx.GetAddressTableLookups().GetTableIDs()` to get the list of all address table IDs.
// NOTE: you can call this once.
func (m *Message) SetAddressTables(tables map[PublicKey][]PublicKey) error {
	if m.addressTables != nil {
		return fmt.Errorf("address tables already set")
	}
	m.addressTables = tables
	return nil
}

// GetAddressTables returns the actual address tables used by this message.
// NOTE: you must have called `SetAddressTable` before being able to use this method.
func (m *Message) GetAddressTables() map[PublicKey][]PublicKey {
	return m.addressTables
}

// GetProgram current program address
func (m Message) GetProgram(idIndex uint16) PublicKey {
	return m.AccountKeys[idIndex]
}

// GetVersion returns the message version.
func (m *Message) GetVersion() MessageVersion {
	return m.version
}

// SetAddressTableLookups (re)sets the lookups used by this message.
func (m *Message) SetAddressTableLookups(lookups []MessageAddressTableLookup) *Message {
	m.AddressTableLookups = lookups
	m.version = MessageVersionV0
	return m
}

// AddAddressTableLookup adds a new lookup to the message.
func (m *Message) AddAddressTableLookup(lookup MessageAddressTableLookup) *Message {
	m.AddressTableLookups = append(m.AddressTableLookups, lookup)
	m.version = MessageVersionV0
	return m
}

func (m *Message) MarshalBinary() ([]byte, error) {
	// VersionV0
	if m.version == MessageVersionV0 {
		return m.MarshalV0()
	}
	// Default is legacy
	return m.MarshalLegacy()
}

func (m *Message) MarshalLegacy() ([]byte, error) {

	buf := []byte{
		m.Header.NumRequiredSignatures,
		m.Header.NumReadonlySignedAccounts,
		m.Header.NumReadonlyUnsignedAccounts,
	}

	encodbin.EncodeCompactU16Length(&buf, len(m.AccountKeys))
	for _, key := range m.AccountKeys {
		buf = append(buf, key[:]...)
	}

	buf = append(buf, m.RecentBlockhash[:]...)

	encodbin.EncodeCompactU16Length(&buf, len(m.Instructions))
	for _, instruction := range m.Instructions {
		buf = append(buf, byte(instruction.ProgramIDIndex))
		encodbin.EncodeCompactU16Length(&buf, len(instruction.Accounts))
		for _, accountIdx := range instruction.Accounts {
			buf = append(buf, byte(accountIdx))
		}

		encodbin.EncodeCompactU16Length(&buf, len(instruction.Data))

		buf = append(buf, instruction.Data...)
	}
	return buf, nil
}

func (m Message) getStaticKeys() (keys []PublicKey) {
	if m.resolved {
		// If the message has been resolved, then the account keys have already
		// been appended to the `AccountKeys` field of the message.
		return m.AccountKeys[:m.numStaticAccounts()]
	}
	return m.AccountKeys
}

func (m *Message) MarshalV0() ([]byte, error) {
	buf := []byte{
		m.Header.NumRequiredSignatures,
		m.Header.NumReadonlySignedAccounts,
		m.Header.NumReadonlyUnsignedAccounts,
	}
	{
		// Encode only the keys that are not in the address table lookups.
		staticAccountKeys := m.getStaticKeys()
		encodbin.EncodeCompactU16Length(&buf, len(staticAccountKeys))
		for _, key := range staticAccountKeys {
			buf = append(buf, key[:]...)
		}

		buf = append(buf, m.RecentBlockhash[:]...)

		encodbin.EncodeCompactU16Length(&buf, len(m.Instructions))
		for _, instruction := range m.Instructions {
			buf = append(buf, byte(instruction.ProgramIDIndex))
			encodbin.EncodeCompactU16Length(&buf, len(instruction.Accounts))
			for _, accountIdx := range instruction.Accounts {
				buf = append(buf, byte(accountIdx))
			}

			encodbin.EncodeCompactU16Length(&buf, len(instruction.Data))
			buf = append(buf, instruction.Data...)
		}
	}
	versionNum := byte(m.version) // TODO: what number is this?
	if versionNum > 127 {
		return nil, fmt.Errorf("invalid message version: %d", m.version)
	}
	buf = append([]byte{byte(versionNum + 127)}, buf...)

	// wite length of address table lookups as u8
	buf = append(buf, byte(len(m.AddressTableLookups)))
	for _, lookup := range m.AddressTableLookups {
		// write account pubkey
		buf = append(buf, lookup.AccountKey[:]...)
		// write writable indexes
		encodbin.EncodeCompactU16Length(&buf, len(lookup.WritableIndexes))
		buf = append(buf, lookup.WritableIndexes...)
		// write readonly indexes
		encodbin.EncodeCompactU16Length(&buf, len(lookup.ReadonlyIndexes))
		buf = append(buf, lookup.ReadonlyIndexes...)
	}

	return buf, nil
}

// numStaticAccounts returns the number of accounts that are always present in the
// account keys list (i.e. all the accounts that are NOT in the lookup table).
func (m Message) numStaticAccounts() int {
	if !m.resolved || m.AddressTableLookups == nil {
		return len(m.AccountKeys)
	}
	// NumLookups
	return len(m.AccountKeys) - m.AddressTableLookups.NumLookups()
}

// Signers returns the pubkeys of all accounts that are signers.
func (m *Message) Signers() []PublicKey {
	out := make([]PublicKey, 0, len(m.AccountKeys))
	for _, a := range m.AccountKeys {
		if m.IsSigner(a) {
			out = append(out, a)
		}
	}
	return out
}

// Writable returns the pubkeys of all accounts that are writable.
func (m *Message) Writable() (out []PublicKey) {
	for _, a := range m.AccountKeys {
		if m.IsWritable(a) {
			out = append(out, a)
		}
	}
	return out
}

func (m *Message) IsSigner(account PublicKey) bool {
	for idx, acc := range m.AccountKeys {
		if acc == account {
			return idx < int(m.Header.NumRequiredSignatures)
		}
	}
	return false
}

func (m *Message) IsWritable(account PublicKey) bool {
	index := 0
	found := false
	for idx, acc := range m.AccountKeys {
		if acc == account {
			found = true
			index = idx
		}
	}
	if !found {
		return false
	}
	h := m.Header
	return (index < int(h.NumRequiredSignatures-h.NumReadonlySignedAccounts)) ||
		((index >= int(h.NumRequiredSignatures)) && (index < len(m.AccountKeys)-int(h.NumReadonlyUnsignedAccounts)))
}

func (m *Message) signerKeys() []PublicKey {
	return m.AccountKeys[0:m.Header.NumRequiredSignatures]
}

type MessageHeader struct {
	// The total number of signatures required to make the transaction valid.
	// The signatures must match the first `numRequiredSignatures` of `message.account_keys`.
	NumRequiredSignatures uint8 `json:"numRequiredSignatures"`

	// The last numReadonlySignedAccounts of the signed keys are read-only accounts.
	// Programs may process multiple transactions that load read-only accounts within
	// a single PoH entry, but are not permitted to credit or debit lamports or modify
	// account data.
	// Transactions targeting the same read-write account are evaluated sequentially.
	NumReadonlySignedAccounts uint8 `json:"numReadonlySignedAccounts"`

	// The last `numReadonlyUnsignedAccounts` of the unsigned keys are read-only accounts.
	NumReadonlyUnsignedAccounts uint8 `json:"numReadonlyUnsignedAccounts"`
}

func (m *Message) UnmarshalWithDecoder(decoder *encodbin.Decoder) (err error) {
	// peek first byte to determine if this is a legacy or v0 message
	versionNum, err := decoder.Peek(1)
	if err != nil {
		return err
	}
	// TODO: is this the right way to determine if this is a legacy or v0 message?
	if versionNum[0] < 127 {
		m.version = MessageVersionLegacy
	} else {
		m.version = MessageVersionV0
	}
	switch m.version {
	case MessageVersionV0:
		return m.UnmarshalV0(decoder)
	case MessageVersionLegacy:
		return m.UnmarshalLegacy(decoder)
	default:
		return fmt.Errorf("invalid message version: %d", m.version)
	}
}

func (m *Message) UnmarshalBase64(b64 string) error {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return m.UnmarshalWithDecoder(encodbin.NewBinDecoder(b))
}

func (m *Message) UnmarshalV0(decoder *encodbin.Decoder) (err error) {
	version, err := decoder.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read message version: %w", err)
	}
	// TODO: check version
	m.version = MessageVersion(version - 127)

	// The middle of the message is the same as the legacy message:
	err = m.UnmarshalLegacy(decoder)
	if err != nil {
		return err
	}

	// Read address table lookups length:
	addressTableLookupsLen, err := decoder.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read address table lookups length: %w", err)
	}
	if addressTableLookupsLen > 0 {
		m.AddressTableLookups = make([]MessageAddressTableLookup, addressTableLookupsLen)
		for i := 0; i < int(addressTableLookupsLen); i++ {
			// read account pubkey
			_, err = decoder.Read(m.AddressTableLookups[i].AccountKey[:])
			if err != nil {
				return fmt.Errorf("failed to read account pubkey: %w", err)
			}

			// read writable indexes
			writableIndexesLen, err := decoder.ReadCompactU16Length()
			if err != nil {
				return fmt.Errorf("failed to read writable indexes length: %w", err)
			}
			m.AddressTableLookups[i].WritableIndexes = make([]byte, writableIndexesLen)
			_, err = decoder.Read(m.AddressTableLookups[i].WritableIndexes)
			if err != nil {
				return fmt.Errorf("failed to read writable indexes: %w", err)
			}

			// read readonly indexes
			readonlyIndexesLen, err := decoder.ReadCompactU16Length()
			if err != nil {
				return fmt.Errorf("failed to read readonly indexes length: %w", err)
			}
			m.AddressTableLookups[i].ReadonlyIndexes = make([]byte, readonlyIndexesLen)
			_, err = decoder.Read(m.AddressTableLookups[i].ReadonlyIndexes)
			if err != nil {
				return fmt.Errorf("failed to read readonly indexes: %w", err)
			}
		}
	}
	return nil
}

func (m *Message) UnmarshalLegacy(decoder *encodbin.Decoder) (err error) {
	{
		m.Header.NumRequiredSignatures, err = decoder.ReadUint8()
		if err != nil {
			return fmt.Errorf("unable to decode m.Header.NumRequiredSignatures: %w", err)
		}
		m.Header.NumReadonlySignedAccounts, err = decoder.ReadUint8()
		if err != nil {
			return fmt.Errorf("unable to decode m.Header.NumReadonlySignedAccounts: %w", err)
		}
		m.Header.NumReadonlyUnsignedAccounts, err = decoder.ReadUint8()
		if err != nil {
			return fmt.Errorf("unable to decode m.Header.NumReadonlyUnsignedAccounts: %w", err)
		}
	}
	{
		numAccountKeys, err := decoder.ReadCompactU16()
		if err != nil {
			return fmt.Errorf("unable to decode numAccountKeys: %w", err)
		}
		m.AccountKeys = make([]PublicKey, numAccountKeys)
		for i := 0; i < numAccountKeys; i++ {
			_, err := decoder.Read(m.AccountKeys[i][:])
			if err != nil {
				return fmt.Errorf("unable to decode m.AccountKeys[%d]: %w", i, err)
			}
		}
	}
	{
		_, err := decoder.Read(m.RecentBlockhash[:])
		if err != nil {
			return fmt.Errorf("unable to decode m.RecentBlockhash: %w", err)
		}
	}
	{
		numInstructions, err := decoder.ReadCompactU16()
		if err != nil {
			return fmt.Errorf("unable to decode numInstructions: %w", err)
		}
		m.Instructions = make([]CompiledInstruction, numInstructions)
		for instructionIndex := 0; instructionIndex < numInstructions; instructionIndex++ {
			programIDIndex, err := decoder.ReadUint8()
			if err != nil {
				return fmt.Errorf("unable to decode m.Instructions[%d].ProgramIDIndex: %w", instructionIndex, err)
			}
			m.Instructions[instructionIndex].ProgramIDIndex = uint16(programIDIndex)

			{
				numAccounts, err := decoder.ReadCompactU16()
				if err != nil {
					return fmt.Errorf("unable to decode numAccounts for ix[%d]: %w", instructionIndex, err)
				}
				m.Instructions[instructionIndex].Accounts = make([]uint16, numAccounts)
				for i := 0; i < numAccounts; i++ {
					accountIndex, err := decoder.ReadUint8()
					if err != nil {
						return fmt.Errorf("unable to decode accountIndex for ix[%d].Accounts[%d]: %w", instructionIndex, i, err)
					}
					m.Instructions[instructionIndex].Accounts[i] = uint16(accountIndex)
				}
			}
			{
				dataLen, err := decoder.ReadCompactU16()
				if err != nil {
					return fmt.Errorf("unable to decode dataLen for ix[%d]: %w", instructionIndex, err)
				}
				dataBytes, err := decoder.ReadNBytes(dataLen)
				if err != nil {
					return fmt.Errorf("unable to decode dataBytes for ix[%d]: %w", instructionIndex, err)
				}
				// setBytes
				m.Instructions[instructionIndex].Data = dataBytes
			}
		}
	}

	return nil
}
