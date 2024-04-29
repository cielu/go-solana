package base

import "github.com/cielu/go-solana/common"

type AccountMeta struct {
	PublicKey  common.Address `json:"publickey"`
	IsWritable bool
	IsSigner   bool
}

type AccountMetaSlice []*AccountMeta

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

func (slice AccountMetaSlice) GetAccounts() []*AccountMeta {
	out := make([]*AccountMeta, 0, len(slice))
	for i := range slice {
		if slice[i] != nil {
			out = append(out, slice[i])
		}
	}
	return out
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

func calcSplitAtLengths(total int, index int) (int, int) {
	if index == 0 {
		return 0, total
	}
	if index > total-1 {
		return total, 0
	}
	return index, total - index
}

func (slice AccountMetaSlice) SplitFrom(index int) (AccountMetaSlice, AccountMetaSlice) {
	if index < 0 {
		panic("negative index")
	}
	if index == 0 {
		return AccountMetaSlice{}, slice
	}
	if index > len(slice)-1 {
		return slice, AccountMetaSlice{}
	}

	firstLen, secondLen := calcSplitAtLengths(len(slice), index)

	first := make(AccountMetaSlice, firstLen)
	copy(first, slice[:index])

	second := make(AccountMetaSlice, secondLen)
	copy(second, slice[index:])

	return first, second
}
