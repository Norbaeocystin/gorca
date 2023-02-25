package gorca

import (
	"context"
	"fmt"
	whirlpool2 "github.com/Norbaeocystin/gorca/whirlpool"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
)

func OpenNewPositionAndAddLiquidity(tickLower,
	tickUpper int32, client rpc.Client, tokenAMax, tokenBMax uint64, liquidity bin.Uint128, tokenAAddress, tokenBAddress, TokenVaultA,
	TokenVaultB, WPPool solana.PublicKey, wallet solana.PrivateKey) (solana.Signature, error) {
	ktas := GetTickArrays(&client, WPPool)
	lowerArray := GetTickArray(tickLower, ktas)
	upperArray := GetTickArray(tickUpper, ktas)
	owner := wallet
	whirlpool2.ProgramID = ORCA_WHIRPOOL_PROGRAM_ID
	positionMint := solana.NewWallet()
	position, bump, _ := solana.FindProgramAddress([][]byte{[]byte("position"), positionMint.PublicKey().Bytes()},
		ORCA_WHIRPOOL_PROGRAM_ID)
	positionMA, bump2, _ := solana.FindProgramAddress([][]byte{[]byte("metadata"),
		ORCA_METADATA_PROGRAM_ID.Bytes(),
		positionMint.PublicKey().Bytes()},
		ORCA_METADATA_PROGRAM_ID)
	/*
		same as Token.getAssociatedTokenAddress(
						    ASSOCIATED_TOKEN_PROGRAM_ID,
						    TOKEN_PROGRAM_ID,
						    positionMintKeypair.publicKey,
						    owner
						  );
	*/
	positionTokenAccount, _, _ := solana.FindProgramAddress([][]byte{owner.PublicKey().Bytes(),
		solana.TokenProgramID.Bytes(),
		positionMint.PublicKey().Bytes(),
	}, solana.SPLAssociatedTokenAccountProgramID)
	i0 := solana.NewInstruction(COMPUTE_BUDGET,
		[]*solana.AccountMeta{},
		// fee 1, u
		[]uint8{0, 32, 161, 7, 0, 1, 0, 0, 0},
	)
	i1 := whirlpool2.NewOpenPositionWithMetadataInstruction(
		whirlpool2.OpenPositionWithMetadataBumps{bump, bump2},
		tickLower,
		tickUpper,
		owner.PublicKey(),
		owner.PublicKey(),
		position,
		positionMint.PublicKey(),
		positionMA,
		positionTokenAccount,
		WPPool,
		solana.TokenProgramID,
		solana.SystemProgramID,
		solana.SysVarRentPubkey,
		solana.SPLAssociatedTokenAccountProgramID,
		ORCA_METADATA_PROGRAM_ID,
		solana.MustPublicKeyFromBase58("3axbTs2z5GBy6usVbNVoqEgZMng3vZvMnAoX29BFfwhr"),
	).Build()
	i2 := whirlpool2.NewIncreaseLiquidityInstruction(
		liquidity,
		tokenAMax,
		tokenBMax,
		WPPool,
		solana.TokenProgramID,
		owner.PublicKey(),
		position,
		positionTokenAccount,
		tokenAAddress,
		tokenBAddress,
		TokenVaultA,
		TokenVaultB,
		lowerArray.Account,
		upperArray.Account,
	).Build()
	recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			i0,
			i1,
			i2,
		},
		recent.Value.Blockhash, //NONCE
		solana.TransactionPayer(owner.PublicKey()),
	)
	// log.Println(tx, err)
	// TODO intiliaze those 2 accounts
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if wallet.PublicKey().Equals(key) {
				return &wallet
			}
			//if position.PublicKey().Equals(key) {
			//	return &position.PrivateKey
			//}
			if positionMint.PublicKey().Equals(key) {
				return &positionMint.PrivateKey
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
	log.Println(position, positionMint.PublicKey(), positionTokenAccount, positionMA)
	return sig, err
}
