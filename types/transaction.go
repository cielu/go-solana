package types

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/core"
	"github.com/cielu/go-solana/crypto"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/mr-tron/base58"
	"sort"
)

type Transaction struct {
	// A list of base-58 encoded signatures applied to the transaction.
	// The list is always of length `message.header.numRequiredSignatures` and not empty.
	// The signature at index `i` corresponds to the public key at index
	// `i` in `message.account_keys`. The first one is used as the transaction id.
	Signatures []common.Signature `json:"signatures"`

	// Defines the content of the transaction.
	Message Message `json:"message"`
}

type CompiledInstruction struct {
	// StackHeight if empty
	StackHeight *uint16 `json:"stackHeight"`
	// Index into the message.accountKeys array indicating the program account that executes this instruction.
	// NOTE: it is actually a uint8, but using a uint16 because uint8 is treated as a byte everywhere,
	// and that can be an issue.
	ProgramIDIndex uint16 `json:"programIdIndex"`
	// List of ordered indices into the message.accountKeys array indicating which accounts to pass to the program.
	// NOTE: it is actually a []uint8, but using a uint16 because []uint8 is treated as a []byte everywhere,
	// and that can be an issue.
	Accounts []uint16 `json:"accounts"`
	// The program input data encoded in a base-58 string.
	Data common.SolData `json:"data"`
}

func NewTransaction(instructions []Instruction, recentBlockHash common.Hash, payer common.Address) (*Transaction, error) {
	if len(instructions) == 0 {
		return nil, fmt.Errorf("requires at-least one instruction to create a transaction")
	}

	feePayer := payer
	if feePayer.IsEmpty() {
		found := false
		for _, act := range instructions[0].Accounts() {
			if act.IsSigner {
				feePayer = act.PublicKey
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("cannot determine fee payer. You can ether pass the fee payer via the 'TransactionWithInstructions' option parameter or it falls back to the first instruction's first signer")
		}
	}

	programIDs := make([]common.Address, 0)
	accounts := []*AccountMeta{}
	for _, instruction := range instructions {
		accounts = append(accounts, instruction.Accounts()...)
		programIDs = core.UniqueAppend(programIDs, instruction.ProgramID())
	}

	// Add programID to the account list
	for _, programID := range programIDs {
		accounts = append(accounts, &AccountMeta{
			PublicKey:  programID,
			IsSigner:   false,
			IsWritable: false,
		})
	}

	// Sort. Prioritizing first by signer, then by writable
	sort.SliceStable(accounts, func(i, j int) bool {
		return accounts[i].less(accounts[j])
	})

	var (
		uniqAccounts    []*AccountMeta
		uniqAccountsMap = map[common.Address]uint64{}
	)
	for _, acc := range accounts {
		if index, found := uniqAccountsMap[acc.PublicKey]; found {
			uniqAccounts[index].IsWritable = uniqAccounts[index].IsWritable || acc.IsWritable
			continue
		}
		uniqAccounts = append(uniqAccounts, acc)
		uniqAccountsMap[acc.PublicKey] = uint64(len(uniqAccounts) - 1)
	}

	// Move fee payer to the front
	feePayerIndex := -1
	for idx, acc := range uniqAccounts {
		if acc.PublicKey == feePayer {
			feePayerIndex = idx
		}
	}

	accountCount := len(uniqAccounts)
	if feePayerIndex < 0 {
		// fee payer is not part of accounts we want to add it
		accountCount++
	}
	finalAccounts := make([]*AccountMeta, accountCount)

	itr := 1
	for idx, uniqAccount := range uniqAccounts {
		if idx == feePayerIndex {
			uniqAccount.IsSigner = true
			uniqAccount.IsWritable = true
			finalAccounts[0] = uniqAccount
			continue
		}
		finalAccounts[itr] = uniqAccount
		itr++
	}

	if feePayerIndex < 0 {
		// fee payer is not part of accounts we want to add it
		feePayerAccount := &AccountMeta{
			PublicKey:  feePayer,
			IsSigner:   true,
			IsWritable: true,
		}
		finalAccounts[0] = feePayerAccount
	}

	message := Message{
		RecentBlockhash: recentBlockHash,
	}
	accountKeyIndex := map[string]uint16{}
	for idx, acc := range finalAccounts {
		message.AccountKeys = append(message.AccountKeys, acc.PublicKey)
		accountKeyIndex[acc.PublicKey.String()] = uint16(idx)
		if acc.IsSigner {
			message.Header.NumRequiredSignatures++
			if !acc.IsWritable {
				message.Header.NumReadonlySignedAccounts++
			}
			continue
		}

		if !acc.IsWritable {
			message.Header.NumReadonlyUnsignedAccounts++
		}
	}

	for txIdx, instruction := range instructions {
		accounts = instruction.Accounts()
		accountIndex := make([]uint16, len(accounts))
		for idx, acc := range accounts {
			accountIndex[idx] = accountKeyIndex[acc.PublicKey.String()]
		}
		data, err := instruction.Data()
		if err != nil {
			return nil, fmt.Errorf("unable to encode instructions [%d]: %w", txIdx, err)
		}
		message.Instructions = append(message.Instructions, CompiledInstruction{
			ProgramIDIndex: accountKeyIndex[instruction.ProgramID().String()],
			Accounts:       accountIndex,
			Data:           common.SolData{RawData: data, Encoding: "base58"},
		})
	}

	return &Transaction{
		Message: message,
	}, nil
}

func (tx *Transaction) MarshalBinary() ([]byte, error) {
	if len(tx.Signatures) == 0 || len(tx.Signatures) != int(tx.Message.Header.NumRequiredSignatures) {
		return nil, errors.New("signature verification failed")
	}

	messageContent, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to encode tx.Message to binary: %w", err)
	}

	var signatureCount []byte
	encodbin.EncodeCompactU16Length(&signatureCount, len(tx.Signatures))
	output := make([]byte, 0, len(signatureCount)+len(signatureCount)*64+len(messageContent))
	output = append(output, signatureCount...)
	for _, sig := range tx.Signatures {
		output = append(output, sig[:]...)
	}
	output = append(output, messageContent...)

	return output, nil
}

// UnmarshalJSON parses the transaction Content
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	// Unmarshal data to []byte
	var data interface{}
	// Unmarshal to map
	if err := json.Unmarshal(input, &data); err != nil {
		return err
	}
	// core.BeautifyConsole("aaa", data)
	// get active type
	switch v := data.(type) {
	// slice
	case []interface{}:
		// none data
		if len(v) == 0 {
			return nil
		}
		// decode to string
		switch v[1] {
		case "base58":
			tx.UnmarshalBase58(v[0].(string))
		case "base64":
			tx.UnmarshalBase64(v[0].(string))
		default:
			return fmt.Errorf("UnmarshalDataByEncoding Err: %s", v[1])
		}
	case map[string]interface{}:
		//
		for k, vv := range v {
			switch k {
			case "message":
				message := vv.(map[string]interface{})
				// Message
				for _, accKeys := range message["accountKeys"].([]interface{}) {
					tx.Message.AccountKeys = append(tx.Message.AccountKeys, common.Base58ToAddress(accKeys.(string)))
				}
				for _, insMap := range message["instructions"].([]interface{}) {
					var (
						accounts []uint16
						stackHeight *uint16
						ins = insMap.(map[string]interface{})
					)
					for _, acc := range ins["accounts"].([]interface{}) {
						accounts = append(accounts, uint16(acc.(float64)))
					}
					// exist stack height
					if ins["stackHeight"] != nil {
						stackHeight = ins["stackHeight"].(*uint16)
					}
					insData, _ := base58.Decode(ins["data"].(string))
					// Instructions
					tx.Message.Instructions = append(tx.Message.Instructions, CompiledInstruction{
						StackHeight:    stackHeight,
						ProgramIDIndex: uint16(ins["programIdIndex"].(float64)),
						Accounts:       accounts,
						Data: common.SolData{
							RawData:  insData,
							Encoding: "base58",
						},
					})
				}
				msgHeader, ok := message["header"].(map[string]interface{})
				if !ok {
					continue
				}
				tx.Message.Header = MessageHeader{
					NumRequiredSignatures:       uint8(msgHeader["numRequiredSignatures"].(float64)),
					NumReadonlySignedAccounts:   uint8(msgHeader["numReadonlySignedAccounts"].(float64)),
					NumReadonlyUnsignedAccounts: uint8(msgHeader["numReadonlyUnsignedAccounts"].(float64)),
				}
				tx.Message.RecentBlockhash = common.Base58ToHash(message["recentBlockhash"].(string))
			case "signatures":
				for _, sig := range vv.([]interface{}) {
					tx.Signatures = append(tx.Signatures, common.Base58ToSignature(sig.(string)))
				}
			}
		}
		return nil
	}
	// UnSerialize transaction content
	return nil
}

// UnmarshalBase64 decodes a base64 encoded transaction.
func (tx *Transaction) UnmarshalBase64(b64 string) error {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return tx.UnmarshalWithDecoder(encodbin.NewBinDecoder(b))
}

// UnmarshalBase58 decodes a base58 encoded transaction.
func (tx *Transaction) UnmarshalBase58(b58 string) error {
	b, _ := base58.Decode(b58)
	return tx.UnmarshalWithDecoder(encodbin.NewBinDecoder(b))
}

func (tx *Transaction) UnmarshalWithDecoder(decoder *encodbin.Decoder) (err error) {
	{
		numSignatures, err := decoder.ReadCompactU16()
		if err != nil {
			return fmt.Errorf("unable to read numSignatures: %w", err)
		}

		tx.Signatures = make([]common.Signature, numSignatures)
		for i := 0; i < numSignatures; i++ {
			_, err := decoder.Read(tx.Signatures[i][:])
			if err != nil {
				return fmt.Errorf("unable to read tx.Signatures[%d]: %w", i, err)
			}
		}
	}
	{
		err := tx.Message.UnmarshalWithDecoder(decoder)
		if err != nil {
			return fmt.Errorf("unable to decode tx.Message: %w", err)
		}
	}
	return nil
}

type privateKeyGetter func(key common.Address) *crypto.Account

func (tx *Transaction) Sign(getter privateKeyGetter) (out []common.Signature, err error) {
	messageContent, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("unable to encode message for signing: %w", err)
	}

	signerKeys := tx.Message.signerKeys()

	for _, key := range signerKeys {
		privateKey := getter(key)
		if privateKey == nil {
			return nil, fmt.Errorf("signer key %q not found. Ensure all the signer keys are in the vault", key.String())
		}

		s := privateKey.Sign(messageContent)

		tx.Signatures = append(tx.Signatures, common.BytesToSignature(s))
	}
	return tx.Signatures, nil
}

func (tx Transaction) ToBase64() (string, error) {
	out, err := tx.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

func (tx Transaction) ToBase58() (string, error) {
	out, err := tx.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base58.Encode(out), nil
}
