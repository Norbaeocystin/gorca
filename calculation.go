package gorcagithub

import "math/big"

func LiquidityInfiniteCurve(amountA, amountB uint64) *big.Int {
	MUL := new(big.Int).Mul(big.NewInt(int64(amountA)), big.NewInt(int64(amountB)))
	mul := new(big.Int)
	SQRT := new(big.Float).Sqrt(new(big.Float).SetInt(MUL))
	SQRT.Int(mul)
	return mul
}
