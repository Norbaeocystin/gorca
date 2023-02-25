package gorca

import (
	"encoding/binary"
	"errors"
	bin "github.com/gagliardetto/binary"
	"math/big"
)

func BigIntToBinUint128(value *big.Int) (bin.Uint128, error) {
	var returnValue bin.Uint128
	bytes := value.Bytes()
	if len(bytes) > 16 {
		return returnValue, errors.New("Overflow, too many bytes from big.Int")
	}
	switch {
	case len(bytes) <= 8:
		low := make([]byte, 8)
		for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
			bytes[i], bytes[j] = bytes[j], bytes[i]
		}
		for idx, byt := range bytes {
			low[idx] = byt
		}
		returnValue.Lo = binary.LittleEndian.Uint64(low)
	case len(bytes) > 8 && len(bytes) < 17:
		for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
			bytes[i], bytes[j] = bytes[j], bytes[i]
		}
		low := make([]byte, 8)
		high := make([]byte, 8)
		highBytes := bytes[8:]
		lowBytes := bytes[:8]
		for idx, byt := range highBytes {
			high[idx] = byt
		}
		for idx, byt := range lowBytes {
			low[idx] = byt
		}
		returnValue.Lo = binary.LittleEndian.Uint64(low)
		returnValue.Hi = binary.LittleEndian.Uint64(high)
	}
	//low := make([]byte, 8)
	//high := make([]byte, 8)
	//log.Println(low, high)
	return returnValue, nil
}
