package gorcagithub

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	whirlpool2 "gorca/whirlpool"
)

func Update(client *rpc.Client, whirlpoolAddress, position solana.PublicKey, owner solana.PrivateKey, positionLowerTick, positionUpperTick int32) solana.Signature {
	whirlpool2.ProgramID = ORCA_WHIRPOOL_PROGRAM_ID
	ktas := GetTickArrays(client, whirlpoolAddress)
	lowerArray := GetTickArray(positionLowerTick, ktas)
	upperArray := GetTickArray(positionUpperTick, ktas)
	inst := whirlpool2.NewUpdateFeesAndRewardsInstruction(
		whirlpoolAddress,
		position,
		lowerArray.Account,
		upperArray.Account,
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
