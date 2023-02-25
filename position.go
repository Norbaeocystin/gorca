package gorcagithub

import (
	"context"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
	"strings"
)

type PositionKeys struct {
	PositionMint         solana.PublicKey
	Position             solana.PublicKey
	PositionMetadata     solana.PublicKey
	PositionTokenAccount solana.PublicKey
	WhirpoolAddress      solana.PublicKey
	OrcaPosition         bool
}

func FindOrcaPositionsForOwner(client *rpc.Client, owner solana.PublicKey) ([]PositionKeys, error) {
	positionsKeys := make([]PositionKeys, 0)
	tokens, err := client.GetTokenAccountsByOwner(context.TODO(), owner,
		&rpc.GetTokenAccountsConfig{
			Mint:      nil,
			ProgramId: solana.TokenProgramID.ToPointer(),
		},

		&rpc.GetTokenAccountsOpts{
			Commitment: "",
			Encoding:   solana.EncodingBase64,
			DataSlice:  nil,
		})
	if err != nil {
		log.Println("fetching positions error", err)
	}
	for _, tk := range tokens.Value {
		var ta token.Account
		borshDec := bin.NewBorshDecoder(tk.Account.Data.GetBinary())
		// log.Println(tk.Pubkey.String(), ta.Mint)
		borshDec.Decode(&ta)
		// log.Println(tk.Pubkey.String(), ta.Mint, len(tk.Account.Data.GetBinary()))
		if ta.Amount == 1 {
			// tkPubkey - positionTokenAccount, ta.Mint - positionMint
			// log.Println(tk.Pubkey, ta.Mint)
			keys, err := FindPublicKeysForPositioMint(client, ta.Mint)
			if err != nil {
				return positionsKeys, err
			}
			// var pubKey solana.PublicKey
			if keys.OrcaPosition == true {
				positionsKeys = append(positionsKeys, keys)
			}
		}
	}
	return positionsKeys, nil
}

// not working if use with SOL, position needs to be open with WSOL
func FindPublicKeysForPositioMint(client *rpc.Client, positionMint solana.PublicKey) (PositionKeys, error) {
	var positionKeys PositionKeys
	positionKeys.OrcaPosition = false
	one := 1
	signature, err := client.GetSignaturesForAddressWithOpts(context.TODO(), positionMint,
		&rpc.GetSignaturesForAddressOpts{
			&one,
			solana.Signature{},
			solana.Signature{},
			rpc.CommitmentFinalized,
			nil,
		})
	if err != nil {
		return positionKeys, err
	}
	version := uint64(0)
	opts := rpc.GetTransactionOpts{
		Encoding:                       solana.EncodingBase64,
		Commitment:                     rpc.CommitmentFinalized,
		MaxSupportedTransactionVersion: &version,
	}
	// log.Println("sign", signature[0].Signature)

	txs, err := client.GetTransaction(context.TODO(), signature[0].Signature, &opts)
	if err != nil {
		log.Println(err)
		return positionKeys, err
	}
	if strings.Contains(strings.Join(txs.Meta.LogMessages, ","), ORCA_WHIRPOOL_PROGRAM_ID.String()) {
		tx, _ := txs.Transaction.GetTransaction()
		accounts, _ := tx.AccountMetaList()
		if len(accounts) > 8 {
			// log.Println("Match tx with signature", signature[0].Signature)
			positionKeys.Position = accounts[2].PublicKey
			positionKeys.PositionTokenAccount = accounts[4].PublicKey
			positionKeys.PositionMetadata = accounts[3].PublicKey
			positionKeys.WhirpoolAddress = accounts[5].PublicKey
			positionKeys.PositionMint = positionMint
			positionKeys.OrcaPosition = true
		}
	}
	return positionKeys, nil
}

func GetPositionData(client *rpc.Client, position solana.PublicKey) PositionData {
	account, _ := client.GetAccountInfoWithOpts(context.TODO(),
		position,
		&rpc.GetAccountInfoOpts{
			Encoding:       solana.EncodingBase64,
			Commitment:     rpc.CommitmentFinalized,
			DataSlice:      nil,
			MinContextSlot: nil,
		},
	)
	var positionData PositionData
	dataPos := account.GetBinary()
	borshDec := bin.NewBorshDecoder(dataPos)
	borshDec.Decode(&positionData)
	// log.Println(positionData)
	return positionData
}
