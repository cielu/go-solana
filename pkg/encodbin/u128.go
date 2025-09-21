package encodbin

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

// Uint128 en
type Uint128 struct {
	Lo         uint64
	Hi         uint64
	Endianness binary.ByteOrder
}

func NewUint128BigEndian() *Uint128 {
	return &Uint128{
		Endianness: binary.BigEndian,
	}
}

func NewUint128LittleEndian() *Uint128 {
	return &Uint128{
		Endianness: binary.LittleEndian,
	}
}

func ReverseBytes(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (u Uint128) IsZero() bool {
	// NOTE: we do not compare against Zero, because that is a global variable
	// that could be modified.
	return u == Uint128{}
}

// Cmp compares u and v and returns:
//
//	-1 if u <  v
//	 0 if u == v
//	+1 if u >  v
func (u Uint128) Cmp(v Uint128) int {
	if u == v {
		return 0
	} else if u.Hi < v.Hi || (u.Hi == v.Hi && u.Lo < v.Lo) {
		return -1
	} else {
		return 1
	}
}

// Cmp64 compares u and v and returns:
//
//	-1 if u <  v
//	 0 if u == v
//	+1 if u >  v
func (u Uint128) Cmp64(v uint64) int {
	if u.Hi == 0 && u.Lo == v {
		return 0
	} else if u.Hi == 0 && u.Lo < v {
		return -1
	} else {
		return 1
	}
}

// Xor returns u^v.
func (u Uint128) Xor(v Uint128) Uint128 {
	return Uint128{u.Lo ^ v.Lo, u.Hi ^ v.Hi, u.Endianness}
}

// Xor64 returns u^v.
func (u Uint128) Xor64(v uint64) Uint128 {
	return Uint128{u.Lo ^ v, u.Hi ^ 0, u.Endianness}
}

// Lsh returns u<<n.
func (u Uint128) Lsh(n uint) (s Uint128) {
	if n > 64 {
		s.Lo = 0
		s.Hi = u.Lo << (n - 64)
	} else {
		s.Lo = u.Lo << n
		s.Hi = u.Hi<<n | u.Lo>>(64-n)
	}
	return
}

// Rsh returns u>>n.
func (u Uint128) Rsh(n uint) (s Uint128) {
	if n > 64 {
		s.Lo = u.Hi >> (n - 64)
		s.Hi = 0
	} else {
		s.Lo = u.Lo>>n | u.Hi<<(64-n)
		s.Hi = u.Hi >> n
	}
	return
}

func getByteOrder(endia binary.ByteOrder) binary.ByteOrder {
	if endia == nil {
		return defaultByteOrder
	}
	return endia
}

func (i Int128) getByteOrder() binary.ByteOrder {
	return getByteOrder(i.Endianness)
}

func (i Float128) getByteOrder() binary.ByteOrder {
	return getByteOrder(i.Endianness)
}

func (u Uint128) Bytes() []byte {
	buf := make([]byte, 16)
	order := getByteOrder(u.Endianness)
	if order == binary.LittleEndian {
		order.PutUint64(buf[:8], u.Lo)
		order.PutUint64(buf[8:], u.Hi)
		ReverseBytes(buf)
	} else {
		order.PutUint64(buf[:8], u.Hi)
		order.PutUint64(buf[8:], u.Lo)
	}
	return buf
}

func (u Uint128) BigInt() *big.Int {
	buf := u.Bytes()
	value := (&big.Int{}).SetBytes(buf)
	return value
}

func (u Uint128) String() string {
	// Same for Int128, Float128
	return u.DecimalString()
}

func (u Uint128) DecimalString() string {
	return u.BigInt().String()
}

func (u Uint128) HexString() string {
	number := u.Bytes()
	return fmt.Sprintf("0x%s", hex.EncodeToString(number))
}

// FromBigInt bigInt to Uint128
func (u *Uint128) FromBigInt(b *big.Int) error {
	if b.Sign() < 0 {
		return fmt.Errorf("cannot assign negative integer: %v", b)
	} else if b.BitLen() > 128 {
		return fmt.Errorf("value overflows Uint128")
	}
	u.Lo = b.Uint64()
	u.Hi = b.Rsh(b, 64).Uint64()
	return nil
}

// FromString parses s as a Uint128 value.
func FromString(s string) (u Uint128, err error) {
	_, err = fmt.Sscan(s, &u)
	return
}

// MarshalText implements encoding.TextMarshaler.
func (u Uint128) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *Uint128) UnmarshalText(b []byte) error {
	_, err := fmt.Sscan(string(b), u)
	return err
}

func (u Uint128) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + u.String() + `"`), nil
}

func (u *Uint128) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return u.unmarshalJSON_hex(s)
	}

	return u.unmarshalJSON_decimal(s)
}

func (u *Uint128) unmarshalJSON_decimal(s string) error {
	parsed, ok := (&big.Int{}).SetString(s, 0)
	if !ok {
		return fmt.Errorf("could not parse %q", s)
	}
	oo := parsed.FillBytes(make([]byte, 16))
	ReverseBytes(oo)

	dec := NewBinDecoder(oo)

	out, err := dec.ReadUint128(getByteOrder(u.Endianness))
	if err != nil {
		return err
	}
	u.Lo = out.Lo
	u.Hi = out.Hi

	return nil
}

func (u *Uint128) unmarshalJSON_hex(s string) error {
	truncatedVal := s[2:]
	if len(truncatedVal) != 16 {
		return fmt.Errorf("uint128 expects 16 characters after 0x, had %v", len(truncatedVal))
	}

	data, err := hex.DecodeString(truncatedVal)
	if err != nil {
		return err
	}

	order := getByteOrder(u.Endianness)
	if order == binary.LittleEndian {
		u.Lo = order.Uint64(data[:8])
		u.Hi = order.Uint64(data[8:])
	} else {
		u.Hi = order.Uint64(data[:8])
		u.Lo = order.Uint64(data[8:])
	}

	return nil
}

func (u *Uint128) UnmarshalWithDecoder(dec *Decoder) error {
	var order binary.ByteOrder
	if dec != nil && dec.currentFieldOpt != nil {
		order = dec.currentFieldOpt.Order
	} else {
		order = getByteOrder(u.Endianness)
	}
	value, err := dec.ReadUint128(order)
	if err != nil {
		return err
	}

	*u = value
	return nil
}

func (u Uint128) MarshalWithEncoder(enc *Encoder) error {
	var order binary.ByteOrder
	if enc != nil && enc.currentFieldOpt != nil {
		order = enc.currentFieldOpt.Order
	} else {
		order = getByteOrder(u.Endianness)
	}
	return enc.WriteUint128(u, order)
}

// Int128
type Int128 Uint128

func (i Int128) BigInt() *big.Int {
	comp := byte(0x80)
	buf := Uint128(i).Bytes()

	var value *big.Int
	if (buf[0] & comp) == comp {
		buf = twosComplement(buf)
		value = (&big.Int{}).SetBytes(buf)
		value = value.Neg(value)
	} else {
		value = (&big.Int{}).SetBytes(buf)
	}
	return value
}

func (i Int128) String() string {
	return Uint128(i).String()
}

func (i Int128) DecimalString() string {
	return i.BigInt().String()
}

func (i Int128) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + Uint128(i).String() + `"`), nil
}

func (i *Int128) UnmarshalJSON(data []byte) error {
	var el Uint128
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}

	out := Int128(el)
	*i = out

	return nil
}

func (i *Int128) UnmarshalWithDecoder(dec *Decoder) error {
	var order binary.ByteOrder
	if dec != nil && dec.currentFieldOpt != nil {
		order = dec.currentFieldOpt.Order
	} else {
		order = i.getByteOrder()
	}
	value, err := dec.ReadInt128(order)
	if err != nil {
		return err
	}

	*i = value
	return nil
}

func (i Int128) MarshalWithEncoder(enc *Encoder) error {
	var order binary.ByteOrder
	if enc != nil && enc.currentFieldOpt != nil {
		order = enc.currentFieldOpt.Order
	} else {
		order = i.getByteOrder()
	}
	return enc.WriteInt128(i, order)
}

type Float128 Uint128

func (i Float128) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + Uint128(i).String() + `"`), nil
}

func (i *Float128) UnmarshalJSON(data []byte) error {
	var el Uint128
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}

	out := Float128(el)
	*i = out

	return nil
}

func (i *Float128) UnmarshalWithDecoder(dec *Decoder) error {
	var order binary.ByteOrder
	if dec != nil && dec.currentFieldOpt != nil {
		order = dec.currentFieldOpt.Order
	} else {
		order = i.getByteOrder()
	}
	value, err := dec.ReadFloat128(order)
	if err != nil {
		return err
	}

	*i = Float128(value)
	return nil
}

func (i Float128) MarshalWithEncoder(enc *Encoder) error {
	var order binary.ByteOrder
	if enc != nil && enc.currentFieldOpt != nil {
		order = enc.currentFieldOpt.Order
	} else {
		order = i.getByteOrder()
	}
	return enc.WriteUint128(Uint128(i), order)
}

func twosComplement(v []byte) []byte {
	buf := make([]byte, len(v))
	for i, b := range v {
		buf[i] = b ^ byte(0xff)
	}
	one := big.NewInt(1)
	value := (&big.Int{}).SetBytes(buf)
	return value.Add(value, one).Bytes()
}
