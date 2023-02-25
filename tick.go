package gorcagithub

import (
	"context"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"gorca/whirlpool"
	"log"
	"sort"
)

/*
tick = 1 + 16 + 16 + 16 + 16 + 48
tick = 113

TICK_ARRAY_SIZE i32 = 88
tick array_size_usize = 88

9944 + 8 = 9952
tickarray = 8 + 36 + (113 * 88 )
tickarray = 9988
*/

type KTAS []KeyedTickArray

func (ta KTAS) Len() int { return len(ta) }
func (ta KTAS) Less(i, j int) bool {
	return ta[i].TickArray.StartTickIndex < ta[j].TickArray.StartTickIndex
}
func (ta KTAS) Swap(i, j int) { ta[i], ta[j] = ta[j], ta[i] }

type KeyedTickArray struct {
	Account   solana.PublicKey
	TickArray whirlpool.TickArray
}

func GetTickArrays(client *rpc.Client, market solana.PublicKey) []KeyedTickArray {
	memcmp := rpc.RPCFilterMemcmp{9956, market.Bytes()}
	//var one uint8
	///// two, three, four uint8
	//one = 1
	//two = 2
	//three = 3
	//four = 4
	filters := []rpc.RPCFilter{
		{ // Memcmp: &memcmp,
			DataSize: uint64(9988),
		},
		{
			Memcmp: &memcmp,
		},
		//{
	}
	opts := rpc.GetProgramAccountsOpts{
		Commitment: "",
		Encoding:   "base64",
		Filters:    filters,
	}
	out, err := client.GetProgramAccountsWithOpts(
		context.TODO(),
		ORCA_WHIRPOOL_PROGRAM_ID,
		&opts,
	)
	if err != nil {
		log.Println("got error during fetching tickArrays", err)
	}
	ktas := make(KTAS, 0)
	for _, acc := range out {
		var kta KeyedTickArray
		var ta whirlpool.TickArray
		b := acc.Account.Data.GetBinary()
		decoder := bin.NewBorshDecoder(b)
		decoder.Decode(&ta)
		kta.TickArray = ta
		kta.Account = acc.Pubkey
		ktas = append(ktas, kta)
	}
	sort.Sort(ktas)
	log.Println("got tickarrays:", len(ktas))
	return ktas
}

func GetStartTickIndex() {

}

func GetTickArray(tick int32, ktas KTAS) KeyedTickArray {
	for idx, kta := range ktas {
		if kta.TickArray.StartTickIndex > tick {
			return ktas[idx-1]
		}

	}
	return ktas[len(ktas)-1]
}
