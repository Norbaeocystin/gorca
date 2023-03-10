package gorca

import (
	"context"
	"fmt"
	whirlpool2 "github.com/Norbaeocystin/gorca/whirlpool"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
)

func BurnPosition(client rpc.Client, position, positionMint, positionTokenAccount solana.PublicKey,
	owner solana.PrivateKey) solana.Signature {
	log.Println("sending burn/close transactions")
	whirlpool2.ProgramID = ORCA_WHIRPOOL_PROGRAM_ID
	inst := whirlpool2.NewClosePositionInstruction(
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
			inst,
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
