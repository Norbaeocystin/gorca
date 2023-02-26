package gorca

import (
	"context"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func GetWhirlpoolData(client *rpc.Client, whirlpoolAddress solana.PublicKey) WhirlpoolData {
	account, _ := client.GetAccountInfoWithOpts(context.TODO(),
		whirlpoolAddress,
		&rpc.GetAccountInfoOpts{
			Encoding:       solana.EncodingBase64,
			Commitment:     rpc.CommitmentFinalized,
			DataSlice:      nil,
			MinContextSlot: nil,
		},
	)
	var wpData WhirlpoolData
	dataPos := account.GetBinary()
	borshDec := bin.NewBorshDecoder(dataPos)
	borshDec.Decode(&wpData)
	// log.Println(wpData)
	return wpData
}
