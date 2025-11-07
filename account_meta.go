package solana

type AccountsSettable interface {
	SetAccounts(accounts []*AccountMeta) error
}

type AccountsGettable interface {
	GetAccounts() (accounts []*AccountMeta)
}

type AccountMeta struct {
	PublicKey  PublicKey `json:"publickey"`
	IsWritable bool
	IsSigner   bool
}

// Meta intializes a new AccountMeta with the provided pubKey.
func Meta(pubKey PublicKey) *AccountMeta {
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

func NewAccountMeta(pubKey PublicKey, WRITE bool, SIGNER bool) *AccountMeta {
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

type AccountMetaSlice []*AccountMeta

func (slice *AccountMetaSlice) Append(account *AccountMeta) {
	*slice = append(*slice, account)
}

func (slice *AccountMetaSlice) SetAccounts(accounts []*AccountMeta) error {
	*slice = accounts
	return nil
}

func (slice AccountMetaSlice) GetAccounts() (accounts []*AccountMeta) {
	// range slice
	for i := range slice {
		if slice[i] != nil {
			accounts = append(accounts, slice[i])
		}
	}
	return
}

// Get returns the AccountMeta at the desired index.
// If the index is not present, it returns nil.
func (slice AccountMetaSlice) Get(index int) *AccountMeta {
	if len(slice) > index {
		return slice[index]
	}
	return nil
}

// GetSigners returns the accounts that are signers.
func (slice AccountMetaSlice) GetSigners() []*AccountMeta {
	signers := make([]*AccountMeta, 0, len(slice))
	for _, ac := range slice {
		if ac.IsSigner {
			signers = append(signers, ac)
		}
	}
	return signers
}

// GetKeys returns the pubkeys of all AccountMeta.
func (slice AccountMetaSlice) GetKeys() (keys []PublicKey) {
	// range slice
	for _, ac := range slice {
		keys = append(keys, ac.PublicKey)
	}
	return
}

func (slice AccountMetaSlice) Len() int {
	return len(slice)
}


