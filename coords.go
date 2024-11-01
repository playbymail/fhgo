// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fhgo

import (
	"fmt"
	"math"
)

type coord_t struct {
	x int
	y int
	z int
}

func (c coord_t) DistanceTo(o coord_t) float64 {
	dx, dy, dz := float64(c.x-o.x), float64(c.y-o.y), float64(c.z-o.z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (c coord_t) IsZero() bool {
	return c.x == 0 && c.y == 0 && c.z == 0
}

func (c coord_t) Less(o coord_t) bool {
	di, dj := c.DistanceTo(coord_t{}), o.DistanceTo(coord_t{})
	if di < dj {
		return true
	} else if di == dj {
		if c.x < o.x {
			return true
		} else if c.x == o.x {
			if c.y < o.y {
				return true
			} else if c.y == o.y {
				return c.z < o.z
			}
		}
	}
	return false
}

func (c coord_t) String() string {
	return fmt.Sprintf("(%5d,%5d,%5d)", c.x, c.y, c.z)
}
