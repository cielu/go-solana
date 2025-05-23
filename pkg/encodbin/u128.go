package encodbin

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

// Uint128
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

func (i Uint128) getByteOrder() binary.ByteOrder {
	if i.Endianness == nil {
		return defaultByteOrder
	}
	return i.Endianness
}

func (i Int128) getByteOrder() binary.ByteOrder {
	return Uint128(i).getByteOrder()
}

func (i Float128) getByteOrder() binary.ByteOrder {
	return Uint128(i).getByteOrder()
}

func (i Uint128) Bytes() []byte {
	buf := make([]byte, 16)
	order := i.getByteOrder()
	if order == binary.LittleEndian {
		order.PutUint64(buf[:8], i.Lo)
		order.PutUint64(buf[8:], i.Hi)
		ReverseBytes(buf)
	} else {
		order.PutUint64(buf[:8], i.Hi)
		order.PutUint64(buf[8:], i.Lo)
	}
	return buf
}

func (i Uint128) BigInt() *big.Int {
	buf := i.Bytes()
	value := (&big.Int{}).SetBytes(buf)
	return value
}

func (i Uint128) String() string {
	// Same for Int128, Float128
	return i.DecimalString()
}

func (i Uint128) DecimalString() string {
	return i.BigInt().String()
}

func (i Uint128) HexString() string {
	number := i.Bytes()
	return fmt.Sprintf("0x%s", hex.EncodeToString(number))
}

func (i Uint128) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + i.String() + `"`), nil
}

func (i *Uint128) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return i.unmarshalJSON_hex(s)
	}

	return i.unmarshalJSON_decimal(s)
}

func (i *Uint128) unmarshalJSON_decimal(s string) error {
	parsed, ok := (&big.Int{}).SetString(s, 0)
	if !ok {
		return fmt.Errorf("could not parse %q", s)
	}
	oo := parsed.FillBytes(make([]byte, 16))
	ReverseBytes(oo)

	dec := NewBinDecoder(oo)

	out, err := dec.ReadUint128(i.getByteOrder())
	if err != nil {
		return err
	}
	i.Lo = out.Lo
	i.Hi = out.Hi

	return nil
}

func (i *Uint128) unmarshalJSON_hex(s string) error {
	truncatedVal := s[2:]
	if len(truncatedVal) != 16 {
		return fmt.Errorf("uint128 expects 16 characters after 0x, had %v", len(truncatedVal))
	}

	data, err := hex.DecodeString(truncatedVal)
	if err != nil {
		return err
	}

	order := i.getByteOrder()
	if order == binary.LittleEndian {
		i.Lo = order.Uint64(data[:8])
		i.Hi = order.Uint64(data[8:])
	} else {
		i.Hi = order.Uint64(data[:8])
		i.Lo = order.Uint64(data[8:])
	}

	return nil
}

func (i *Uint128) UnmarshalWithDecoder(dec *Decoder) error {
	var order binary.ByteOrder
	if dec != nil && dec.currentFieldOpt != nil {
		order = dec.currentFieldOpt.Order
	} else {
		order = i.getByteOrder()
	}
	value, err := dec.ReadUint128(order)
	if err != nil {
		return err
	}

	*i = value
	return nil
}

func (i Uint128) MarshalWithEncoder(enc *Encoder) error {
	var order binary.ByteOrder
	if enc != nil && enc.currentFieldOpt != nil {
		order = enc.currentFieldOpt.Order
	} else {
		order = i.getByteOrder()
	}
	return enc.WriteUint128(i, order)
}

// SetBigInt bigInt to Uint128
func (i *Uint128) SetBigInt(b *big.Int) error {
	if b.Sign() < 0 {
		return fmt.Errorf("cannot assign negative integer: %v", b)
	} else if b.BitLen() > 128 {
		return fmt.Errorf("value overflows Uint128")
	}
	i.Lo = b.Uint64()
	i.Hi = b.Rsh(b, 64).Uint64()
	return nil
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
