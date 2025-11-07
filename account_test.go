package solana

import (
	"crypto/ed25519"
	"fmt"
	"testing"
)

func TestAccount(t *testing.T) {
	// generate a new account
	account, err := GenerateAccount()
	// generate failed
	if err != nil {
		t.Errorf("GenearteAaccount Failed: %s", err.Error())
	}
	var (
		privAddr PublicKey
		//message  = []byte("hello word")
	)
	copy(privAddr[:], account.PrivateKey.Public().(ed25519.PublicKey))
	// Got PublicKey from private key
	if account.PublicKey != privAddr {
		t.Errorf("account PublicKey not eq priv PublicKey. Want: %s, Got: %s", account.PublicKey, privAddr)
	}
	// AccountFromBytes
	account2, err := AccountFromBytes(account.PrivateKey)
	if err != nil {
		t.Errorf("AccountFromBytes Failed: %s", err.Error())
	}
	//
	if account.PublicKey != account2.PublicKey {
		t.Errorf("account PublicKey not eq account2. Want: %s, Got: %s", account.PublicKey, account2.PublicKey)
	}
	// AccountFromSeed
	account3, err := AccountFromSeed(account.PrivateKey.Seed())
	if err != nil {
		t.Errorf("AccountFromSeed Failed: %s", err.Error())
	}
	if account.PublicKey != account3.PublicKey {
		t.Errorf("account PublicKey not eq account3. Want: %s, Got: %s", account.PublicKey, account3.PublicKey)
	}
	// GenerateBase58PrvKey
	base58Key, err := GenerateBase58PrvKey(account)
	if err != nil {
		t.Errorf("GenerateBase58PrvKey Failed: %s", err.Error())
	}
	fmt.Println("base58 Key:", base58Key)
	// AccountFromBase58
	account4, err := AccountFromBase58Key(base58Key)
	if err != nil {
		t.Errorf("AccountFromBase58 Failed: %s", err.Error())
	}
	if account.PublicKey != account4.PublicKey {
		t.Errorf("account PublicKey not eq account4. Want: %s, Got: %s", account.PublicKey, account4.PublicKey)
	}
	// GenerateHexPrvKey
	hexKey, err := GenerateHexPrvKey(account)
	if err != nil {
		t.Errorf("GenerateHexPrvKey Failed: %s", err.Error())
	}
	fmt.Println("Hex Key:", hexKey)
	// AccountFromBase58
	account5, err := AccountFromHexKey(hexKey)
	if err != nil {
		t.Errorf("AccountFromHex Failed: %s", err.Error())
	}
	if account.PublicKey != account5.PublicKey {
		t.Errorf("account PublicKey not eq account5. Want: %s, Got: %s", account.PublicKey, account5.PublicKey)
	}
	mnemonic := "letter advice cage absurd amount doctor acoustic avoid letter advice cage above"
	// Account
	account6, err := AccountFromMnemonic(mnemonic, true)
	// account6
	if account.PublicKey != account6.PublicKey {
		//t.Errorf("account PublicKey not eq account6. Want: %s, Got: %s", account.PublicKey, account6.PublicKey)
	}
	fmt.Println("account6:", account6.PublicKey)
}
