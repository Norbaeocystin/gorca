package gorca

import (
	"context"
	"fmt"
	whirlpool2 "github.com/Norbaeocystin/gorca/whirlpool"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func Swap(client *rpc.Client, amountIn, amountOutMin uint64, sqrtPriceLimit bin.Uint128, amountSpecifiedIsInput, AtoB bool, tokenA, tokenB,
	tokenAVault, tokenBVault, oracle, wpPool solana.PublicKey, tickArrays [2]solana.PublicKey, owner solana.PrivateKey) (solana.Signature, error) {
	whirlpool2.ProgramID = ORCA_WHIRPOOL_PROGRAM_ID
	i0 := solana.NewInstruction(COMPUTE_BUDGET,
		[]*solana.AccountMeta{},
		// fee 1, u
		[]uint8{0, 32, 161, 7, 0, 1, 0, 0, 0},
	)

	i := whirlpool2.NewSwapInstruction(
		amountIn,
		amountOutMin,
		sqrtPriceLimit,
		amountSpecifiedIsInput,
		AtoB,
		solana.TokenProgramID,
		owner.PublicKey(),
		wpPool,
		tokenA,
		tokenAVault,
		tokenB,
		tokenBVault,
		tickArrays[0],
		tickArrays[0],
		tickArrays[1],
		oracle,
	).Build()
	recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			i0, i,
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
