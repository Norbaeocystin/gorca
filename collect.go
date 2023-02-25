package gorcagithub

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	whirlpool2 "gorca/whirlpool"
)

func Collect(client *rpc.Client, whirlpoolAddress, position, positionTokenAccount,
	tokenOwnerAccountA, tokenVaultA, tokenOwnerAccountB, tokenVaultB solana.PublicKey, owner solana.PrivateKey) solana.Signature {
	whirlpool2.ProgramID = ORCA_WHIRPOOL_PROGRAM_ID
	i := whirlpool2.NewCollectFeesInstruction(
		whirlpoolAddress,
		owner.PublicKey(),
		position,
		positionTokenAccount,
		tokenOwnerAccountA,
		tokenVaultA,
		tokenOwnerAccountB,
		tokenVaultB,
		solana.TokenProgramID,
	).Build()
	recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			i,
		},
		recent.Value.Blockhash, //NONCE
		solana.TransactionPayer(owner.PublicKey()),
	)
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
