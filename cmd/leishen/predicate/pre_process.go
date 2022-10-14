package predicate

import (
	"leishen/pkg/edb"
	"leishen/pkg/etypes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

func RemoveZeroRecord(records []Record) []Record {
	var result []Record
	for _, record := range records {
		if record.Amount().Cmp(big.NewInt(0)) != 0 {
			result = append(result, record)
		}
	}

	return result
}

func UnifiedAddressName(contractDB, platformNameDB ethdb.Database, fltx *etypes.FlashLoan, records []Record) []Record {
	var result []Record
	for _, record := range records {
		from := getAddressName(contractDB, platformNameDB, fltx, record.FromAddress())
		to := getAddressName(contractDB, platformNameDB, fltx, record.ToAddress())
		result = append(result, WrapWithPlatformName(record, from, to))
	}

	return result
}

const (
	Borrower = "Borrower"
	WETH     = "ETH"
	WBTC     = "BTC"
)

var (
	WETHAddr = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	WBTCAddr = common.HexToAddress("0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599")
)

func getAddressName(contractDB, platformNameDB ethdb.Database, fltx *etypes.FlashLoan, addr common.Address) string {
	var addrName string
	for _, borrower := range fltx.Types2Loaner {
		if addr == borrower {
			return Borrower
		}

		if addr == WETHAddr {
			return WETH
		}

		if addr == WBTCAddr {
			return WBTC
		}
	}

	addr_platform_name := edb.ReadPlatformName(platformNameDB, addr)
	if addr_platform_name == nil || len(addr_platform_name.NameMap) != 1 {
		eoa := getEOAPath(contractDB, addr)
		addrName = eoa.Hex()
	} else {
		addrName = addr_platform_name.Names()[0]
	}

	return addrName
}

func getEOAPath(contractDB ethdb.Database, addr common.Address) common.Address {
	for {
		contract := edb.ReadContract(contractDB, &addr)
		if contract == nil {
			return addr
		}

		addr = contract.Creator
	}
}

func RemoveInternalRecord(records []Record) []Record {
	var result []Record
	for _, record := range records {
		if record.From() != record.To() {
			result = append(result, record)
		}
	}

	return result
}

var (
	ChangeThreshold = big.NewFloat(0.01)
)

func MergeChange(records []Record) []Record {
	var result []Record
	for i, j, k := 0, 1, 2; k < len(records); i, j, k = i+1, j+1, k+1 {
		a := records[i]
		b := records[j]
		c := records[k]

		if !(a.Token() == b.Token() && a.Token() == c.Token()) {
			result = append(result, a)
			if k+1 >= len(records) {
				result = append(result, b)
				result = append(result, c)
			}
			continue
		}

		if !(a.To() == b.From() && a.To() == c.From() && b.To() == WETH) {
			result = append(result, a)
			if k+1 >= len(records) {
				result = append(result, b)
				result = append(result, c)
			}
			continue
		}

		d := big.NewInt(0).Sub(a.Amount(), c.Amount())

		e := big.NewFloat(0).SetInt(d)
		base := big.NewFloat(0).SetInt(b.Amount())
		f := big.NewFloat(0).Sub(e, base)
		f = f.Abs(f)

		if f.Quo(f, base).Cmp(ChangeThreshold) < 0 {
			i, j, k = i+2, j+2, k+2
			result = append(result, WrapWithChange(a, b, c, b.Amount()))
			if k < len(records) && k+1 >= len(records) {
				result = append(result, records[j])
				result = append(result, records[k])
			} else if j < len(records) && j+1 >= len(records) {
				result = append(result, records[j])
			}
		} else {
			result = append(result, a)
			if k+1 >= len(records) {
				result = append(result, b)
				result = append(result, c)
			}
		}
	}

	return result
}

func RemoveWETHAndWBTC(records []Record) []Record {
	var result []Record
	for _, record := range records {
		if record.ToAddress() != WETHAddr && record.ToAddress() != WBTCAddr {
			result = append(result, record)
		}
	}

	return result
}

var (
	BlackHoleAddr = common.Address{}
)

func MergeBlackHoleTransferred(records []Record) []Record {
	var result []Record
	for i, j := 0, 1; j < len(records); i, j = i+1, j+1 {
		a := records[i]
		b := records[j]

		if a.Token() != b.Token() {
			result = append(result, a)
			if j+1 >= len(records) {
				result = append(result, b)
			}
			continue
		}

		if !(a.From() == b.From() && b.ToAddress() == BlackHoleAddr) {
			result = append(result, a)
			if j+1 >= len(records) {
				result = append(result, b)
			}
			continue
		}

		value := big.NewInt(0).Add(a.Amount(), b.Amount())
		result = append(result, WrapWithBlackHole(a, b, value))
		i, j = i+1, j+1
		if j < len(records) && j+1 >= len(records) {
			result = append(result, records[j])
		}
	}

	return result
}

var (
	TransitThreshold = big.NewFloat(0.01)
)

func MergeTransitTransferred(records []Record) []Record {
	var result []Record
	for i, j := 0, 1; j < len(records); i, j = i+1, j+1 {
		a := records[i]
		b := records[j]

		if a.Token() != b.Token() {
			result = append(result, a)
			if j+1 >= len(records) {
				result = append(result, b)
			}
			continue
		}

		if a.To() != b.From() {
			result = append(result, a)
			if j+1 >= len(records) {
				result = append(result, b)
			}
			continue
		}

		d := big.NewInt(0).Sub(a.Amount(), b.Amount())

		e := big.NewFloat(0).SetInt(d)
		base := big.NewFloat(0).SetInt(a.Amount())
		f := big.NewFloat(0).Sub(e, base)
		f = f.Abs(f)

		if f.Quo(f, base).Cmp(TransitThreshold) < 0 {
			i, j = i+1, j+1
			result = append(result, WrapWithTransit(a, b, a.Amount()))
			if j < len(records) && j+1 >= len(records) {
				result = append(result, records[j])
			}
		} else {
			result = append(result, a)
			if j+1 >= len(records) {
				result = append(result, b)
			}
		}
	}

	return result
}
