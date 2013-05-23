package main

import "sync"

type Reservoir struct {
	samples []interface{}
	total   uint64
	rng     RandomNumberGenerator

	rwMutex sync.RWMutex
}

func (r *Reservoir) Add(sample interface{}) bool {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	r.total += 1
	s_len, s_cap := len(r.samples), cap(r.samples)
	if s_len < s_cap {
		r.samples = append(r.samples, sample)
		return true
	}

	random := r.rng.Next(r.total)
	if random >= uint64(cap(r.samples)) {
		return false
	}
	r.samples[random] = sample
	return true
}

func (r *Reservoir) GetSamples() []interface{} {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	snapshot := make([]interface{}, len(r.samples))
	copy(snapshot, r.samples)
	return snapshot
}

func NewReservoirSampler(size int64, rng RandomNumberGenerator) *Reservoir {
	return &Reservoir{samples: make([]interface{}, 0, size), rng: rng}
}
