package arena

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	ref "github.com/aperturerobotics/protobuf-go-lite/testproto/arena/ref"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/timestamppb"
)

// testRef builds a ref to digest i.
func testRef(i int) *ref.Ref {
	return &ref.Ref{Digest: &ref.Digest{Kind: 1, Sum: bytes.Repeat([]byte{byte(i)}, 32)}}
}

// testPage builds a page of n entries using only fields the arena slabs.
func testPage(n int) *Page {
	page := &Page{Level: 1, LowerBound: []byte("a")}
	for i := range n {
		page.Entries = append(page.Entries, &Entry{
			Key:   fmt.Appendf(nil, "key/%08d", i),
			Hash:  bytes.Repeat([]byte{byte(i)}, 16),
			Size:  uint64(i),
			Child: testRef(i),
		})
	}
	return page
}

// fullPage builds a page using every field.
func fullPage() *Page {
	page := testPage(3)
	version := int64(7)
	page.Label = "label"
	page.Kind = Kind_KIND_LEAF
	page.Sealed = true
	page.Sizes = []uint64{1, 300, 70000}
	page.Tags = [][]byte{[]byte("x"), {}, []byte("zz")}
	page.Named = map[string]*Entry{"one": {Key: []byte("k1")}, "empty": {}}
	page.Extra = &Page_Pinned{Pinned: &Entry{Key: []byte("pin"), Value: testRef(9)}}
	page.Version = &version
	page.Updated = &timestamppb.Timestamp{Seconds: 1700000000, Nanos: 5}
	page.RootRef = &ref.Ref{Digest: &ref.Digest{}, Name: "root"}
	page.Children = []*Page{testPage(2), {}}
	page.Weight = 0.5
	page.Entries[1].Value = testRef(4)
	page.Entries[1].Extra = []*ref.Ref{testRef(5), {}, testRef(6)}
	page.Entries[2].Labels = []string{"l1", ""}
	page.Entries[2].Key = []byte{}
	return page
}

// checkArena requires UnmarshalVTArena to match UnmarshalVT on data, field by
// field, and to retain no part of data.
func checkArena(t *testing.T, data []byte) {
	// Decode with both decoders and require the same result.
	t.Helper()
	var want, got Page
	wantErr := want.UnmarshalVT(data)
	gotErr := got.UnmarshalVTArena(data)
	if (wantErr == nil) != (gotErr == nil) {
		t.Fatalf("UnmarshalVTArena error = %v, UnmarshalVT error = %v", gotErr, wantErr)
	}
	if wantErr != nil {
		return
	}
	if !reflect.DeepEqual(&got, &want) {
		t.Fatalf("UnmarshalVTArena = %v, want %v", &got, &want)
	}

	// Overwrite the input: the decoded page must not change.
	for i := range data {
		data[i] ^= 0xff
	}
	if !reflect.DeepEqual(&got, &want) {
		t.Fatal("decoded page aliases its input")
	}
}

// inSlab reports whether entries sit in one array, as the arena carves them.
func inSlab(entries []*Entry) bool {
	first := reflect.ValueOf(entries[0]).Pointer()
	size := reflect.TypeFor[Entry]().Size()
	for i, ent := range entries {
		if reflect.ValueOf(ent).Pointer() != first+uintptr(i)*size {
			return false
		}
	}
	return true
}

// TestArenaMatchesUnmarshalVT checks the arena decoder on encoded pages, and
// that pages with entries take the arena path.
func TestArenaMatchesUnmarshalVT(t *testing.T) {
	pages := map[string]*Page{
		"empty": {},
		"slab":  testPage(32),
		"full":  fullPage(),
	}
	for name, page := range pages {
		t.Run(name, func(t *testing.T) {
			// Encode the page.
			data, err := page.MarshalVT()
			if err != nil {
				t.Fatal(err)
			}

			// Require the arena path, then the same page as UnmarshalVT.
			var arena Page
			if err := arena.UnmarshalVTArena(data); err != nil {
				t.Fatal(err)
			}
			if len(page.Entries) != 0 && !inSlab(arena.Entries) {
				t.Fatal("page left the arena path")
			}
			checkArena(t, data)
		})
	}
}

// TestArenaUsesSlabs checks that a page's entries sit in one slab and that its
// allocations do not grow with its entries.
func TestArenaUsesSlabs(t *testing.T) {
	// Decode a page and require its entries to be adjacent.
	data, err := testPage(4).MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	var page Page
	if err := page.UnmarshalVTArena(data); err != nil {
		t.Fatal(err)
	}
	if !inSlab(page.Entries) {
		t.Fatal("entries are not carved from one slab")
	}

	// Require the same allocation count at 8 and 128 entries.
	allocs := func(n int) float64 {
		data, err := testPage(n).MarshalVT()
		if err != nil {
			t.Fatal(err)
		}
		return testing.AllocsPerRun(20, func() {
			var page Page
			if err := page.UnmarshalVTArena(data); err != nil {
				t.Fatal(err)
			}
		})
	}
	if small, large := allocs(8), allocs(128); small != large {
		t.Fatalf("allocations grow with entries: %v at 8, %v at 128", small, large)
	}
}

// TestArenaAppendIsolated checks that appending to a decoded bytes field does
// not write into the field after it.
func TestArenaAppendIsolated(t *testing.T) {
	// Decode a page.
	data, err := testPage(2).MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	var page Page
	if err := page.UnmarshalVTArena(data); err != nil {
		t.Fatal(err)
	}

	// Append to a key and require the hash after it unchanged.
	hashBefore := bytes.Clone(page.Entries[0].Hash)
	_ = append(page.Entries[0].Key, bytes.Repeat([]byte{0xee}, 64)...)
	if !bytes.Equal(page.Entries[0].Hash, hashBefore) {
		t.Fatal("append to a key overwrote the entry hash")
	}
}

// TestArenaNonEmptyMerges checks that decoding into a set page merges as
// UnmarshalVT does.
func TestArenaNonEmptyMerges(t *testing.T) {
	// Encode a page.
	data, err := fullPage().MarshalVT()
	if err != nil {
		t.Fatal(err)
	}

	// Decode it into a set page with both decoders.
	want, got := &Page{Level: 9, Tags: [][]byte{[]byte("t")}}, &Page{Level: 9, Tags: [][]byte{[]byte("t")}}
	if err := want.UnmarshalVT(data); err != nil {
		t.Fatal(err)
	}
	if err := got.UnmarshalVTArena(data); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatal("set page decode differs from UnmarshalVT")
	}
}

// FuzzArena checks that the arena decoder matches UnmarshalVT on any input.
func FuzzArena(f *testing.F) {
	// Seed with encoded pages.
	for _, page := range []*Page{{}, testPage(3), fullPage()} {
		data, err := page.MarshalVT()
		if err != nil {
			f.Fatal(err)
		}
		f.Add(data)
	}

	// Repeated singular messages, unknown fields, wrong wire types.
	f.Add([]byte{0x7a, 0x00, 0x7a, 0x02, 0x12, 0x00})
	f.Add([]byte{0x1a, 0x02, 0x22, 0x00, 0x1a, 0x04, 0x22, 0x02, 0x0a, 0x00})
	f.Add([]byte{0x1a, 0x02, 0x48, 0x01})
	f.Add([]byte{0x18, 0x01})
	f.Add([]byte{0x1a, 0x03, 0x22, 0x01})
	f.Add([]byte{0x82, 0x01, 0x02, 0x1a, 0x00, 0x82, 0x01, 0x00})

	// Compare the decoders on each input.
	f.Fuzz(func(t *testing.T, data []byte) {
		checkArena(t, data)
	})
}

// BenchmarkArena compares the arena decoder with UnmarshalVT.
func BenchmarkArena(b *testing.B) {
	for _, n := range []int{8, 32, 128} {
		data, err := testPage(n).MarshalVT()
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("n=%d/generated", n), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(data)))
			for range b.N {
				var page Page
				if err := page.UnmarshalVT(data); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(fmt.Sprintf("n=%d/arena", n), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(data)))
			for range b.N {
				var page Page
				if err := page.UnmarshalVTArena(data); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
