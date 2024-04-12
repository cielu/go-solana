package types

import (
	"errors"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/bincode"
)

type NewTransactionParam struct {
	Signatures []Account
	Message    MessageBody
}

type NewMessageParam struct {
	FeePayer                   common.Address
	Instructions               []Instruction
	RecentBlockhash            common.Hash
	AddressLookupTableAccounts []AddressLookupTableAccount
}

type Transaction struct {
	Signatures []common.SignatureUndefinedLength
	Message    MessageBody
}

func NewTransaction(val NewTransactionParam) (Transaction, error) {
	signatures := make([]common.SignatureUndefinedLength, 0, val.Message.Header.NumRequiredSignatures)
	for i := uint8(0); i < val.Message.Header.NumRequiredSignatures; i++ {
		signatures = append(signatures, make([]byte, 64))
	}
	m := map[common.Address]uint8{}
	for i := uint8(0); i < val.Message.Header.NumRequiredSignatures; i++ {
		m[val.Message.Accounts[i]] = i
	}
	data, err := val.Message.Serialize()

	if err != nil {
		return Transaction{}, fmt.Errorf("failed to serialize message, err: %v", err)
	}
	for _, signer := range val.Signatures {
		idx, ok := m[signer.PublicKey]
		if !ok {
			return Transaction{}, fmt.Errorf("%w, %v is not a signer")
		}
		signatures[idx] = signer.Sign(data)
	}

	return Transaction{
		Signatures: signatures,
		Message:    val.Message,
	}, nil
}

func (tx *Transaction) Serialize() ([]byte, error) {
	if len(tx.Signatures) == 0 || len(tx.Signatures) != int(tx.Message.Header.NumRequiredSignatures) {
		return nil, errors.New("Signature verification failed")
	}

	signatureCount := bincode.UintToVarLenBytes(uint64(len(tx.Signatures)))
	messageData, err := tx.Message.Serialize()
	if err != nil {
		return nil, err
	}

	output := make([]byte, 0, len(signatureCount)+len(signatureCount)*64+len(messageData))
	output = append(output, signatureCount...)
	for _, sig := range tx.Signatures {
		output = append(output, sig...)
	}
	output = append(output, messageData...)

	return output, nil
}
