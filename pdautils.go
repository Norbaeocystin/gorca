package gorca

/*
TODO implement logic from https://github.com/orca-so/whirlpools/blob/main/sdk/src/utils/public/pda-utils.ts
*/

import "github.com/gagliardetto/solana-go"

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
//func GetWhirlPool() {}
//func GetPosition() {}
//func GetPositionMetadata() {}
//func GetTickArray(){}
//func GetTickArrayFromTickIndex(){}
//func GetTickArrayFromSqrtPrice(){}
//func GetFeeTier(){}
