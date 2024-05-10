package protobuf_go_lite

import "testing"

type testCase struct {
	val int
}

func (t *testCase) EqualVT(ot *testCase) bool {
	if t == ot {
		return true
	}
	if (ot == nil) != (t == nil) {
		return false
	}
	if ot == nil {
		return true
	}
	return t.val == ot.val
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

func TestIsEqualVTSlice(t *testing.T) {
	testCases := []struct {
		s1, s2 []*testCase
		expect bool
	}{
		{
			s1:     []*testCase{{val: 1}, {val: 2}},
			s2:     []*testCase{{val: 1}, {val: 2}},
			expect: true,
		},
		{
			s1:     []*testCase{{val: 1}, {val: 2}},
			s2:     []*testCase{{val: 1}, {val: 3}},
			expect: false,
		},
		{
			s1:     []*testCase{{val: 1}, {val: 2}},
			s2:     []*testCase{{val: 1}, {val: 2}, {val: 3}},
			expect: false,
		},
		{
			s1:     []*testCase{{val: 1}, nil},
			s2:     []*testCase{{val: 1}, nil},
			expect: true,
		},
		{
			s1:     []*testCase{{val: 1}, nil},
			s2:     []*testCase{{val: 1}, {val: 2}},
			expect: false,
		},
		{
			s1:     []*testCase{},
			s2:     []*testCase{},
			expect: true,
		},
		{
			s1:     nil,
			s2:     nil,
			expect: true,
		},
	}

	for _, tc := range testCases {
		actual := IsEqualVTSlice(tc.s1, tc.s2)
		if actual != tc.expect {
			t.Errorf("IsEqualVTSlice(%v, %v) = %v; want %v", tc.s1, tc.s2, actual, tc.expect)
		}
	}
}
