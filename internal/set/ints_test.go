// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package set

import (
	"math/rand/v2"
	"testing"
)

const maxLimit = 1024

var toSet, toClear [maxLimit]bool

func init() {
	var seed [32]byte

	r := rand.New(rand.NewChaCha8(seed)) //nolint:gosec
	for i := range maxLimit {
		toSet[i] = r.Int32N(2) == 0
		toClear[i] = r.Int32N(2) == 0
	}
}

func TestInts(t *testing.T) {
	ns := new(Ints)

	// Check that set starts empty.
	wantLen := 0
	if ns.Len() != wantLen {
		t.Errorf("init: Len() = %d, want %d", ns.Len(), wantLen)
	}
	for i := range maxLimit {
		if ns.Has(uint64(i)) { //nolint:gosec
			t.Errorf("init: Has(%d) = true, want false", i)
		}
	}

	// Set some numbers.
	for i, b := range toSet[:maxLimit] {
		if b {
			ns.Set(uint64(i)) //nolint:gosec
			wantLen++
		}
	}

	// Check that integers were set.
	if ns.Len() != wantLen {
		t.Errorf("after Set: Len() = %d, want %d", ns.Len(), wantLen)
	}
	for i := range maxLimit {
		if got := ns.Has(uint64(i)); got != toSet[i] { //nolint:gosec
			t.Errorf("after Set: Has(%d) = %v, want %v", i, got, !got)
		}
	}

	// Clear some numbers.
	for i, b := range toClear[:maxLimit] {
		if b {
			ns.Clear(uint64(i)) //nolint:gosec
			if toSet[i] {
				wantLen--
			}
		}
	}

	// Check that integers were cleared.
	if ns.Len() != wantLen {
		t.Errorf("after Clear: Len() = %d, want %d", ns.Len(), wantLen)
	}
	for i := range maxLimit {
		if got := ns.Has(uint64(i)); got != toSet[i] && !toClear[i] { //nolint:gosec
			t.Errorf("after Clear: Has(%d) = %v, want %v", i, got, !got)
		}
	}
}
