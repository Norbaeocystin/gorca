package gorca

import "github.com/gagliardetto/solana-go"

var ORCA_WHIRPOOL_PROGRAM_ID = solana.MustPublicKeyFromBase58("whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc")
var ORCA_METADATA_PROGRAM_ID = solana.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

var WSOL = "So11111111111111111111111111111111111111112"
var USDC = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
var WHIRPOOL_WSOLUSDC_005 = "7qbRF6YsyGuLUVs6Y1q64bdVrfe4ZcUUz1JRdoVNUJnm"
var WHIRPOOL_WSOLUSDC_03 = "HJPjoWUrhoZzkNfRpHuieeFk9WcZWjwy6PBjZ81ngndJ"
var ORACLE_WSOLUSDC_03 = solana.MustPublicKeyFromBase58("4GkRbcYg1VKsZropgai4dMf2Nj2PkXNLf43knFpavrSi")
var ORACLE_WSOLUSDC_005 = solana.MustPublicKeyFromBase58("6vK8gSiRHSnZzAa5JsvBF2ej1LrxpRX21Y185CzP4PeA")
var COMPUTE_BUDGET = solana.MustPublicKeyFromBase58("ComputeBudget111111111111111111111111111111")
var API_WHILPOOLS = "https://api.mainnet.orca.so/v1/whirlpool/list"

const TICK_ARRAY_SIZE = 88
const MIN_TICK_INDEX = -443636
const MAX_TICK_INDEX = 443636
