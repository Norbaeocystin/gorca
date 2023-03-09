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
	// TOO COSTLY
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
	KTAS          KTAS
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

// calculate derived key
func (m Market) GetTickAccount(tick int32) solana.PublicKey {
	startTick := GetStartTickIndex(tick, m.WhirlpoolData.TickSpacing)
	key, _, _ := GetTickArrayViaFPA(m.ProgramId, m.MarketId, startTick)
	return key
}

func (m Market) SwapAtoBExactInputInstruction(amount, otherAmountThreshold uint64, sqrtPriceLimit bin.Uint128, owner, ownerTokenAAddress, ownerTokenBAddress, kta0, kta1, kta2 solana.PublicKey) solana.Instruction {
	whirlpool.ProgramID = m.ProgramId
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
	whirlpool.ProgramID = m.ProgramId
	price := CalculatePriceFromSQRPriceQ64(m.WhirlpoolData.SqrtPrice.BigInt())
	priceWithSlippage := price - (price * (slippagePCT / 100))
	tick := PriceToTick(priceWithSlippage)
	otherAmountThreshold := uint64(float64(amount) * priceWithSlippage)
	sqrtPriceLimit, _ := BigIntToBinUint128(CalculateSqrtPriceQ64(big.NewFloat(priceWithSlippage)))
	kta0 := m.GetTickAccount(m.WhirlpoolData.TickCurrentIndex)
	kta1 := m.GetTickAccount(tick)
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
		kta1,
		m.Oracle,
	).Build()
}

func (m Market) SwapAtoBExactInputInstructionWithSlippageUsePrice(amount uint64, price, slippagePCT float64, owner, ownerTokenAAddress, ownerTokenBAddress solana.PublicKey) solana.Instruction {
	whirlpool.ProgramID = m.ProgramId
	currentTick := (PriceToTick(price) / int32(m.WhirlpoolData.TickSpacing)) * int32(m.WhirlpoolData.TickSpacing)
	priceWithSlippage := price - (price * (slippagePCT / 100))
	tick := PriceToTick(priceWithSlippage)
	// tickNormalized := (tick / int32(m.WhirlpoolData.TickSpacing)) * int32(m.WhirlpoolData.TickSpacing)
	otherAmountThreshold := uint64(float64(amount) * priceWithSlippage)
	sqrtPriceLimit, _ := BigIntToBinUint128(CalculateSqrtPriceQ64(big.NewFloat(priceWithSlippage)))
	kta0 := m.GetTickAccount(currentTick)
	kta1 := m.GetTickAccount(tick)
	// kta0, kta1, kta2 := m.GetKtasForTicks(currentTick, tickNormalized)
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
		kta1,
		m.Oracle,
	).Build()
}

func (m Market) SwapAtoBExactOutputInstruction(amount, otherAmountThreshold uint64, sqrtPriceLimit bin.Uint128, owner, ownerTokenAAddress, ownerTokenBAddress, kta0, kta1, kta2 solana.PublicKey) solana.Instruction {
	whirlpool.ProgramID = m.ProgramId
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
	whirlpool.ProgramID = m.ProgramId
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
	whirlpool.ProgramID = m.ProgramId
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
	whirlpool.ProgramID = m.ProgramId
	price := CalculatePriceFromSQRPriceQ64(m.WhirlpoolData.SqrtPrice.BigInt())
	priceWithSlippage := price + (price * (slippagePCT / 100))
	tick := PriceToTick(priceWithSlippage)
	otherAmountThreshold := uint64(float64(amount) / priceWithSlippage)
	sqrtPriceLimit, _ := BigIntToBinUint128(CalculateSqrtPriceQ64(big.NewFloat(priceWithSlippage)))
	kta0 := m.GetTickAccount(m.WhirlpoolData.TickCurrentIndex)
	kta1 := m.GetTickAccount(tick)
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
		kta1,
		m.Oracle,
	).Build()
}

// price needs to be modified to  (price * decimalsB)/decimalsA
func (m Market) SwapBtoAExactInputInstructionWithSlippageUsePrice(amount uint64, price, slippagePCT float64, owner, ownerTokenAAddress, ownerTokenBAddress solana.PublicKey) solana.Instruction {
	whirlpool.ProgramID = m.ProgramId
	currentTick := (PriceToTick(price) / int32(m.WhirlpoolData.TickSpacing)) * int32(m.WhirlpoolData.TickSpacing)
	priceWithSlippage := price + (price * (slippagePCT / 100))
	tick := PriceToTick(priceWithSlippage)
	otherAmountThreshold := uint64(float64(amount) / priceWithSlippage)
	sqrtPriceLimit, _ := BigIntToBinUint128(CalculateSqrtPriceQ64(big.NewFloat(priceWithSlippage)))
	kta0 := m.GetTickAccount(currentTick)
	kta1 := m.GetTickAccount(tick)
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
		kta1,
		m.Oracle,
	).Build()
}
