package protobuf_go_lite

import (
	"bytes"
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

func (t *testCase) CloneVT() *testCase {
	if t == nil {
		return nil
	}
	out := *t
	return &out
}

func (t *testCase) CloneMessageVT() CloneMessage {
	return t.CloneVT()
}

func (t *testCase) SizeVT() int {
	return 0
}

func (t *testCase) MarshalToSizedBufferVT([]byte) (int, error) {
	return 0, nil
}

func (t *testCase) MarshalVT() ([]byte, error) {
	return nil, nil
}

func (t *testCase) UnmarshalVT([]byte) error {
	return nil
}

func (t *testCase) Reset() {
	*t = testCase{}
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

func TestCloneHelpers(t *testing.T) {
	v := 7
	vp := ClonePtr(&v)
	*vp = 8
	if v != 7 {
		t.Fatalf("ClonePtr aliases source")
	}
	if ClonePtr[int](nil) != nil {
		t.Fatalf("ClonePtr(nil) != nil")
	}

	bytesValue := []byte{1, 2}
	bytesClone := CloneBytes(bytesValue)
	bytesClone[0] = 9
	if bytesValue[0] != 1 {
		t.Fatalf("CloneBytes aliases source")
	}

	ints := []int{1, 2}
	intsClone := CloneSlice(ints)
	intsClone[0] = 9
	if ints[0] != 1 {
		t.Fatalf("CloneSlice aliases source")
	}

	intMap := map[string]int{"a": 1}
	intMapClone := CloneMap(intMap)
	intMapClone["a"] = 9
	if intMap["a"] != 1 {
		t.Fatalf("CloneMap aliases source")
	}

	if CloneBytesSlice([][]byte(nil)) != nil {
		t.Fatalf("CloneBytesSlice(nil) != nil")
	}
	bytesSlice := [][]byte{{1}, nil}
	bytesSliceClone := CloneBytesSlice(bytesSlice)
	bytesSliceClone[0][0] = 9
	if bytesSlice[0][0] != 1 || bytesSliceClone[1] != nil {
		t.Fatalf("CloneBytesSlice did not preserve deep copy and nil element")
	}

	bytesMap := map[string][]byte{"a": {1}, "nil": nil}
	bytesMapClone := CloneBytesMap(bytesMap)
	bytesMapClone["a"][0] = 9
	if bytesMap["a"][0] != 1 || bytesMapClone["nil"] != nil {
		t.Fatalf("CloneBytesMap did not preserve deep copy and nil value")
	}

	if CloneVTSlice([]*testCase(nil)) != nil {
		t.Fatalf("CloneVTSlice(nil) != nil")
	}
	msgSlice := []*testCase{{val: 1}, nil}
	msgSliceClone := CloneVTSlice(msgSlice)
	msgSliceClone[0].val = 9
	if msgSlice[0].val != 1 || msgSliceClone[1] != nil {
		t.Fatalf("CloneVTSlice did not preserve deep copy and nil element")
	}

	msgMap := map[string]*testCase{"a": {val: 1}, "nil": nil}
	msgMapClone := CloneVTMap(msgMap)
	msgMapClone["a"].val = 9
	if msgMap["a"].val != 1 {
		t.Fatalf("CloneVTMap aliases message values")
	}
	if _, ok := msgMapClone["nil"]; !ok || msgMapClone["nil"] != nil {
		t.Fatalf("CloneVTMap did not preserve nil-valued key")
	}
}

func TestEqualHelpers(t *testing.T) {
	if !EqualPtr(ptr(1), ptr(1)) || EqualPtr(ptr(1), ptr(2)) || EqualPtr(ptr(1), nil) {
		t.Fatalf("EqualPtr mismatch")
	}
	if !EqualBytes(nil, []byte{}) {
		t.Fatalf("EqualBytes should treat nil and empty as equal")
	}
	if EqualBytesPresent(nil, []byte{}) {
		t.Fatalf("EqualBytesPresent should distinguish nil and empty")
	}
	if !EqualSlice([]int{1, 2}, []int{1, 2}) || EqualSlice([]int{1}, []int{2}) {
		t.Fatalf("EqualSlice mismatch")
	}
	if !EqualMap(map[string]int{"a": 1}, map[string]int{"a": 1}) || EqualMap(map[string]int{"a": 1}, map[string]int{"a": 2}) {
		t.Fatalf("EqualMap mismatch")
	}
	if !EqualBytesSlice([][]byte{{1}, nil}, [][]byte{{1}, []byte{}}) {
		t.Fatalf("EqualBytesSlice should use implicit bytes equality")
	}
	if !EqualBytesMap(map[string][]byte{"a": nil}, map[string][]byte{"a": {}}) {
		t.Fatalf("EqualBytesMap should use implicit bytes equality")
	}

	empty := func() *testCase { return &testCase{} }
	if !EqualVTImplicit((*testCase)(nil), &testCase{}, empty) {
		t.Fatalf("EqualVTImplicit should treat nil as empty")
	}
	if !EqualVTSliceImplicit([]*testCase{{val: 1}, nil}, []*testCase{{val: 1}, {}}, empty) {
		t.Fatalf("EqualVTSliceImplicit mismatch")
	}
	if !EqualVTMapImplicit(map[string]*testCase{"a": nil}, map[string]*testCase{"a": {}}, empty) {
		t.Fatalf("EqualVTMapImplicit mismatch")
	}
}

func TestSizeHelpers(t *testing.T) {
	if got, want := SizeVarintValue(1, uint32(300)), 1+SizeOfVarint(300); got != want {
		t.Fatalf("SizeVarintValue = %d, want %d", got, want)
	}
	if SizeVarintNonZero(1, int32(0)) != 0 || SizeVarintPtr[int32](1, nil) != 0 {
		t.Fatalf("varint zero or nil size mismatch")
	}
	if got, want := SizeVarintPacked(1, []uint32{1, 300}), SizeBytesValue(1, SizeOfVarint(1)+SizeOfVarint(300)); got != want {
		t.Fatalf("SizeVarintPacked = %d, want %d", got, want)
	}
	sint := int32(-1)
	if got, want := SizeZigzagValue(1, sint), 1+SizeOfZigzag(uint64(sint)); got != want {
		t.Fatalf("SizeZigzagValue = %d, want %d", got, want)
	}
	if SizeFixed32NonZero(1, float32(0)) != 0 || SizeFixed64NonZero(1, float64(0)) != 0 {
		t.Fatalf("fixed zero size mismatch")
	}
	if got, want := SizeFixed32Packed(1, []uint32{1, 2}), SizeBytesValue(1, 8); got != want {
		t.Fatalf("SizeFixed32Packed = %d, want %d", got, want)
	}
	if SizeBoolNonZero(1, false) != 0 || SizeBoolValue(1) != 2 {
		t.Fatalf("bool size mismatch")
	}
	if got, want := SizeStringValue(1, "abc"), SizeBytesValue(1, 3); got != want {
		t.Fatalf("SizeStringValue = %d, want %d", got, want)
	}
	if SizeStringNonEmpty(1, "") != 0 || SizeStringPtr(1, (*string)(nil)) != 0 {
		t.Fatalf("string empty or nil size mismatch")
	}
	if got, want := SizeBytesSlice(1, [][]byte{{1, 2}, nil}), SizeBytesValue(1, 2)+SizeBytesValue(1, 0); got != want {
		t.Fatalf("SizeBytesSlice = %d, want %d", got, want)
	}
	if SizeBytesNonEmpty(1, nil) != 0 || SizeBytesPresent(1, nil) != 0 {
		t.Fatalf("bytes empty or nil size mismatch")
	}
	if SizeMessage(1, 3) != SizeBytesValue(1, 3) || SizeGroup(1, 3) != 5 {
		t.Fatalf("message or group size mismatch")
	}
}

func TestEncodeHelpers(t *testing.T) {
	check := func(name string, size int, encode func([]byte, int) int, want []byte) {
		t.Helper()
		buf := make([]byte, size)
		if got := encode(buf, len(buf)); got != 0 {
			t.Fatalf("%s offset = %d, want 0", name, got)
		}
		if !bytes.Equal(buf, want) {
			t.Fatalf("%s bytes = %v, want %v", name, buf, want)
		}
	}

	check("varint", 2, func(buf []byte, offset int) int {
		return EncodeVarint(buf, offset, 300)
	}, []byte{0xac, 0x02})
	check("raw bytes", 2, func(buf []byte, offset int) int {
		return EncodeRawBytes(buf, offset, []byte{0x11, 0x22})
	}, []byte{0x11, 0x22})
	check("fixed32", 4, func(buf []byte, offset int) int {
		return EncodeFixed32(buf, offset, 0x12345678)
	}, []byte{0x78, 0x56, 0x34, 0x12})
	check("fixed64", 8, func(buf []byte, offset int) int {
		return EncodeFixed64(buf, offset, 0x1234567890abcdef)
	}, []byte{0xef, 0xcd, 0xab, 0x90, 0x78, 0x56, 0x34, 0x12})
	check("bool true", 1, func(buf []byte, offset int) int {
		return EncodeBool(buf, offset, true)
	}, []byte{1})
	check("bool false", 1, func(buf []byte, offset int) int {
		return EncodeBool(buf, offset, false)
	}, []byte{0})
	check("string", 4, func(buf []byte, offset int) int {
		return EncodeString(buf, offset, "abc")
	}, []byte{0x03, 'a', 'b', 'c'})
	check("bytes", 3, func(buf []byte, offset int) int {
		return EncodeBytes(buf, offset, []byte{0x11, 0x22})
	}, []byte{0x02, 0x11, 0x22})
	check("zigzag32", 1, func(buf []byte, offset int) int {
		return EncodeZigzag32(buf, offset, int32(-1))
	}, []byte{0x01})
	check("zigzag64", 1, func(buf []byte, offset int) int {
		return EncodeZigzag64(buf, offset, int64(-1))
	}, []byte{0x01})
	check("varint packed", 4, func(buf []byte, offset int) int {
		return EncodeVarintPacked(buf, offset, []int32{1, 300})
	}, []byte{0x03, 0x01, 0xac, 0x02})
	check("zigzag32 packed", 4, func(buf []byte, offset int) int {
		return EncodeZigzag32Packed(buf, offset, []int32{-1, 150})
	}, []byte{0x03, 0x01, 0xac, 0x02})
	check("zigzag64 packed", 4, func(buf []byte, offset int) int {
		return EncodeZigzag64Packed(buf, offset, []int64{-1, 150})
	}, []byte{0x03, 0x01, 0xac, 0x02})
}

func ptr[T any](v T) *T {
	return &v
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
	for i := range 8 {
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
