// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fhgo

import "math"

func distanceBetween(s1, s2 *star_data_t) float64 {
	return s1.distanceBetween(s2)
}

func (s *star_data_t) distanceBetween(s2 *star_data_t) float64 {
	return s.distanceTo(s2.x, s2.y, s2.z)
}

func (s *star_data_t) distanceTo(x, y, z int) float64 {
	dx := float64(x - s.x)
	dy := float64(y - s.y)
	dz := float64(z - s.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
