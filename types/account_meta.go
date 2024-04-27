package types

import "github.com/cielu/go-solana/common"

type AccountMeta struct {
	PublicKey  common.Address `json:"publickey"`
	IsWritable bool
	IsSigner   bool
}


type AccountsSettable interface {
	SetAccounts(accounts []*AccountMeta) error
}

type AccountsGettable interface {
	GetAccounts() (accounts []*AccountMeta)
}

// Meta intializes a new AccountMeta with the provided pubKey.
func Meta(pubKey common.Address) *AccountMeta {
	return &AccountMeta{
		PublicKey: pubKey,
	}
}

// WRITE sets IsWritable to true.
func (meta *AccountMeta) WRITE() *AccountMeta {
	meta.IsWritable = true
	return meta
}

// SIGNER sets IsSigner to true.
func (meta *AccountMeta) SIGNER() *AccountMeta {
	meta.IsSigner = true
	return meta
}

func NewAccountMeta(pubKey common.Address, WRITE bool, SIGNER bool) *AccountMeta {
	return &AccountMeta{
		PublicKey:  pubKey,
		IsWritable: WRITE,
		IsSigner:   SIGNER,
	}
}

func (a *AccountMeta) less(act *AccountMeta) bool {
	if a.IsSigner != act.IsSigner {
		return a.IsSigner
	}
	if a.IsWritable != act.IsWritable {
		return a.IsWritable
	}
	return false
}
