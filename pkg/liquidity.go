package gorca

import (
	"context"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"gorca/pkg/whirlpool"
	"log"
	"math"
	"math/big"
)

var Q64BI = new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)

func CalculateLiquidity(tokenAMax, tokenBMax uint64, tick, upperTick, lowerTick int64) *big.Int {
	low := CalculateSqrtPriceQ64(big.NewFloat(math.Pow(1.0001, float64(lowerTick))))
	current := CalculateSqrtPriceQ64(big.NewFloat(math.Pow(1.0001, float64(tick))))
	upper := CalculateSqrtPriceQ64(big.NewFloat(math.Pow(1.0001, float64(upperTick))))
	var diff0 *big.Int
	if current.Cmp(upper) == 1 {
		diff0 = new(big.Int).Sub(current, upper)
	} else {
		diff0 = new(big.Int).Sub(upper, current)
	}
	var diff1 *big.Int
	if current.Cmp(low) == 1 {
		diff1 = new(big.Int).Sub(current, low)
	} else {
		diff1 = new(big.Int).Sub(low, current)
	}
	liq0 := new(big.Int).Div(new(big.Int).Div(new(big.Int).Mul(big.NewInt(int64(tokenAMax)), new(big.Int).Mul(current, upper)), Q64BI), diff0)
	liq1 := new(big.Int).Div(new(big.Int).Mul(big.NewInt(int64(tokenBMax)), Q64BI), diff1)
	if liq0.Cmp(liq1) == -1 {
		return liq0
	}
	return liq1
}

func IncreaseLiquidityWSOL(client *rpc.Client, tokenAMax, tokenBMax uint64,
	position, positionTokenAccount, tokenAAddress, tokenBAddress, whirlpoolAddress solana.PublicKey,
	owner solana.PrivateKey, positionLowerTick, positionUpperTick int32) solana.Signature {
	// (token A * token B)^(1/2) / (pool token supply)
	whirlpool.ProgramID = ORCA_WHIRPOOL_PROGRAM_ID
	ktas := GetTickArrays(client, whirlpoolAddress)
	lowerArray := GetTickArray(positionLowerTick, ktas)
	upperArray := GetTickArray(positionUpperTick, ktas)
	wpPoolData := GetWhirlpoolData(*client, whirlpoolAddress)
	liq := LiquidityInfiniteCurve(tokenAMax, tokenBMax)
	liquiditySupply, err := BigIntToBinUint128(liq)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(liquiditySupply)
	i3 := whirlpool.NewIncreaseLiquidityInstruction(
		liquiditySupply,
		tokenAMax,
		tokenBMax,
		whirlpoolAddress,
		solana.TokenProgramID,
		owner.PublicKey(),
		position,
		positionTokenAccount,
		tokenAAddress,
		tokenBAddress,
		*wpPoolData.TokenVaultA,
		*wpPoolData.TokenVaultB,
		lowerArray.Account,
		upperArray.Account,
	).Build()
	//i4 := token.NewCloseAccountInstruction(
	//	newWSOL.PublicKey(),
	//	owner.PublicKey(),
	//	owner.PublicKey(),
	//	nil, //[]solana.PublicKey{newWSOL.PublicKey(), owner.PublicKey()},
	//).Build()
	recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			i3, // i1, i2, i3, i4,
		},
		recent.Value.Blockhash, //NONCE
		solana.TransactionPayer(owner.PublicKey()),
	)
	// log.Println(tx, err)
	// TODO intiliaze those 2 accounts
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if owner.PublicKey().Equals(key) {
				return &owner
			}
			return nil
		},
	)
	if err != nil {
		panic(fmt.Errorf("unable to sign transaction: %w", err))
	}
	sig, err := client.SendTransactionWithOpts(context.TODO(), tx,
		rpc.TransactionOpts{
			Encoding:            "",
			SkipPreflight:       false,
			PreflightCommitment: "",
			MaxRetries:          nil,
			MinContextSlot:      nil,
		},
	)
	if err != nil {
		panic(err)
	}
	// log.Println(sig)
	return sig
}

func DecreaseUpdateCollectBurn(client *rpc.Client, liquidity bin.Uint128, tokenAMin, tokenBMin uint64,
	position, positionMint, positionTokenAccount, tokenAAddress, tokenBAddress, tokenVaultA, tokenVaultB, whirlpoolAddress solana.PublicKey,
	owner solana.PrivateKey, positionLowerTick, positionUpperTick int32) (solana.Signature, error) {
	whirlpool.ProgramID = ORCA_WHIRPOOL_PROGRAM_ID
	ktas := GetTickArrays(client, whirlpoolAddress)
	lowerArray := GetTickArray(positionLowerTick, ktas)
	upperArray := GetTickArray(positionUpperTick, ktas)
	GetPositionData(client, position)
	// lamports := 2000000 + tokenAMax
	// DECREASE LIQUIDITY after update fees
	i0 := solana.NewInstruction(COMPUTE_BUDGET,
		[]*solana.AccountMeta{},
		// fee 1, u
		[]uint8{0, 32, 161, 7, 0, 1, 0, 0, 0},
	)
	i2 := whirlpool.NewDecreaseLiquidityInstruction(
		liquidity,
		tokenAMin,
		tokenBMin,
		whirlpoolAddress,
		solana.TokenProgramID,
		owner.PublicKey(),
		position,
		positionTokenAccount,
		tokenAAddress,
		tokenBAddress,
		tokenVaultA,
		tokenVaultB,
		lowerArray.Account, // tickarray,
		upperArray.Account, //tickarray,
	).Build()
	// UPDATE
	i1 := whirlpool.NewUpdateFeesAndRewardsInstruction(
		whirlpoolAddress,
		position,
		lowerArray.Account, // tickarray,
		upperArray.Account, //tickarray,
	).Build()
	// COLLECT FEES
	i3 := whirlpool.NewCollectFeesInstruction(
		whirlpoolAddress,
		owner.PublicKey(),
		position,
		positionTokenAccount,
		tokenAAddress,
		tokenVaultA,
		tokenBAddress,
		tokenVaultB,
		solana.TokenProgramID,
	).Build()
	// BURN
	i4 := whirlpool.NewClosePositionInstruction(
		owner.PublicKey(),
		owner.PublicKey(),
		position,
		positionMint,
		positionTokenAccount,
		solana.TokenProgramID,
	).Build()
	recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			// first updatefees after removing liquidity
			i0, i1, i2, i3, i4,
		},
		recent.Value.Blockhash, //NONCE
		solana.TransactionPayer(owner.PublicKey()),
	)
	// log.Println(tx, err)
	// TODO intiliaze those 2 accounts
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if owner.PublicKey().Equals(key) {
				return &owner
			}
			return nil
		},
	)
	if err != nil {
		panic(fmt.Errorf("unable to sign transaction: %w", err))
	}
	sig, err := client.SendTransactionWithOpts(context.TODO(), tx,
		rpc.TransactionOpts{
			Encoding:            "",
			SkipPreflight:       false,
			PreflightCommitment: "",
			MaxRetries:          nil,
			MinContextSlot:      nil,
		},
	)
	// log.Println(sig)
	return sig, err
}
