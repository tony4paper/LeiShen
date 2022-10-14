package edb

import "encoding/binary"

const (
	keyLenthUint64 = 8
)

func EncodeUint64(number uint64) []byte {
	enc := make([]byte, keyLenthUint64)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

func DecodeUint64(data []byte) *uint64 {
	if len(data) != keyLenthUint64 {
		return nil
	}
	i := binary.BigEndian.Uint64(data)
	return &i
}
