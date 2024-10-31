// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package prng implements a PRNG to return a random int between 1 and max, inclusive.
// It uses the so-called "Algorithm M" method, which is a combination of the congruential and shift-register methods.
package prng

import (
	"fmt"
	"strconv"
)

type PRNG struct {
	seed uint64
}

var defaultHistoricalSeedValue uint64 = 1924085713

var defaultPRNG = &PRNG{
	seed: defaultHistoricalSeedValue,
}

func DefaultHistoricalSeedValue() uint64 {
	return defaultHistoricalSeedValue
}

func DefaultPRNG() *PRNG {
	return defaultPRNG
}

// Rand returns a random int between 1 and max, inclusive.
func Rand(max int) int {
	return defaultPRNG.IntN(max)
}

func SetSeed(seed uint64) {
	defaultPRNG.seed = seed
}

func String() string {
	return fmt.Sprintf("%016x", defaultPRNG.seed)
}

func New(seed uint64) *PRNG {
	return &PRNG{seed: seed}
}

// IntN returns a random int between 1 and max, inclusive.
func (p *PRNG) IntN(max int) int {
	// for congruential method, multiply previous value by the prime number 16417.
	var congResult uint64 = p.seed + (p.seed << 5) + p.seed<<14

	// for shift-register method, use shift-right 15 and shift-left 17 with no-carry addition (i.e., exclusive-or)
	var shiftResult uint64 = (p.seed >> 15) ^ p.seed
	shiftResult ^= (shiftResult << 17)

	// save seed for next iteration
	p.seed = congResult ^ shiftResult

	// avoid returning the low-order bits
	return int(((p.seed&0x0000FFFF)*uint64(max))>>16) + 1
}

func (p *PRNG) SetState(state string) {
	p.seed, _ = strconv.ParseUint(state, 16, 64)
}

func (p *PRNG) String() string {
	return fmt.Sprintf("%x", p.seed)
}
