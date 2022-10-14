package predicate

import (
	"fmt"
	"math/big"
)

func NewTrade(a, b Record) Trade {
	return Trade{
		A: a,
		B: b,
	}
}

type Trade struct {
	A Record // A -> B, A.Token(), A.Amount()
	B Record // B -> A, B.Token(), B.Amount()
}

func (t *Trade) SwappedToken() string {
	if t.A.Token() < t.B.Token() {
		return fmt.Sprintf("(%s, %s)", t.A.Token(), t.B.Token())
	}
	return fmt.Sprintf("(%s, %s)", t.B.Token(), t.A.Token())
}

func (t *Trade) Paticipant() string {
	if t.A.From() < t.A.To() {
		return fmt.Sprintf("(%s, %s)", t.A.From(), t.A.To())
	}

	return fmt.Sprintf("(%s, %s)", t.A.To(), t.A.From())
}

func (t *Trade) SwappedTokenAndPaticipant() string {
	return fmt.Sprintf("(%s, %s)", t.SwappedToken(), t.Paticipant())
}

func (t *Trade) Sender(direct bool) string {
	var a, b string
	if t.A.Token() < t.B.Token() {
		a = t.A.From()
		b = t.A.To()
	} else {
		a = t.A.To()
		b = t.A.From()
	}

	if direct {
		return a
	}
	return b
}

func (t *Trade) Receiver(direct bool) string {
	var a, b string
	if t.A.Token() < t.B.Token() {
		a = t.A.From()
		b = t.A.To()
	} else {
		a = t.A.To()
		b = t.A.From()
	}

	if direct {
		return b
	}
	return a
}

func (t *Trade) ExchangeRate(direct bool) *big.Float {
	var a, b *big.Float
	if t.A.Token() < t.B.Token() {
		a = big.NewFloat(0).SetInt(t.A.Amount())
		b = big.NewFloat(0).SetInt(t.B.Amount())
	} else {
		a = big.NewFloat(0).SetInt(t.B.Amount())
		b = big.NewFloat(0).SetInt(t.A.Amount())
	}

	if direct {
		return a.Quo(a, b)
	}
	return b.Quo(b, a)
}

type PeakOrTrough int

const (
	Peak   PeakOrTrough = 1
	Nomal  PeakOrTrough = 0
	Trough PeakOrTrough = -1
)

func (v *PeakOrTrough) IsPeak() bool {
	return *v == Peak
}

func (v *PeakOrTrough) IsNomal() bool {
	return *v == Nomal
}

func (v *PeakOrTrough) IsTrough() bool {
	return *v == Trough
}

func NextPeak(arr []PeakOrTrough) int {
	for i, v := range arr {
		if v.IsPeak() {
			return i
		}
	}

	return -1
}

func NextTrough(arr []PeakOrTrough) int {
	for i, v := range arr {
		if v.IsTrough() {
			return i
		}
	}

	return -1
}

func IsPeakOrTrough(trades []Trade, direct bool, i int) PeakOrTrough {
	if IsPeak(trades, direct, i) {
		return Peak
	}

	if IsTrough(trades, direct, i) {
		return Trough
	}

	return Nomal
}

func IsPeak(trades []Trade, direct bool, i int) bool {
	prev, next := i-1, i+1
	if prev >= 0 && trades[prev].ExchangeRate(direct).Cmp(trades[i].ExchangeRate(direct)) > 0 {
		return false
	}

	if next < len(trades) && trades[next].ExchangeRate(direct).Cmp(trades[i].ExchangeRate(direct)) > 0 {
		return false
	}

	return true
}

func IsTrough(trades []Trade, direct bool, i int) bool {
	prev, next := i-1, i+1
	if prev >= 0 && trades[prev].ExchangeRate(direct).Cmp(trades[i].ExchangeRate(direct)) < 0 {
		return false
	}

	if next < len(trades) && trades[next].ExchangeRate(direct).Cmp(trades[i].ExchangeRate(direct)) < 0 {
		return false
	}

	return true
}
