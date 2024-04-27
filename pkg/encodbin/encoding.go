package encodbin

import "encoding/binary"

type Encoding int

const (
	EncodingBin Encoding = iota
	EncodingCompactU16
	EncodingBorsh
)

func (enc Encoding) String() string {
	switch enc {
	case EncodingBin:
		return "Bin"
	case EncodingCompactU16:
		return "CompactU16"
	case EncodingBorsh:
		return "Borsh"
	default:
		return ""
	}
}

func (enc Encoding) IsBorsh() bool {
	return enc == EncodingBorsh
}

func (enc Encoding) IsBin() bool {
	return enc == EncodingBin
}

func (enc Encoding) IsCompactU16() bool {
	return enc == EncodingCompactU16
}

func isValidEncoding(enc Encoding) bool {
	switch enc {
	case EncodingBin, EncodingCompactU16, EncodingBorsh:
		return true
	default:
		return false
	}
}

var LE binary.ByteOrder = binary.LittleEndian
var BE binary.ByteOrder = binary.BigEndian
