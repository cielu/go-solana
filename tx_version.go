package solana

import "strconv"

type TxVersion int

const (
	LegacyTransactionVersion TxVersion = -1
	legacyVersion                      = `"legacy"`
)

func (ver *TxVersion) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	s := string(b)
	if s == "null" || s == `""` || s == legacyVersion {
		*ver = LegacyTransactionVersion
		return nil
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*ver = TxVersion(v)
	return nil
}

func (ver TxVersion) MarshalJSON() ([]byte, error) {
	if ver == LegacyTransactionVersion {
		return []byte(legacyVersion), nil
	} else {
		return []byte(strconv.Itoa(int(ver))), nil
	}
}
