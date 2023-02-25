package gorcagithub

import (
	"context"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
	"math/big"
)

func FindATA(wallet solana.PrivateKey, mint solana.PublicKey) solana.PublicKey {
	ataAccount, _, _ := solana.FindAssociatedTokenAddress(wallet.PublicKey(), mint)
	return ataAccount
}

func GetBalance(client *rpc.Client, tokenAddress solana.PublicKey) (uint64, error) {
	result, err := client.GetTokenAccountBalance(context.TODO(), tokenAddress, rpc.CommitmentSingleGossip)
	if err != nil {
		log.Println("fetching balance returned", err)
		return 1, err
	}
	amount, ok := new(big.Int).SetString(result.Value.Amount, 10)
	log.Println("got amount which can be used for liquidation", amount, ok)
	liquidityAmount := amount.Uint64()
	return liquidityAmount, nil
}

func GetTokenBalance(client *rpc.Client, tokenPublicMint, owner solana.PublicKey) (solana.PublicKey, solana.PublicKey,
	uint64) {
	tokens, _ := client.GetTokenAccountsByOwner(context.TODO(), owner,
		&rpc.GetTokenAccountsConfig{
			Mint:      nil,
			ProgramId: solana.TokenProgramID.ToPointer(),
		},
		&rpc.GetTokenAccountsOpts{
			Commitment: "",
			Encoding:   solana.EncodingBase64,
			DataSlice:  nil,
		})
	for _, tk := range tokens.Value {
		var ta token.Account
		borshDec := bin.NewBorshDecoder(tk.Account.Data.GetBinary())
		// log.Println(tk.Pubkey.String(), ta.Mint)
		borshDec.Decode(&ta)
		if ta.Mint == tokenPublicMint {
			return tk.Pubkey, ta.Mint, ta.Amount
		}
		return solana.PublicKey{}, solana.PublicKey{}, 0
	}
	return solana.PublicKey{}, solana.PublicKey{}, 0
}
