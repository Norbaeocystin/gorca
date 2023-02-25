// Code Generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package whirlpool

import (
	"bytes"
	"fmt"
	ag_binary "github.com/gagliardetto/binary"
)

func encodeT(data interface{}, buf *bytes.Buffer) error {
	if err := ag_binary.NewBorshEncoder(buf).Encode(data); err != nil {
		return fmt.Errorf("unable to encode instruction: %w", err)
	}
	return nil
}

func decodeT(dst interface{}, data []byte) error {
	return ag_binary.NewBorshDecoder(data).Decode(dst)
}
