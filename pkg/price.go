package gorca

import (
	"log"
	"math"
	"math/big"
)

var Q64 = new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil))

func EncodeSqrtRatioX64(amount1 *big.Int, amount0 *big.Int) *big.Int {
	numerator := new(big.Int).Lsh(amount1, 128)
	denominator := amount0
	ratioX128 := new(big.Int).Div(numerator, denominator)
	return new(big.Int).Sqrt(ratioX128)
}

// calculate squared root price with Q96
// example of input 0.86 matic/usdc what is 0.86e-12
func CalculateSqrtPriceQ64(price *big.Float) *big.Int {
	priceSqrt := new(big.Float).Sqrt(price)
	sqrtQ64Price := new(big.Float).Mul(priceSqrt, Q64)
	result := new(big.Int)
	sqrtQ64Price.Int(result)
	return result
}

func CalculatePriceFromSQRPriceQ64(priceSQRTQ64 *big.Int) float64 {
	floatPriceSQRTQ64 := new(big.Float).SetInt(priceSQRTQ64)
	floatPriceSQRT := new(big.Float).Quo(floatPriceSQRTQ64, Q64)
	floatPrice := new(big.Float).Mul(floatPriceSQRT, floatPriceSQRT)
	price, _ := floatPrice.Float64()
	return price
}

func GetTicksForPrice(price *big.Float, tick, tickSpacing *big.Int, spreadPCT *big.Float) (*big.Int, *big.Int) {
	higherPCT := new(big.Float).Quo(new(big.Float).Add(big.NewFloat(100), spreadPCT), big.NewFloat(100))
	lowerPCT := new(big.Float).Quo(new(big.Float).Sub(big.NewFloat(100), spreadPCT), big.NewFloat(100))
	// log.Println("pct", lowerPCT.String(), higherPCT.String())
	// log.Println(higherPCT, lowerPCT)
	priceHigher := new(big.Float).Mul(price, higherPCT)
	priceLower := new(big.Float).Mul(price, lowerPCT)
	log.Println("pl", "ph", priceLower, priceHigher)
	tickNormalizedFee := new(big.Int).Mul(new(big.Int).Div(tick, tickSpacing), tickSpacing)
	// log.Println("normalized tick", tickNormalizedFee.String())
	lower := tickNormalizedFee
	// log.Println("lower", tickNormalizedFee)
	for {
		TICK := lower.Int64()
		TICKfloat := float64(TICK)
		// time.Sleep(5 * time.Second)
		price := math.Pow(1.0001, TICKfloat)
		priceCalculated := big.NewFloat(price)
		if priceLower.Cmp(priceCalculated) == 1 {
			break
		}
		lower.Sub(lower, tickSpacing)
	}
	tickNormalizedFee = new(big.Int).Mul(new(big.Int).Div(tick, tickSpacing), tickSpacing)
	higher := tickNormalizedFee
	// log.Println("higher", tickNormalizedFee)
	for {
		TICK := higher.Int64()
		TICKfloat := float64(TICK)
		price := math.Pow(1.0001, TICKfloat)
		priceCalculated := big.NewFloat(price)
		// log.Println("ph", "pc", priceHigher, priceCalculated)
		if priceHigher.Cmp(priceCalculated) == -1 {
			break
		}
		// log.Println(higher, priceHigher, priceCalculated, price)
		higher.Add(higher, tickSpacing)
	}
	return lower, higher
}
