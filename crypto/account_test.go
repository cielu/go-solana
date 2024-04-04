package crypto

import (
	"crypto/ed25519"
	"fmt"
	"github.com/cielu/go-solana/common"
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
		privAddr common.Address
		//message  = []byte("hello word")
	)
	copy(privAddr[:], account.PrivateKey.Public().(ed25519.PublicKey))
	// Got address from private key
	if account.Address != privAddr {
		t.Errorf("account address not eq priv address. Want: %s, Got: %s", account.Address, privAddr)
	}
	// AccountFromBytes
	account2, err := AccountFromBytes(account.PrivateKey)
	if err != nil {
		t.Errorf("AccountFromBytes Failed: %s", err.Error())
	}
	//
	if account.Address != account2.Address {
		t.Errorf("account address not eq account2. Want: %s, Got: %s", account.Address, account2.Address)
	}
	// AccountFromSeed
	account3, err := AccountFromSeed(account.PrivateKey.Seed())
	if err != nil {
		t.Errorf("AccountFromSeed Failed: %s", err.Error())
	}
	if account.Address != account3.Address {
		t.Errorf("account address not eq account3. Want: %s, Got: %s", account.Address, account3.Address)
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
	if account.Address != account4.Address {
		t.Errorf("account address not eq account4. Want: %s, Got: %s", account.Address, account4.Address)
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
	if account.Address != account5.Address {
		t.Errorf("account address not eq account5. Want: %s, Got: %s", account.Address, account5.Address)
	}
	mnemonic := "letter advice cage absurd amount doctor acoustic avoid letter advice cage above"
	// Account
	account6, err := AccountFromMnemonic(mnemonic, true)
	// account6
	if account.Address != account6.Address {
		//t.Errorf("account address not eq account6. Want: %s, Got: %s", account.Address, account6.Address)
	}
	fmt.Println("account6:", account6.Address)
}
