package gorca

import (
	"encoding/json"
	"github.com/gagliardetto/solana-go"
	"io/ioutil"
	"net/http"
	"sort"
)

// get whirlpools data via orca api
func GetWhirlpoolsViaApi() (WhirlpoolsApi, error) {
	var whirlpoolsApi WhirlpoolsApi
	client := http.Client{}
	request, err := http.NewRequest("GET", API_WHILPOOLS, nil)
	if err != nil {
		return whirlpoolsApi, err
	}
	request.Header.Set("Content-Type", "application/json")
	// Make request
	response, err := client.Do(request)
	if err != nil {
		return whirlpoolsApi, err
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return whirlpoolsApi, err
	}
	defer response.Body.Close()
	err = json.Unmarshal(bodyBytes, &whirlpoolsApi)
	return whirlpoolsApi, err
}

type Whirlpool struct {
	Address string `json:"address"`
	TokenA  struct {
		Mint        string `json:"mint"`
		Symbol      string `json:"symbol"`
		Name        string `json:"name"`
		Decimals    int    `json:"decimals"`
		LogoURI     string `json:"logoURI"`
		CoingeckoID string `json:"coingeckoId"`
		Whitelisted bool   `json:"whitelisted"`
		PoolToken   bool   `json:"poolToken"`
		Wrapper     string `json:"wrapper,omitempty"`
	} `json:"tokenA,omitempty"`
	TokenB struct {
		Mint        string `json:"mint"`
		Symbol      string `json:"symbol"`
		Name        string `json:"name"`
		Decimals    int    `json:"decimals"`
		LogoURI     string `json:"logoURI"`
		CoingeckoID string `json:"coingeckoId"`
		Whitelisted bool   `json:"whitelisted"`
		PoolToken   bool   `json:"poolToken"`
		Wrapper     string `json:"wrapper,omitempty"`
	} `json:"tokenB,omitempty"`
	Whitelisted      bool    `json:"whitelisted"`
	TickSpacing      int     `json:"tickSpacing"`
	Price            float64 `json:"price"`
	LpFeeRate        float64 `json:"lpFeeRate"`
	ProtocolFeeRate  float64 `json:"protocolFeeRate"`
	WhirlpoolsConfig string  `json:"whirlpoolsConfig"`
	ModifiedTimeMs   int64   `json:"modifiedTimeMs,omitempty"`
	Tvl              float64 `json:"tvl,omitempty"`
	Volume           struct {
		Day   float64 `json:"day"`
		Week  float64 `json:"week"`
		Month float64 `json:"month"`
	} `json:"volume,omitempty"`
	VolumeDenominatedA struct {
		Day   float64 `json:"day"`
		Week  float64 `json:"week"`
		Month float64 `json:"month"`
	} `json:"volumeDenominatedA,omitempty"`
	VolumeDenominatedB struct {
		Day   float64 `json:"day"`
		Week  float64 `json:"week"`
		Month float64 `json:"month"`
	} `json:"volumeDenominatedB,omitempty"`
	PriceRange struct {
		Day struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"day"`
		Week struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"week"`
		Month struct {
			Min float64 `json:"min"`
			Max float64 `json:"max"`
		} `json:"month"`
	} `json:"priceRange,omitempty"`
	FeeApr struct {
		Day   float64 `json:"day"`
		Week  float64 `json:"week"`
		Month float64 `json:"month"`
	} `json:"feeApr,omitempty"`
	Reward0Apr struct {
		Day   float64 `json:"day"`
		Week  float64 `json:"week"`
		Month float64 `json:"month"`
	} `json:"reward0Apr,omitempty"`
	Reward1Apr struct {
		Day   float64 `json:"day"`
		Week  float64 `json:"week"`
		Month float64 `json:"month"`
	} `json:"reward1Apr,omitempty"`
	Reward2Apr struct {
		Day   float64 `json:"day"`
		Week  float64 `json:"week"`
		Month float64 `json:"month"`
	} `json:"reward2Apr,omitempty"`
	TotalApr struct {
		Day   float64 `json:"day"`
		Week  float64 `json:"week"`
		Month float64 `json:"month"`
	} `json:"totalApr,omitempty"`
}

type WhirlpoolsApi struct {
	Whirlpools []Whirlpool `json:"whirlpools"`
	HasMore    bool        `json:"hasMore"`
}

func (w Whirlpool) HasSymbol(symbol string) bool {
	if w.TokenA.Symbol == symbol || w.TokenB.Symbol == symbol {
		return true
	}
	return false
}

func (w Whirlpool) HasMintKey(mint solana.PublicKey) bool {
	if w.TokenA.Mint == mint.String() || w.TokenB.Mint == mint.String() {
		return true
	}
	return false
}

func (w Whirlpool) HasMintKeyAsString(mint string) bool {
	if w.TokenA.Mint == mint || w.TokenB.Mint == mint {
		return true
	}
	return false
}

// find pools with symbola symbolb
func (wapi WhirlpoolsApi) GetPoolsViaSymbols(symbolA, symbolB string) []Whirlpool {
	found := make([]Whirlpool, 0)
	for _, whirlpool := range wapi.Whirlpools {
		if whirlpool.HasSymbol(symbolA) && whirlpool.HasSymbol(symbolB) {
			found = append(found, whirlpool)
		}
	}
	sort.SliceStable(found, func(i, j int) bool {
		return int(found[i].Tvl) > int(found[j].Tvl)
	})
	return found
}

// find pools with mints
func (wapi WhirlpoolsApi) GetPoolsViaMints(mintA, mintB solana.PublicKey) []Whirlpool {
	found := make([]Whirlpool, 0)
	for _, whirlpool := range wapi.Whirlpools {
		if whirlpool.HasMintKey(mintA) && whirlpool.HasMintKey(mintB) {
			found = append(found, whirlpool)
		}
	}
	sort.SliceStable(found, func(i, j int) bool {
		return int(found[i].Tvl) > int(found[j].Tvl)
	})
	return found
}

// find pools with symbola symbolb
func (wapi WhirlpoolsApi) GetPoolsViaMintsAsStrins(mintA, mintB string) []Whirlpool {
	found := make([]Whirlpool, 0)
	for _, whirlpool := range wapi.Whirlpools {
		if whirlpool.HasMintKeyAsString(mintA) && whirlpool.HasMintKeyAsString(mintB) {
			found = append(found, whirlpool)
		}
	}
	sort.SliceStable(found, func(i, j int) bool {
		return int(found[i].Tvl) > int(found[j].Tvl)
	})
	return found
}
