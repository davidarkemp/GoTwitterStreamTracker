package main

import (
	"math/rand"
	"time"
)

type RandomNumberGenerator interface {
	Next(uint64) uint64
}

type PredefinedRandom struct {
	numbers []uint64
	current uint64
}

func (p *PredefinedRandom) Next(top uint64) uint64 {
	index := p.current
	p.current += 1
	return p.numbers[index]
}

func NewPredefinedRandom(numbers []uint64) *PredefinedRandom {
	return &PredefinedRandom{numbers: numbers}
}

type PseudoRangomNumberGenerator struct {
}

func (p *PseudoRangomNumberGenerator) Next(top uint64) uint64 {
	return uint64(rand.Int63n(int64(top)))
}

func NewPseudoRangomNumberGenerator() *PseudoRangomNumberGenerator {
	rand.Seed(time.Now().UTC().UnixNano())
	return &PseudoRangomNumberGenerator{}
}
