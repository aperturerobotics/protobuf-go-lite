package vtcompare

import "testing"

type testCase struct {
	val int
}

func (t *testCase) EqualVT(ot *testCase) bool {
	return t == ot
}

func TestCompareVT(t *testing.T) {
	t1, t2 := &testCase{val: 1}, &testCase{val: 2}
	cmp := CompareEqualVT[*testCase]()
	if cmp(t1, t2) {
		t.Fail()
	}
	if !cmp(t1, t1) {
		t.Fail()
	}
	if !cmp(nil, nil) {
		t.Fail()
	}
	if cmp(t1, nil) {
		t.Fail()
	}
}
