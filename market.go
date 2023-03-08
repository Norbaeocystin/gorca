package gorca

import (
	"github.com/Norbaeocystin/gorca/whirlpool"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"math/big"
)

func NewMarket(client *rpc.Client, programId, marketId solana.PublicKey) Market {
	var m Market
	m.Client = client
	m.ProgramId = programId
	m.MarketId = marketId
	m.SetKtas()
	m.SetData()
	m.SetOracle()
	return m
}

type Market struct {
	ProgramId     solana.PublicKey
	MarketId      solana.PublicKey
	WhirlpoolData WhirlpoolData
	Oracle        solana.PublicKey
	KTAS          []KeyedTickArray
	Client        *rpc.Client
}

func (m Market) FetchData() WhirlpoolData {
	data := GetWhirlpoolData(m.Client, m.MarketId)
	return data
}

func (m *Market) SetData() {
	m.WhirlpoolData = m.FetchData()
}

func (m *Market) SetOracle() {
	m.Oracle, _, _ = GetOracle(m.ProgramId, m.MarketId)
}

func (m Market) FetchKtas() []KeyedTickArray {
	ktas := GetTickArrays(m.Client, m.ProgramId, m.MarketId)
	return ktas
}

func (m *Market) SetKtas() {
	m.KTAS = m.FetchKtas()
}

func (m Market) GetKtasForTicks(currentTick, tickForSqrtPriceLimit int32) (solana.PublicKey, solana.PublicKey, solana.PublicKey) {
	keyTickForSqrtPriceLimit := GetTickArray(tickForSqrtPriceLimit, m.KTAS).Account
	return GetTickArray(currentTick, m.KTAS).Account, keyTickForSqrtPriceLimit, keyTickForSqrtPriceLimit
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

func (m Market) SwapAtoBExactInputInstructionWithSlippageUseWPData(amount uint64, slippagePCT float64, owner, ownerTokenAAddress, ownerTokenBAddress solana.PublicKey) solana.Instruction {
	price := CalculatePriceFromSQRPriceQ64(m.WhirlpoolData.SqrtPrice.BigInt())
	priceWithSlippage := price - ((price * 100) / slippagePCT)
	tick := PriceToTick(priceWithSlippage)
	tickNormalized := (tick / int32(m.WhirlpoolData.TickSpacing)) * int32(m.WhirlpoolData.TickSpacing)
	otherAmountThreshold := uint64(float64(amount) * priceWithSlippage)
	sqrtPriceLimit, _ := BigIntToBinUint128(CalculateSqrtPriceQ64(big.NewFloat(priceWithSlippage)))
	kta0, kta1, kta2 := m.GetKtasForTicks(m.WhirlpoolData.TickCurrentIndex, tickNormalized)
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

func (m Market) SwapBtoAExactInputInstructionWithSlippageUseWPData(amount uint64, slippagePCT float64, owner, ownerTokenAAddress, ownerTokenBAddress solana.PublicKey) solana.Instruction {
	price := CalculatePriceFromSQRPriceQ64(m.WhirlpoolData.SqrtPrice.BigInt())
	priceWithSlippage := price + ((price * 100) / slippagePCT)
	tick := PriceToTick(priceWithSlippage)
	tickNormalized := (tick / int32(m.WhirlpoolData.TickSpacing)) * int32(m.WhirlpoolData.TickSpacing)
	otherAmountThreshold := uint64(float64(amount) / priceWithSlippage)
	sqrtPriceLimit, _ := BigIntToBinUint128(CalculateSqrtPriceQ64(big.NewFloat(priceWithSlippage)))
	kta0, kta1, kta2 := m.GetKtasForTicks(m.WhirlpoolData.TickCurrentIndex, tickNormalized)
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
