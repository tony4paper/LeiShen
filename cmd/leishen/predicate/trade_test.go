package predicate

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func newTestOrderRecord(a, b, t string) Record {
	return &testOrderRecord{a, b, t}
}

type testOrderRecord struct {
	a, b, t string
}

func (r *testOrderRecord) Token() string {
	return r.t
}

func (r *testOrderRecord) From() string {
	return r.a
}

func (r *testOrderRecord) To() string {
	return r.b
}

func (r *testOrderRecord) Amount() *big.Int {
	panic("not implemented") // TODO: Implement
}

func (r *testOrderRecord) TokenAddress() common.Address {
	panic("not implemented") // TODO: Implement
}

func (r *testOrderRecord) FromAddress() common.Address {
	panic("not implemented") // TODO: Implement
}

func (r *testOrderRecord) ToAddress() common.Address {
	panic("not implemented") // TODO: Implement
}

func TestTradeOrder(t *testing.T) {
	r1 := newTestOrderRecord("a", "b", "t1")
	r2 := newTestOrderRecord("b", "a", "t2")
	trade1 := Trade{r1, r2}
	trade2 := Trade{r2, r1}

	if trade1.SwappedToken() != trade2.SwappedToken() {
		t.Errorf("SwappedToken does not satisfy semantics")
	}

	if trade1.Paticipant() != trade2.Paticipant() {
		t.Errorf("Paticipant does not satisfy semantics")
	}

	if trade1.SwappedTokenAndPaticipant() != trade2.SwappedTokenAndPaticipant() {
		t.Errorf("SwappedTokenAndPaticipant does not satisfy semantics")
	}
}
