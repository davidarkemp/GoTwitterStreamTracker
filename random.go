package main

import (
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

func NewPredefinedRandom(numbers []uint64) *PredefinedRandom{
	return &PredefinedRandom{ numbers: numbers }
}
