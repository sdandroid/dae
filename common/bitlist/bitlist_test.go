/*
 * SPDX-License-Identifier: AGPL-3.0-only
 * Copyright (c) 2022-2025, daeuniverse Organization <dae@v2raya.org>
 */

package bitlist

import (
	"fmt"
	"testing"
)

func TestBitList6(t *testing.T) {
	bm := NewCompactBitList(6)
	bm.Set(1, 0b110010)
	if v := bm.Get(1); v != 0b110010 {
		t.Fatal(fmt.Errorf("expect 0b%08b, got 0b%08b", 0b110010, v))
	}
	bm.Tighten()
	if v := bm.Get(1); v != 0b110010 {
		t.Fatal(fmt.Errorf("expect 0b%08b, got 0b%08b", 0b110010, v))
	}
	bm.Set(13, 0b110010)
	if v := bm.Get(13); v != 0b110010 {
		t.Fatal(fmt.Errorf("expect 0b%08b, got 0b%08b", 0b110010, v))
	}
	bm.Tighten()
	if bm.b.Cap() != 6 {
		t.Fatal("failed to tighten", bm.b.Cap())
	}
	if v := bm.Get(13); v != 0b110010 {
		t.Fatal(fmt.Errorf("expect 0b%08b, got 0b%08b", 0b110010, v))
	}
	bm.Append(0b110010)
	if v := bm.Get(14); v != 0b110010 {
		t.Fatal(fmt.Errorf("expect 0b%08b, got 0b%08b", 0b110010, v))
	}
	if bm.b.Cap() != 6 {
		t.Fatal("unexpected grow behavior", bm.b.Cap())
	}
	bm.Tighten()
	if bm.b.Cap() != 6 {
		t.Fatal("failed to tighten", bm.b.Cap())
	}
}

func TestBitList19(t *testing.T) {
	bm := NewCompactBitList(19)
	bm.Set(1, 0b1110010110010110010)
	if v := bm.Get(1); v != 0b1110010110010110010 {
		t.Fatal(fmt.Errorf("expect 0b%019b, got 0b%019b", 0b1110010110010110010, v))
	}
	bm.Tighten()
	if v := bm.Get(1); v != 0b1110010110010110010 {
		t.Fatal(fmt.Errorf("expect 0b%019b, got 0b%019b", 0b1110010110010110010, v))
	}
	bm.Set(13, 0b1110010110010110010)
	if v := bm.Get(13); v != 0b1110010110010110010 {
		t.Fatal(fmt.Errorf("expect 0b%019b, got 0b%019b", 0b1110010110010110010, v))
	}
	bm.Tighten()
	if bm.b.Cap() != 17 {
		t.Fatal("failed to tighten", bm.b.Cap())
	}
	if v := bm.Get(13); v != 0b1110010110010110010 {
		t.Fatal(fmt.Errorf("expect 0b%019b, got 0b%019b", 0b1110010110010110010, v))
	}
	bm.Append(0b1110010110010110010)
	if v := bm.Get(14); v != 0b1110010110010110010 {
		t.Fatal(fmt.Errorf("expect 0b%019b, got 0b%019b", 0b1110010110010110010, v))
	}
	if bm.b.Cap() != 35 {
		t.Fatal("unexpected grow behavior", bm.b.Cap())
	}
	bm.Tighten()
	if bm.b.Cap() != 18 {
		t.Fatal("failed to tighten", bm.b.Cap())
	}
	bm.Set(1, 0b0000000000000000000)
	if v := bm.Get(1); v != 0b0000000000000000000 {
		t.Fatal(fmt.Errorf("expect 0b%019b, got 0b%019b", 0b0000000000000000000, v))
	}
	bm.Set(2, 0b1111111111111111111)
	if v := bm.Get(2); v != 0b1111111111111111111 {
		t.Fatal(fmt.Errorf("expect 0b%019b, got 0b%019b", 0b1111111111111111111, v))
	}
	if v := bm.Get(1); v != 0b0000000000000000000 {
		t.Fatal(fmt.Errorf("expect 0b%019b, got 0b%019b", 0b0000000000000000000, v))
	}
}

func TestNewCompactBitListWithCapacity(t *testing.T) {
	bm := NewCompactBitListWithCapacity(6, 100)
	if bm.b.Cap() != 38 {
		t.Fatal("unexpected preallocated capacity", bm.b.Cap())
	}
}

func TestNewCompactBitListWithSize(t *testing.T) {
	bm := NewCompactBitListWithSize(6, 100)
	if bm.b.Len() != 38 {
		t.Fatal("unexpected preallocated size", bm.b.Len())
	}
	for i := 0; i < 100; i++ {
		bm.Append(uint64(i % 64))
	}
	for i := 0; i < 100; i++ {
		if got, want := bm.Get(i), uint64(i%64); got != want {
			t.Fatalf("index %d: expect %d, got %d", i, want, got)
		}
	}
}

func TestAppend(t *testing.T) {
	for bitSize := 1; bitSize <= 64; bitSize++ {
		bm := NewCompactBitListWithCapacity(bitSize, 100)
		mask := ^uint64(0)
		if bitSize < 64 {
			mask = 1<<bitSize - 1
		}
		for i := 0; i < 100; i++ {
			bm.Append(uint64(i) & mask)
		}
		for i := 0; i < 100; i++ {
			if got, want := bm.Get(i), uint64(i)&mask; got != want {
				t.Fatalf("bit size %d index %d: expect %d, got %d", bitSize, i, want, got)
			}
		}
	}
}
