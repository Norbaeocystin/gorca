package gorca

/*
TODO implement logic from https://github.com/orca-so/whirlpools/blob/main/sdk/src/utils/public/pda-utils.ts
*/

import (
	"encoding/binary"
	"github.com/gagliardetto/solana-go"
)

const PDA_WHIRLPOOL_SEED = "whirlpool"

const PDA_POSITION_SEED = "position"

const PDA_METADATA_SEED = "metadata"

const PDA_TICK_ARRAY_SEED = "tick_array"

const PDA_FEE_TIER_SEED = "fee_tier"

const PDA_ORACLE_SEED = "oracle"

func GetOracle(programID, whirlpooAddress solana.PublicKey) (solana.PublicKey, uint8, error) {
	address, bump, err := solana.FindProgramAddress([][]byte{[]byte(PDA_ORACLE_SEED), whirlpooAddress.Bytes()}, programID)
	return address, bump, err
}

// TODO implement
func GetWhirlPool(programID, configKey, mintA, mintB solana.PublicKey, tickSpacing uint16) (solana.PublicKey, uint8, error) {
	buff := [2]byte{}
	binary.LittleEndian.PutUint16(buff[:], tickSpacing)
	address, bump, err := solana.FindProgramAddress([][]byte{[]byte(PDA_WHIRLPOOL_SEED), configKey.Bytes(), mintA.Bytes(), mintB.Bytes(), buff[:]}, programID)
	return address, bump, err
}

func GetPosition(programID, positionMintKey solana.PublicKey) (solana.PublicKey, uint8, error) {
	address, bump, err := solana.FindProgramAddress([][]byte{[]byte(PDA_POSITION_SEED), positionMintKey.Bytes()}, programID)
	return address, bump, err
}

//func GetPositionMetadata(positionMintKey solana.PublicKey) (solana.PublicKey, uint8, error) {
//	address, bump, err := solana.FindProgramAddress([][]byte{[]byte(PDA_METADATA_SEED), METADATA_PROGRAM_ADDRESS.Bytes(), positionMintKey.Bytes()}, METADATA_PROGRAM_ADDRESS)
//	return address, bump, err
//}

func GetTickArrayViaFPA(programID, whirlpoolAddress solana.PublicKey, startTick int64) (solana.PublicKey, uint8, error) {
	address, bump, err := solana.FindProgramAddress([][]byte{[]byte(PDA_TICK_ARRAY_SEED), whirlpoolAddress.Bytes(), []byte(string(startTick))}, programID)
	return address, bump, err
}

//func GetTickArrayFromTickIndex(){}
//func GetTickArrayFromSqrtPrice(){}
//func GetFeeTier(){}
