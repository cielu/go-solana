// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package encodbin

import (
	"bytes"
	"fmt"
)

type EncoderDecoder interface {
	BinaryMarshaler
	BinaryUnmarshaler
}

func MarshalBin(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := NewBinEncoder(buf)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}

func UnmarshalBin(v interface{}, b []byte) error {
	decoder := NewBinDecoder(b)
	return decoder.Decode(v)
}

type byteCounter struct {
	count uint64
}

func (c *byteCounter) Write(p []byte) (n int, err error) {
	c.count += uint64(len(p))
	return len(p), nil
}

// BinByteCount computes the byte count size for the received populated structure. The reported size
// is the one for the populated structure received in arguments. Depending on how serialization of
// your fields is performed, size could vary for different structure.
func BinByteCount(v interface{}) (uint64, error) {
	counter := byteCounter{}
	err := NewBinEncoder(&counter).Encode(v)
	if err != nil {
		return 0, fmt.Errorf("encode %T: %w", v, err)
	}
	return counter.count, nil
}

// MustBinByteCount acts just like BinByteCount but panics if it encounters any encoding errors.
func MustBinByteCount(v interface{}) uint64 {
	count, err := BinByteCount(v)
	if err != nil {
		panic(err)
	}
	return count
}

