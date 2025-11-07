package solana

type RawTransaction struct {
	instructions []Instruction
	blockHash    Hash
	payer        PublicKey
	signers      []string
}

func NewRawTransaction(blockHash Hash, payer string, inst []Instruction, signers []string) *RawTransaction {
	return &RawTransaction{
		instructions: inst,
		blockHash:    blockHash,
		payer:        StrToPublicKey(payer),
		signers:      signers,
	}
}
