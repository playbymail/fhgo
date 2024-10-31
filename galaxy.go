// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fhgo

import (
	"fmt"
	"log"
	"math"
	"time"
)

type GalaxyData = galaxy_data_t

type galaxy_data_t struct {
	d_num_species int // Design number of species in galaxy
	num_species   int // Actual number of species allocated
	radius        int // Galactic radius in parsecs
	turn_number   int // Current turn number
}

func CreateGalaxy(path string, galacticRadius, desiredNumStars, desiredNumSpecies int, seed uint64) *GalaxyData {
	if galacticRadius < MIN_RADIUS || galacticRadius > MAX_RADIUS {
		log.Fatalf("error: galaxy must have a radius between %d and %d parsecs.\n", MIN_RADIUS, MAX_RADIUS)
	}
	if desiredNumStars < MIN_STARS || desiredNumStars > MAX_STARS {
		log.Fatalf("error: galaxy must have between %d and %d star systems\n", MIN_STARS, MAX_STARS)
	}
	if desiredNumSpecies < MIN_SPECIES || desiredNumSpecies > MAX_SPECIES {
		log.Fatalf("error: galaxy must have between %d and %d species\n", MIN_SPECIES, MAX_SPECIES)
	}

	fmt.Printf(" info: radius      %6d\n", galacticRadius)
	fmt.Printf(" info: stars       %6d\n", desiredNumStars)
	fmt.Printf(" info: species     %6d\n", desiredNumSpecies)

	// The probability of a star system existing at any particular set of x,y,z coordinates
	// is the volume of the cluster divided by the desired number of stars.
	volume := 4 * math.Pi * math.Pow(float64(galacticRadius), 3) / 3
	starsPerCubicParsec := float64(desiredNumStars) / volume
	fmt.Printf("       volume of cluster  == %12.5f cubic parsecs\n", volume)
	fmt.Printf("       density of cluster == %12.5f stars per cubic parsec\n", starsPerCubicParsec)
	fmt.Printf("       minimum density    == %12.5f\n", 1.0/3200.0)
	fmt.Printf("       maximum density    == %12.5f\n", 1.0/50.0)

	/* Get the number of cubic parsecs within a sphere with a radius of galacticRadius parsecs.
	 * Again, use long values to prevent loss of data by compilers that use 16-bit ints. */
	galactic_diameter := 2 * galacticRadius
	galactic_volume := (4 * 314 * galacticRadius * galacticRadius * galacticRadius) / 300

	// The chance_of_star is the probability of a star system existing at any particular set of x,y,z coordinates.
	// It's the volume of the cluster divided by the desired number of stars.
	chance_of_star := galactic_volume / desiredNumStars
	fmt.Printf("       galactic_volume    == %6d cubic parsecs\n", galactic_volume)
	fmt.Printf("       desiredNumStars    == %6d stars\n", desiredNumStars)
	fmt.Printf("       chance_of_star     == %6d\n", chance_of_star)
	if chance_of_star < 50 {
		log.Fatalf("error: galactic radius is too small for %d stars\n", desiredNumStars)
	} else if chance_of_star > 3200 {
		log.Fatalf("error: galactic radius is too large for %d stars\n", desiredNumStars)
	}

	type coord_t struct {
		x int
		y int
		z int
	}

	started := time.Now()

	// initialize star location data
	var starCoords []coord_t
	var star_here [MAX_DIAMETER][MAX_DIAMETER]bool

	// randomly assign stars to locations within the galactic cluster.
	galacticRadiusSquared := galacticRadius * galacticRadius
	for len(starCoords) < desiredNumStars {
		// randomly place a star
		coords := coord_t{
			x: rnd(galactic_diameter) - 1 - galacticRadius,
			y: rnd(galactic_diameter) - 1 - galacticRadius,
			z: rnd(galactic_diameter) - 1 - galacticRadius,
		}

		// if there's already a star at this x, y, loop and try again.
		// yeah, we allow only one star per x,y coordinate.
		if star_here[coords.x+galacticRadius][coords.y+galacticRadius] {
			continue
		}

		// check that the coordinate is within the galactic boundary.
		// if it isn't, loop back to the top of the loop and try again.
		sq_distance_from_center := (coords.x * coords.x) + (coords.y * coords.y) + (coords.z * coords.z)
		if sq_distance_from_center >= galacticRadiusSquared {
			continue
		}
		// otherwise, add the star to the list of stars.
		starCoords = append(starCoords, coords)
		// and mark the location as having a star.
		star_here[coords.x+galacticRadius][coords.y+galacticRadius] = true
	}
	fmt.Printf("       number of stars    == %6d in %v\n", len(starCoords), time.Since(started))

	g := &galaxy_data_t{
		d_num_species: desiredNumSpecies,
		num_species:   desiredNumSpecies,
		radius:        galacticRadius,
		turn_number:   0,
	}

	return g
}
