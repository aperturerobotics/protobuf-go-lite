package protobuf_go_lite

import (
	"io"
	"math"
	"testing"
)

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

func TestDecodeVarint(t *testing.T) {
	// Test basic varint decoding
	buf := AppendVarint(nil, 300)
	v, idx, err := DecodeVarint(buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if v != 300 {
		t.Errorf("DecodeVarint got %d, want 300", v)
	}
	if idx != len(buf) {
		t.Errorf("DecodeVarint idx got %d, want %d", idx, len(buf))
	}

	// Test empty buffer
	_, _, err = DecodeVarint(nil, 0)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("DecodeVarint empty got %v, want ErrUnexpectedEOF", err)
	}
}

func TestDecodeVarintTyped(t *testing.T) {
	buf := AppendVarint(nil, 0xFFFFFFFF)

	v32, _, err := DecodeVarintUint32(buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if v32 != 0xFFFFFFFF {
		t.Errorf("DecodeVarintUint32 got %d, want 0xFFFFFFFF", v32)
	}

	buf = AppendVarint(nil, uint64(math.MaxInt64))
	v64, _, err := DecodeVarintInt64(buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if v64 != math.MaxInt64 {
		t.Errorf("DecodeVarintInt64 got %d, want MaxInt64", v64)
	}
}

func TestDecodeVarintBool(t *testing.T) {
	bufTrue := AppendVarint(nil, 1)
	bufFalse := AppendVarint(nil, 0)

	b, _, err := DecodeVarintBool(bufTrue, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("DecodeVarintBool got false, want true")
	}

	b, _, err = DecodeVarintBool(bufFalse, 0)
	if err != nil {
		t.Fatal(err)
	}
	if b {
		t.Error("DecodeVarintBool got true, want false")
	}
}

func TestDecodeSint(t *testing.T) {
	// sint32: zigzag encode -1 (should be 1 in wire format)
	buf := AppendVarint(nil, 1) // zigzag(-1) = 1
	v, _, err := DecodeSint32(buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if v != -1 {
		t.Errorf("DecodeSint32 got %d, want -1", v)
	}

	// sint64: zigzag encode -1
	v64, _, err := DecodeSint64(buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if v64 != -1 {
		t.Errorf("DecodeSint64 got %d, want -1", v64)
	}
}

func TestDecodeFixed(t *testing.T) {
	// Fixed32
	buf32 := []byte{0x78, 0x56, 0x34, 0x12}
	v32, idx, err := DecodeFixed32(buf32, 0)
	if err != nil {
		t.Fatal(err)
	}
	if v32 != 0x12345678 {
		t.Errorf("DecodeFixed32 got %#x, want 0x12345678", v32)
	}
	if idx != 4 {
		t.Errorf("DecodeFixed32 idx got %d, want 4", idx)
	}

	// Fixed64
	buf64 := []byte{0xEF, 0xCD, 0xAB, 0x90, 0x78, 0x56, 0x34, 0x12}
	v64, idx, err := DecodeFixed64(buf64, 0)
	if err != nil {
		t.Fatal(err)
	}
	if v64 != 0x1234567890ABCDEF {
		t.Errorf("DecodeFixed64 got %#x, want 0x1234567890ABCDEF", v64)
	}
	if idx != 8 {
		t.Errorf("DecodeFixed64 idx got %d, want 8", idx)
	}

	// Test truncated buffer
	_, _, err = DecodeFixed32(buf32[:2], 0)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("DecodeFixed32 truncated got %v, want ErrUnexpectedEOF", err)
	}
}

func TestDecodeFloat(t *testing.T) {
	// Float32
	buf := make([]byte, 4)
	val := float32(3.14)
	bits := math.Float32bits(val)
	buf[0] = byte(bits)
	buf[1] = byte(bits >> 8)
	buf[2] = byte(bits >> 16)
	buf[3] = byte(bits >> 24)

	f, idx, err := DecodeFloat32(buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if f != val {
		t.Errorf("DecodeFloat32 got %v, want %v", f, val)
	}
	if idx != 4 {
		t.Errorf("DecodeFloat32 idx got %d, want 4", idx)
	}

	// Float64
	buf64 := make([]byte, 8)
	val64 := 3.14159265359
	bits64 := math.Float64bits(val64)
	for i := 0; i < 8; i++ {
		buf64[i] = byte(bits64 >> (i * 8))
	}

	f64, idx, err := DecodeFloat64(buf64, 0)
	if err != nil {
		t.Fatal(err)
	}
	if f64 != val64 {
		t.Errorf("DecodeFloat64 got %v, want %v", f64, val64)
	}
	if idx != 8 {
		t.Errorf("DecodeFloat64 idx got %d, want 8", idx)
	}
}

func TestDecodeString(t *testing.T) {
	// Create length-prefixed string
	s := "hello"
	buf := AppendVarint(nil, uint64(len(s)))
	buf = append(buf, s...)

	result, idx, err := DecodeString(buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result != s {
		t.Errorf("DecodeString got %q, want %q", result, s)
	}
	if idx != len(buf) {
		t.Errorf("DecodeString idx got %d, want %d", idx, len(buf))
	}

	// Test empty string
	emptyBuf := AppendVarint(nil, 0)
	result, _, err = DecodeString(emptyBuf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Errorf("DecodeString empty got %q, want empty", result)
	}
}

func TestDecodeBytes(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	buf := AppendVarint(nil, uint64(len(data)))
	buf = append(buf, data...)

	// With copy
	result, idx, err := DecodeBytes(buf, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != string(data) {
		t.Errorf("DecodeBytes got %v, want %v", result, data)
	}
	if idx != len(buf) {
		t.Errorf("DecodeBytes idx got %d, want %d", idx, len(buf))
	}

	// Without copy
	result2, _, err := DecodeBytes(buf, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if string(result2) != string(data) {
		t.Errorf("DecodeBytes nocopy got %v, want %v", result2, data)
	}
}

func TestDecodeStringUnsafe(t *testing.T) {
	s := "hello world"
	buf := AppendVarint(nil, uint64(len(s)))
	buf = append(buf, s...)

	result, idx, err := DecodeStringUnsafe(buf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result != s {
		t.Errorf("DecodeStringUnsafe got %q, want %q", result, s)
	}
	if idx != len(buf) {
		t.Errorf("DecodeStringUnsafe idx got %d, want %d", idx, len(buf))
	}

	// Test empty string
	emptyBuf := AppendVarint(nil, 0)
	result, _, err = DecodeStringUnsafe(emptyBuf, 0)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Errorf("DecodeStringUnsafe empty got %q, want empty", result)
	}
}
