// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fhgo

import "github.com/playbymail/fhgo/prng"

// rnd returns a random int between 1 and max, inclusive.
// It uses the so-called "Algorithm M" method, which is a combination of the congruential and shift-register methods.
func rnd(max int) int {
	return prng.Rand(max)
}
