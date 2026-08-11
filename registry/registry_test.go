package registry

import (
	"errors"
	"sync"
	"testing"

	protobuf_go_lite "github.com/aperturerobotics/protobuf-go-lite"
	"github.com/aperturerobotics/protobuf-go-lite/types/known/anypb"
)

type stubMsg struct{}

func (m *stubMsg) Reset()      {}
func (m *stubMsg) SizeVT() int { return 0 }
func (m *stubMsg) MarshalVT() ([]byte, error) {
	return nil, nil
}
func (m *stubMsg) MarshalToSizedBufferVT([]byte) (int, error) { return 0, nil }
func (m *stubMsg) UnmarshalVT([]byte) error                   { return nil }

func TestRegistryRegisterAndLookup(t *testing.T) {
	r := NewRegistry()
	r.Register(Entry{
		FullName: "demo.ClickEvents",
		New:      func() protobuf_go_lite.Message { return &stubMsg{} },
		Options: []Option{
			{Name: "annotations.msg_opts.tenant", Values: []string{"BLINKIT"}},
			{Name: "annotations.msg_opts.table", Values: []string{"click_events"}},
			{Name: "annotations.msg_opts.tenants", Values: []string{"BLINKIT", "ZOMATO"}},
			{Name: "annotations.msg_opts.empty"},
		},
	})

	msg, ok := r.NewByName("demo.ClickEvents")
	if !ok || msg == nil {
		t.Fatal("NewByName failed")
	}
	msg, ok = r.NewByTypeURL("type.googleapis.com/demo.ClickEvents")
	if !ok || msg == nil {
		t.Fatal("NewByTypeURL failed")
	}
	msg, ok = r.NewByTypeURL("demo.ClickEvents")
	if !ok || msg == nil {
		t.Fatal("NewByTypeURL by name failed")
	}

	entries := r.All()
	if len(entries) != 1 {
		t.Fatalf("All len=%d", len(entries))
	}
	tenant, ok := entries[0].Option("annotations.msg_opts.tenant")
	if !ok || tenant != "BLINKIT" {
		t.Fatalf("tenant=%q ok=%v", tenant, ok)
	}
	if _, ok := entries[0].Option("missing"); ok {
		t.Fatal("missing option should be absent")
	}
	if _, ok := entries[0].Option("annotations.msg_opts.empty"); ok {
		t.Fatal("option without values should be absent")
	}

	var saw string
	r.Range(func(e Entry) bool {
		saw = e.FullName
		return true
	})
	if saw != "demo.ClickEvents" {
		t.Fatalf("Range saw %q", saw)
	}

	fn, err := r.FindMessageByURL("type.googleapis.com/demo.ClickEvents")
	if err != nil || fn == nil || fn() == nil {
		t.Fatalf("FindMessageByURL err=%v", err)
	}
	_, err = r.FindMessageByURL("type.googleapis.com/demo.Missing")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	_, err = anypb.UnmarshalNew(
		&anypb.Any{TypeUrl: "type.googleapis.com/demo.Missing"},
		"",
		r,
	)
	if !errors.Is(err, anypb.ErrNotFound) {
		t.Fatalf("UnmarshalNew error = %v, want anypb.ErrNotFound", err)
	}
}

func TestRegistryDuplicatePanics(t *testing.T) {
	r := NewRegistry()
	e := Entry{
		FullName: "demo.Dup",
		New:      func() protobuf_go_lite.Message { return &stubMsg{} },
	}
	r.Register(e)
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	r.Register(e)
}

func TestRegistryConcurrentReads(t *testing.T) {
	r := NewRegistry()
	r.Register(Entry{
		FullName: "demo.A",
		New:      func() protobuf_go_lite.Message { return &stubMsg{} },
	})
	r.Register(Entry{
		FullName: "demo.B",
		New:      func() protobuf_go_lite.Message { return &stubMsg{} },
	})

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = r.All()
				r.Range(func(Entry) bool { return true })
				_, _ = r.NewByName("demo.A")
				_, _ = r.NewByTypeURL("type.googleapis.com/demo.B")
			}
		}()
	}
	wg.Wait()
}

func TestRegistryAllSortedAndCopied(t *testing.T) {
	r := NewRegistry()
	r.Register(Entry{FullName: "demo.B", New: func() protobuf_go_lite.Message { return &stubMsg{} }})
	r.Register(Entry{
		FullName: "demo.A",
		New:      func() protobuf_go_lite.Message { return &stubMsg{} },
		Options:  []Option{{Name: "k", Values: []string{"v"}}},
	})
	all := r.All()
	if all[0].FullName != "demo.A" || all[1].FullName != "demo.B" {
		t.Fatalf("unsorted: %#v", all)
	}
	all[0].Options[0].Values[0] = "mutated"
	all2 := r.All()
	if all2[0].Options[0].Values[0] != "v" {
		t.Fatal("All should return copied metadata")
	}
}
