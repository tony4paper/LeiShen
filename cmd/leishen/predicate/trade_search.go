package predicate

import (
	"github.com/ethereum/go-ethereum/common"
)

func SearchTrades(records []Record) []Trade {
	var result []Trade
	for i, j := 0, 1; j < len(records); i, j = i+1, j+1 {
		a := records[i]
		b := records[j]

		k := j + 1
		if k < len(records) {
			c := records[k]
			trades := Get2Vs1Trade(a, b, c)
			if len(trades) > 0 {
				result = append(result, trades...)
				i, j = i+2, j+2
				continue
			}

			trades = Get1Vs2Trade(a, b, c)
			if len(trades) > 0 {
				result = append(result, trades...)
				i, j = i+2, j+2
				continue
			}
		}

		trades := Is1Vs1Trade(a, b)
		if len(trades) > 0 {
			result = append(result, trades...)
			i, j = i+1, j+1
			continue
		}
	}

	return result
}

var (
	ZeroAddr = common.Address{}
)

// A -> B a_1 t_1
// A -> B a_2 t_2
// 0 -> A a_3 t_3
// (A, t1) <-> (B, t3) a1/a3
// (A, t2) <-> (B, t3) a2/a3
// (A, t1) <-> (B, t2) a1/a2
func Get2Vs1Trade(a, b, c Record) []Trade {
	var result []Trade
	ok := a.From() == b.From() && a.From() == c.To() && a.To() == b.To() &&
		a.FromAddress() != ZeroAddr && a.ToAddress() != ZeroAddr &&
		c.FromAddress() == ZeroAddr &&
		a.Token() != b.Token() && a.Token() != c.Token() && b.Token() != c.Token()

	if ok {
		A := a.From()
		B := a.To()

		a = WrapWithTrade(a, A, B)
		b = WrapWithTrade(b, A, B)
		c = WrapWithTrade(c, B, A)

		result = append(result, Trade{a, c})
		result = append(result, Trade{b, c})
		result = append(result, Trade{a, b})
	}

	return result
}

// A -> B a_1 t_1
// B -> A a_2 t_2
// B -> A a_3 t_3
//
// A -> B a_1 t_1
// 0 -> B a_2 t_2
// 0 -> A a_3 t_3
//
// A -> 0 a_1 t_1
// B -> A a_2 t_2
// B -> A a_3 t_3
//
// (A, t1) <-> (B, t2) a1/a2
// (A, t1) <-> (B, t3) a1/a3
// (A, t2) <-> (B, t3) a2/a3
func Get1Vs2Trade(a, b, c Record) []Trade {
	var result []Trade

	if !(a.Token() != b.Token() && a.Token() != c.Token() && b.Token() != c.Token()) {
		return result
	}

	// A -> B a_1 t_1
	// B -> A a_2 t_2
	// B -> A a_3 t_3
	pattern1 := a.From() == b.To() && a.From() == c.To() &&
		a.To() == b.From() && a.To() == c.From() &&
		a.FromAddress() != ZeroAddr && a.ToAddress() != ZeroAddr
	if pattern1 {
		A := a.From()
		B := a.To()

		a = WrapWithTrade(a, A, B)
		b = WrapWithTrade(b, B, A)
		c = WrapWithTrade(c, B, A)

		result = append(result, Trade{a, b})
		result = append(result, Trade{a, c})
		result = append(result, Trade{b, c})
		return result
	}

	// A -> B a_1 t_1
	// 0 -> B a_2 t_2
	// 0 -> A a_3 t_3
	pattern2 := a.From() == c.To() && a.To() == b.To() &&
		b.FromAddress() == ZeroAddr && c.FromAddress() == ZeroAddr &&
		a.FromAddress() != ZeroAddr && a.ToAddress() != ZeroAddr
	if pattern2 {
		A := a.From()
		B := a.To()

		a = WrapWithTrade(a, A, B)
		b = WrapWithTrade(b, A, B)
		c = WrapWithTrade(c, B, A)

		result = append(result, Trade{a, b})
		result = append(result, Trade{a, c})
		result = append(result, Trade{b, c})
		return result
	}

	// A -> 0 a_1 t_1
	// B -> A a_2 t_2
	// B -> A a_3 t_3
	pattern3 := a.From() == b.To() && a.From() == c.To() &&
		b.From() == c.From() && a.ToAddress() == ZeroAddr &&
		a.FromAddress() != ZeroAddr && b.FromAddress() != ZeroAddr
	if pattern3 {
		A := a.From()
		B := b.From()

		a = WrapWithTrade(a, A, B)
		b = WrapWithTrade(b, B, A)
		c = WrapWithTrade(c, B, A)

		result = append(result, Trade{a, b})
		result = append(result, Trade{a, c})
		result = append(result, Trade{b, c})
	}

	return result
}

// A -> B a_1 t_1
// B -> A a_2 t_2
//
// 0 -> B a_1 t_1
// B -> A a_2 t_2
//
// A -> B a_1 t_1
// 0 -> A a_2 t_2
//
// A -> 0 a_1 t_1
// B -> A a_2 t_2
//
// A -> B a_1 t_1
// B -> 0 a_2 t_2
//
// (A ,t1) <-> (B, t2) a1/a2
func Is1Vs1Trade(a, b Record) []Trade {
	var result []Trade

	if a.Token() == b.Token() {
		return result
	}

	// A -> B a_1 t_1
	// B -> A a_2 t_2
	pattern1 := a.From() == b.To() && a.To() == b.From() &&
		a.FromAddress() != ZeroAddr && a.ToAddress() != ZeroAddr
	if pattern1 {
		A := a.From()
		B := b.From()

		a = WrapWithTrade(a, A, B)
		b = WrapWithTrade(b, B, A)

		result = append(result, Trade{a, b})
		return result
	}

	// 0 -> B a_1 t_1
	// B -> A a_2 t_2
	pattern2_1 := a.To() == b.From() &&
		a.FromAddress() == ZeroAddr &&
		b.FromAddress() != ZeroAddr && b.ToAddress() != ZeroAddr
	if pattern2_1 {
		A := b.To()
		B := b.From()

		a = WrapWithTrade(a, A, B)
		b = WrapWithTrade(b, B, A)

		result = append(result, Trade{a, b})
		return result
	}

	// A -> B a_1 t_1
	// 0 -> A a_2 t_2
	pattern2_2 := a.From() == b.To() &&
		b.FromAddress() == ZeroAddr &&
		a.FromAddress() != ZeroAddr && a.ToAddress() != ZeroAddr
	if pattern2_2 {
		A := a.From()
		B := a.To()

		a = WrapWithTrade(a, A, B)
		b = WrapWithTrade(b, B, A)

		result = append(result, Trade{a, b})
		return result
	}

	// A -> 0 a_1 t_1
	// B -> A a_2 t_2
	pattern3_1 := a.From() == b.To() &&
		a.ToAddress() == ZeroAddr &&
		b.FromAddress() != ZeroAddr && b.ToAddress() != ZeroAddr
	if pattern3_1 {
		A := a.From()
		B := b.From()

		a = WrapWithTrade(a, A, B)
		b = WrapWithTrade(b, B, A)

		result = append(result, Trade{a, b})
		return result
	}

	// A -> B a_1 t_1
	// B -> 0 a_2 t_2
	pattern3_2 := a.To() == b.From() &&
		b.ToAddress() == ZeroAddr &&
		a.FromAddress() != ZeroAddr && a.ToAddress() != ZeroAddr
	if pattern3_2 {
		A := a.From()
		B := b.From()

		a = WrapWithTrade(a, A, B)
		b = WrapWithTrade(b, B, A)

		result = append(result, Trade{a, b})
	}

	return result
}
