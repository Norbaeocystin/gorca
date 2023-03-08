package gorca

import (
	"github.com/Norbaeocystin/gorca/whirlpool"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Market struct {
	ProgramId     solana.PublicKey
	MarketId      solana.PublicKey
	WhirlpoolData WhirlpoolData
	Oracle        solana.PublicKey
	KTAS          []KeyedTickArray
	Client        *rpc.Client
}

func (m *Market) SetOracle() {
	m.Oracle, _, _ = GetOracle(m.ProgramId, m.MarketId)
}

func (m Market) FetchKtas() []KeyedTickArray {
	ktas := GetTickArrays(m.Client, m.ProgramId, m.MarketId)
	return ktas
}

func (m *Market) FetchAndSetKtas() {
	m.KTAS = m.FetchKtas()
}

func (m Market) SwapAtoBExactInputInstruction(amount, otherAmountThreshold uint64, sqrtPriceLimit bin.Uint128, owner, ownerTokenAAddress, ownerTokenBAddress, kta0, kta1, kta2 solana.PublicKey) solana.Instruction {
	return whirlpool.NewSwapInstruction(
		amount,
		otherAmountThreshold,
		sqrtPriceLimit,
		true,
		true,
		solana.TokenProgramID,
		owner,
		m.MarketId,
		ownerTokenAAddress,
		*m.WhirlpoolData.TokenVaultA,
		ownerTokenBAddress,
		*m.WhirlpoolData.TokenVaultB,
		kta0,
		kta1,
		kta2,
		m.Oracle,
	).Build()
}

func (m Market) SwapAtoBExactOutputInstruction(amount, otherAmountThreshold uint64, sqrtPriceLimit bin.Uint128, owner, ownerTokenAAddress, ownerTokenBAddress, kta0, kta1, kta2 solana.PublicKey) solana.Instruction {
	return whirlpool.NewSwapInstruction(
		amount,
		otherAmountThreshold,
		sqrtPriceLimit,
		false,
		true,
		solana.TokenProgramID,
		owner,
		m.MarketId,
		ownerTokenAAddress,
		*m.WhirlpoolData.TokenVaultA,
		ownerTokenBAddress,
		*m.WhirlpoolData.TokenVaultB,
		kta0,
		kta1,
		kta2,
		m.Oracle,
	).Build()
}

func (m Market) SwapBToAExactInputInstruction(amount, otherAmountThreshold uint64, sqrtPriceLimit bin.Uint128, owner, ownerTokenAAddress, ownerTokenBAddress, kta0, kta1, kta2 solana.PublicKey) solana.Instruction {
	return whirlpool.NewSwapInstruction(
		amount,
		otherAmountThreshold,
		sqrtPriceLimit,
		true,
		false,
		solana.TokenProgramID,
		owner,
		m.MarketId,
		ownerTokenAAddress,
		*m.WhirlpoolData.TokenVaultA,
		ownerTokenBAddress,
		*m.WhirlpoolData.TokenVaultB,
		kta0,
		kta1,
		kta2,
		m.Oracle,
	).Build()
}

func (m Market) SwapBToAExactOutputInstruction(amount, otherAmountThreshold uint64, sqrtPriceLimit bin.Uint128, owner, ownerTokenAAddress, ownerTokenBAddress, kta0, kta1, kta2 solana.PublicKey) solana.Instruction {
	return whirlpool.NewSwapInstruction(
		amount,
		otherAmountThreshold,
		sqrtPriceLimit,
		false,
		false,
		solana.TokenProgramID,
		owner,
		m.MarketId,
		ownerTokenAAddress,
		*m.WhirlpoolData.TokenVaultA,
		ownerTokenBAddress,
		*m.WhirlpoolData.TokenVaultB,
		kta0,
		kta1,
		kta2,
		m.Oracle,
	).Build()
}
