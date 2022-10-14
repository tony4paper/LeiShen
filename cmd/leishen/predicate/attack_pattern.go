package predicate

import (
	"math/big"
)

// Bzx-2 0x762881b07feb63c436dee38edd4ff1f7a74c33091e534af56c9f7d49b5ecac15 18
// Balancer 0x013be97768b702fe8eccef1a40544d5ecb3c1961ad5f87fee4d16fdc08c78106 18
type RaisingPrice struct {
	Direct        bool
	RisingCount   int
	Trades        []Trade
	PeakOrTroughs []PeakOrTrough
}

func SearchRaisingPrice(trades []Trade) []RaisingPrice {
	var result []RaisingPrice
	token_map := make(map[string][]Trade)
	for _, trade := range trades {
		token_map[trade.SwappedTokenAndPaticipant()] = append(token_map[trade.SwappedTokenAndPaticipant()], trade)
	}

	for _, trades := range token_map {
		result = append(result, SearchRaisingPriceInSameNamedTrades(trades, false)...)
		result = append(result, SearchRaisingPriceInSameNamedTrades(trades, true)...)
	}

	return result
}

const raisingRepeatThreshold = 5 - 1

func SearchRaisingPriceInSameNamedTrades(trades []Trade, derect bool) []RaisingPrice {
	var result []RaisingPrice
	peak_or_troughs := []PeakOrTrough{}
	for i := range trades {
		peak_or_troughs = append(peak_or_troughs, IsPeakOrTrough(trades, derect, i))
	}

	for i := 0; i < len(peak_or_troughs); {
		first_trough := NextTrough(peak_or_troughs[i:])
		if first_trough == -1 {
			break
		}
		first_trough += i

		second_trough := NextTrough(peak_or_troughs[first_trough+1:])
		if second_trough == -1 {
			break
		}
		second_trough += first_trough + 1

		peak := NextPeak(peak_or_troughs[first_trough : second_trough+1])
		if peak != -1 && peak >= raisingRepeatThreshold {
			result = append(result, RaisingPrice{
				Direct:        derect,
				RisingCount:   peak,
				Trades:        trades[first_trough : second_trough+1],
				PeakOrTroughs: peak_or_troughs[first_trough : second_trough+1],
			})
		}

		i = second_trough
	}

	return result
}

// Eminence 0x3503253131644dd9f52802d071de74e456570374d586ddd640159cf6fb9b8ad8 3
// Eminence 0x045b60411af18114f1986957a41296ba2a97ccff75a9b38af818800ea9da0b2a 3
// Eminence 0x4f0f495dbcb58b452f268b9149a418524e43b13b55e780673c10b3b755340317 3
// Harvest finance 0x35f8d2f572fceaac9288e5d462117850ef2694786992a8c3f6d02612277b0877 3
// Harvest finance 0x0fc6d2ca064fc841bc9b1c1fad1fbb97bcea5c9a1b2b66ef837f1227e06519a6 3
type MultiRoundPriceFluctuation struct {
	Direct        bool
	Round         int
	Trades        []Trade
	PeakOrTroughs []PeakOrTrough
}

func SearchMultiRoundPriceFluctuation(trades []Trade) []MultiRoundPriceFluctuation {
	var result []MultiRoundPriceFluctuation
	token_map := make(map[string][]Trade)
	for _, trade := range trades {
		token_map[trade.SwappedTokenAndPaticipant()] = append(token_map[trade.SwappedTokenAndPaticipant()], trade)
	}

	for _, trades := range token_map {
		result = append(result, SearchMultiRoundPriceFluctuationInSameNamedTrades(trades, false)...)
		result = append(result, SearchMultiRoundPriceFluctuationInSameNamedTrades(trades, true)...)
	}

	return result
}

const roundThreshold = 3

func SearchMultiRoundPriceFluctuationInSameNamedTrades(trades []Trade, derect bool) []MultiRoundPriceFluctuation {
	var result []MultiRoundPriceFluctuation
	peak_or_troughs := []PeakOrTrough{}
	for i := range trades {
		peak_or_troughs = append(peak_or_troughs, IsPeakOrTrough(trades, derect, i))
	}

	normal_cnt := 0
	for _, v := range peak_or_troughs {
		if v.IsNomal() {
			normal_cnt++
		}

	}

	round := (len(peak_or_troughs) - normal_cnt) / 2
	if round >= roundThreshold {
		result = append(result, MultiRoundPriceFluctuation{
			Direct:        derect,
			Round:         round,
			Trades:        trades,
			PeakOrTroughs: peak_or_troughs,
		})
	}

	return result
}

// Warp Finance(0x8bb8dc5c7c830bac85fa48acad2505e9300a91c3ff239c9517d0cae33b595090) 360%
// Cheese Bank(0x600a869aa3a259158310a233b815ff67ca41eab8961a49918c2031297a02f1cc) 94
type CollateralizingLiquidity struct {
	TradeI          []Record
	RecordI         Record
	TradeJ          []Record
	RecordJ         Record
	FluctuationRate *big.Float
}

var (
	priceFluctuation = big.NewFloat(1)
)

func SearchCollateralizingLiquidity(records []Record) []CollateralizingLiquidity {
	var result []CollateralizingLiquidity

	records_with_borrower := []Record{}
	for _, record := range records {
		if record.From() == Borrower || record.To() == Borrower {
			records_with_borrower = append(records_with_borrower, record)
		}
	}

out:
	for i, j, k := 0, 1, 2; k < len(records_with_borrower); i, j, k = i+1, j+1, k+1 {
		a := records_with_borrower[i]
		b := records_with_borrower[j]
		c := records_with_borrower[k]

		trades := Get2Vs1Trade(a, b, c)
		if len(trades) != 3 {
			continue
		}
		trade_i := trades[2] // 推导 trade
		for ii := k + 1; ii < len(records_with_borrower); ii++ {
			record_ii := records_with_borrower[ii]

			if record_ii.From() == Borrower && record_ii.Token() == c.Token() {

				for jj, kk := ii+1, ii+2; kk < len(records_with_borrower); jj, kk = jj+1, kk+1 {
					d := records_with_borrower[jj]
					e := records_with_borrower[kk]
					trades = Is1Vs1Trade(d, e)
					if len(trades) == 1 {
						trade_j := trades[0]

						if trade_i.SwappedTokenAndPaticipant() == trade_j.SwappedTokenAndPaticipant() {

							for ll := kk + 1; ll < len(records_with_borrower); ll++ {
								f := records_with_borrower[ll]

								if f.From() == record_ii.To() && f.To() == record_ii.From() {
									var r *big.Float
									r1 := rateFluctuation(trade_i.ExchangeRate(true), trade_j.ExchangeRate(true))
									r2 := rateFluctuation(trade_i.ExchangeRate(true), trade_j.ExchangeRate(true))
									if r1.Cmp(r1) > 0 {
										r = r1
									} else {
										r = r2
									}

									if r.Cmp(priceFluctuation) > 0 {
										result = append(result, CollateralizingLiquidity{
											TradeI:          []Record{a, b, c},
											RecordI:         record_ii,
											TradeJ:          []Record{d, e},
											RecordJ:         f,
											FluctuationRate: r1,
										})
										i, j, k = ll-2, ll-1, ll
										continue out
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}

var (
	buySellPriceFluctuation2 = big.NewFloat(0.3)
)

// 低买高卖
// Yearn finance 0xb094d168dd90fcd0946016b19494a966d3d2c348f57b890410c51425d89166e8
// Yearn finance 0x6dc268706818d1e6503739950abc5ba2211fc6b451e54244da7b1e226b12e027
type BuyLowSellHigh struct {
	Direct bool
	Trades []Trade
}

func SearchBuyLowSellHigh(trades []Trade) []BuyLowSellHigh {
	var result []BuyLowSellHigh
	token_map := make(map[string][]Trade)
	for _, trade := range trades {
		token_map[trade.SwappedToken()] = append(token_map[trade.SwappedToken()], trade)
	}

	for _, trades := range token_map {
		result = append(result, SearchBuyLowSellHighInSameNamedTrades(trades, false)...)
		result = append(result, SearchBuyLowSellHighInSameNamedTrades(trades, true)...)
	}

	return result
}

func SearchBuyLowSellHighInSameNamedTrades(trades []Trade, derect bool) []BuyLowSellHigh {
	var result []BuyLowSellHigh

	for i := 0; i < len(trades); i++ {
		a := trades[i]
		if a.Sender(derect) == Borrower {
			for j := i + 1; j < len(trades); j++ {
				b := trades[j]
				if b.Sender(derect) == Borrower {
					r := rateFluctuation(b.ExchangeRate(derect), a.ExchangeRate(derect))
					if r.Cmp(buySellPriceFluctuation2) > 0 {
						result = append(result, BuyLowSellHigh{
							Direct: derect,
							Trades: trades[i : j+1],
						})
						i = j
						break
					}
				}
			}
		}
	}

	return result
}
