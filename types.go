package gorca

import (
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type Account struct {
}

type LastUpdate struct {
	Slot  bin.Uint64
	Stale uint8
}

type WhirlpoolData struct {
	//// Account                    solana.PublicKey
	// Version    uint8
	// LastUpdate LastUpdate
	// Padding          [13]byte
	Padding          [8]byte
	WhirlpoolsConfig *solana.PublicKey // 32
	WhirlpoolBump    [1]uint8          // 1 - 33
	TickSpacing      uint16
	TickSpacingSeed  [2]uint8 // 2 - 37
	// Padding [45]byte
	FeeRate uint16 // 4 - 41
	// Padding         [47]byte
	ProtocolFeeRate uint16 // 4 - 45
	// Padding   [49]byte
	Liquidity bin.Uint128 // 16 - 57
	// Padding   [65]byte
	SqrtPrice bin.Uint128 // 16 - 69
	// Padding          [81]byte
	TickCurrentIndex int32 // 4 - 73
	// Padding          [85]byte
	ProtocolFeeOwedA uint64 // 8 - 81
	// Padding          [93]byte
	ProtocolFeeOwedB uint64 // 8 - 89
	// Padding                    [101]byte
	TokenMintA                 *solana.PublicKey
	TokenVaultA                *solana.PublicKey
	FeeGrowthGlobalA           bin.Uint128
	TokenMintB                 *solana.PublicKey
	TokenVaultB                *solana.PublicKey
	FeeGrowthGlobalB           bin.Uint128
	RewardLastUpdatedTimestamp uint64
	RewardInfos                [3]WhirlpoolRewardInfoData
}

type WhirlpoolRewardInfoData struct {
	Mint                  solana.PublicKey
	Vault                 solana.PublicKey
	Authority             solana.PublicKey
	EmissionsPerSecondX64 bin.Uint128
	GrowthGlobalX64       bin.Uint128
}

type PositionData struct {
	Padding              [8]byte
	Whirlpool            *solana.PublicKey
	PositionMint         *solana.PublicKey
	Liquidity            bin.Uint128
	TickLowerIndex       int32
	TickUpperIndex       int32
	FeeGrowthCheckpointA bin.Uint128
	FeeOwedA             uint64
	FeeGrowthCeckpointB  bin.Uint128
	FeeOwedB             uint64
	RewardInfos          [3]PositionRewardInfoData
}

type PositionRewardInfoData struct {
	GrowthInsideCheckpoint bin.Uint128
	AmountOwed             uint64
}

type TickArrayData struct {
	Whirpool       *solana.PublicKey
	StartTickIndex int64
	Ticks          []TickData
}

type TickData struct {
	Initialized         bool
	LiquidityNet        bin.Uint128
	LiquidityGross      bin.Uint128
	FeeGrowthOutsideA   bin.Uint128
	FeeGrowthOutsideB   bin.Uint128
	RewardGrowthOutside []bin.Uint128
}

type WhirlpoolConfigData struct {
	FeeAuthority                 solana.PublicKey
	CollectProcolFeesAuthority   solana.PublicKey
	RewardEmissionSuperAuthority solana.PublicKey
	DefaultFeeRate               uint64
	DefaultProtocolFeeRate       uint64
}

type FeeTierData struct {
	WhirlpoolsConfig solana.PublicKey
	TickSpacing      uint64
	DefaultFeeRate   uint64
}
