package common

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
)

func TestAddress(t *testing.T) {

	tests := []struct {
		addr string
		want Address
	}{
		{
			addr: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // usdc
			want: Base58ToAddress("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
		},
	}

	for _, test := range tests {
		// base58 address
		addr := Base58ToAddress(test.addr)

		if addr != test.want {
			t.Errorf("Go Address Err ==> Got %s, Want: %s", addr, test.want)
		}

		if addr.String() != test.addr {
			t.Errorf("Go Address Err ==> Got %s, Want: %s", addr, test.want)
		}
	}
	// Random a pub key
	pub, prv, _ := ed25519.GenerateKey(rand.Reader)
	// --> set Bytes
	var (
		addr1, addr2 Address
		pubKey       = make([]byte, AddressLength)
	)
	addr1.SetBytes(pub)
	//
	copy(pubKey[:], prv.Public().(ed25519.PublicKey))
	//
	addr2.SetBytes(pubKey)
	//
	if addr1 != addr2 {
		t.Errorf("pub address not eq prv address. Got addr1: %s, addr2: %s", addr1, addr2)
	}
	// log addr1 and addr2
	t.Logf("addr1: %s, addr2: %s", addr1, addr2)
}
