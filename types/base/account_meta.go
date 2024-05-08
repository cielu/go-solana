package base

import "github.com/cielu/go-solana/common"

type AccountMeta struct {
	PublicKey  common.Address `json:"publickey"`
	IsWritable bool
	IsSigner   bool
}

// type AccountsSettable interface {
// 	SetAccounts(accounts []*AccountMeta) error
// }
//
// type AccountsGettable interface {
// 	GetAccounts() (accounts []*AccountMeta)
// }

func GetAccounts(slice []*AccountMeta) (accounts []*AccountMeta) {
	// range slice
	for _, meta := range slice {
		if meta != nil {
			accounts = append(accounts, meta)
		}
	}
	return
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

func (meta *AccountMeta) Less(act *AccountMeta) bool {
	if meta.IsSigner != act.IsSigner {
		return meta.IsSigner
	}
	if meta.IsWritable != act.IsWritable {
		return meta.IsWritable
	}
	return false
}
