// Package registry indexes generated protobuf-go-lite message constructors and
// their statically generated custom options.
package registry

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
	"sync"

	protobuf_go_lite "github.com/aperturerobotics/protobuf-go-lite"
	anypb_resolver "github.com/aperturerobotics/protobuf-go-lite/types/known/anypb/resolver"
)

// ErrNotFound is returned when a message type is not registered.
var ErrNotFound = anypb_resolver.ErrNotFound

// Option contains the flattened values of one custom message option.
type Option struct {
	// Name is the fully qualified option name followed by its message field path.
	Name string
	// Values contains repeated values in declaration order.
	Values []string
}

// Entry describes one generated protobuf message type.
type Entry struct {
	// FullName is the fully qualified protobuf message name.
	FullName string
	// TypeURL is the canonical google.protobuf.Any type URL.
	TypeURL string
	// New constructs an empty message.
	New func() protobuf_go_lite.Message
	// Options contains the message's flattened custom options.
	Options []Option
}

// Option returns the first value registered under name.
func (e Entry) Option(name string) (string, bool) {
	for _, opt := range e.Options {
		if opt.Name == name && len(opt.Values) != 0 {
			return opt.Values[0], true
		}
	}
	return "", false
}

// Registry indexes message constructors by full name and type URL.
type Registry struct {
	mu      sync.RWMutex
	byName  map[string]Entry
	byURL   map[string]Entry
	entries []Entry
}

var _ anypb_resolver.AnyTypeResolver = (*Registry)(nil)

// NewRegistry constructs an empty message registry.
func NewRegistry() *Registry {
	return &Registry{
		byName: make(map[string]Entry),
		byURL:  make(map[string]Entry),
	}
}

// Register adds e to the registry. It copies the option metadata and panics if
// FullName is empty, New is nil, or FullName is already registered.
func (r *Registry) Register(e Entry) {
	if e.FullName == "" {
		panic("registry: Entry.FullName is required")
	}
	if e.New == nil {
		panic("registry: Entry.New is required")
	}
	if e.TypeURL == "" {
		e.TypeURL = "type.googleapis.com/" + e.FullName
	}

	e = copyEntry(e)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byName[e.FullName]; ok {
		panic("registry: duplicate registration for " + strconv.Quote(e.FullName))
	}
	r.byName[e.FullName] = e
	r.byURL[e.TypeURL] = e
	r.byURL[e.FullName] = e
	r.entries = append(r.entries, e)
	slices.SortFunc(r.entries, func(a, b Entry) int {
		return cmp.Compare(a.FullName, b.FullName)
	})
}

// All returns a full-name-sorted snapshot of the registered entries.
func (r *Registry) All() []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Entry, len(r.entries))
	for i, e := range r.entries {
		out[i] = copyEntry(e)
	}
	return out
}

// Range calls fn for each entry in full-name order until fn returns false.
// The callback runs without holding the registry lock.
func (r *Registry) Range(fn func(Entry) bool) {
	for _, e := range r.All() {
		if !fn(e) {
			return
		}
	}
}

// NewByName constructs the message registered under its full protobuf name.
func (r *Registry) NewByName(name string) (protobuf_go_lite.Message, bool) {
	r.mu.RLock()
	e, ok := r.byName[name]
	r.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return e.New(), true
}

// NewByTypeURL constructs the message named by url.
func (r *Registry) NewByTypeURL(url string) (protobuf_go_lite.Message, bool) {
	r.mu.RLock()
	e, ok := r.lookupByTypeURLLocked(url)
	r.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return e.New(), true
}

// FindMessageByURL returns the constructor named by url.
func (r *Registry) FindMessageByURL(url string) (func() protobuf_go_lite.Message, error) {
	r.mu.RLock()
	e, ok := r.lookupByTypeURLLocked(url)
	r.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	newFn := e.New
	return newFn, nil
}

func (r *Registry) lookupByTypeURLLocked(url string) (Entry, bool) {
	if e, ok := r.byURL[url]; ok {
		return e, true
	}
	name := url
	if i := strings.LastIndex(url, "/"); i >= 0 {
		name = url[i+1:]
	}
	e, ok := r.byName[name]
	return e, ok
}

func copyEntry(e Entry) Entry {
	out := Entry{
		FullName: e.FullName,
		TypeURL:  e.TypeURL,
		New:      e.New,
	}
	if len(e.Options) > 0 {
		out.Options = make([]Option, len(e.Options))
		for i, opt := range e.Options {
			out.Options[i] = Option{
				Name:   opt.Name,
				Values: slices.Clone(opt.Values),
			}
		}
	}
	return out
}

var defaultRegistry = NewRegistry()

// Register adds e to the default registry.
func Register(e Entry) { defaultRegistry.Register(e) }

// All returns a snapshot of the default registry.
func All() []Entry { return defaultRegistry.All() }

// Range calls fn for each entry in the default registry.
func Range(fn func(Entry) bool) { defaultRegistry.Range(fn) }

// NewByName constructs a message from the default registry by full name.
func NewByName(name string) (protobuf_go_lite.Message, bool) {
	return defaultRegistry.NewByName(name)
}

// NewByTypeURL constructs a message from the default registry by type URL.
func NewByTypeURL(url string) (protobuf_go_lite.Message, bool) {
	return defaultRegistry.NewByTypeURL(url)
}

// Default returns the registry populated by generated init functions.
func Default() *Registry { return defaultRegistry }
