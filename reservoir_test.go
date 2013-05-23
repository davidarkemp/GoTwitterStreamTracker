package main

import (
	"testing"
)

func TestReservoir(t *testing.T) {
	rng := NewPredefinedRandom([]uint64{2, 3, 1, 0})
	rsv := NewReservoirSampler(2, rng)

	if !rsv.Add("first") {
		t.Error("can't add first element")
	}
	t.Logf("samples are %v", rsv.GetSamples())
	if !rsv.Add("second") {
		t.Error("can't add second element")
	}

	if rsv.Add("third") {
		t.Error("third added erroneously")
	}
	if rsv.Add("fourth") {
		t.Error("fourth added erroneously")
	}

	if !rsv.Add("fifth") {
		t.Error("didn't add fifth element")
	}
	if !rsv.Add("sixth") {
		t.Error("didn't add sixth element")
	}

	samples := rsv.GetSamples()
	if samples[0] != "sixth" {
		t.Error("sampling went wrong, we got %v and not 'sixth'", samples)
	}
}
