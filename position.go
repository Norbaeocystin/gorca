package gorca

import (
	"context"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"log"
)

func FindOrcaPositionsForOwner(client *rpc.Client, owner, pool solana.PublicKey) ([]PositionData, error) {
	positions := make([]PositionData, 0)
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
			pk, _, _ := GetPosition(ORCA_WHIRPOOL_PROGRAM_ID, ta.Mint)
			position := GetPositionData(client, pk)
			if position.Whirlpool != nil && position.Whirlpool.Equals(pool) {
				positions = append(positions, position)
			}
		}
	}
	return positions, nil
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
