// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fhgo

import (
	"fmt"
	"log"
	"math"
	"sort"
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

	started := time.Now()

	// initialize star location data
	var starList []coord_t
	var star_here [MAX_DIAMETER][MAX_DIAMETER]bool
	origin := coord_t{x: 0, y: 0, z: 0}

	// randomly assign stars to locations within the galactic cluster.
	maxDistance := float64(galacticRadius)
	for len(starList) < desiredNumStars {
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
		if origin.DistanceTo(coords) > maxDistance {
			continue
		}
		// otherwise, add the star to the list of stars.
		starList = append(starList, coords)
		// and mark the location as having a star.
		star_here[coords.x+galacticRadius][coords.y+galacticRadius] = true
	}
	fmt.Printf("       number of stars    == %6d in %v\n", len(starList), time.Since(started))

	sort.Slice(starList, func(i, j int) bool {
		return starList[i].Less(starList[j])
	})

	g := &galaxy_data_t{
		d_num_species: desiredNumSpecies,
		num_species:   desiredNumSpecies,
		radius:        galacticRadius,
		turn_number:   0,
	}

	for n, coords := range starList {
		fmt.Printf("star %6d: %s %12.4f\n", n+1, coords.String(), origin.DistanceTo(coords))
		star := &star_data_t{
			x: coords.x,
			y: coords.y,
			z: coords.z,
		}

		// Determine type of star. Make MAIN_SEQUENCE the most common star type.
		// Type of star determines number of dice rolled when generating planets.
		var numberOfDice int
		switch rnd(10) {
		case 1:
			star.type_ = DWARF
			numberOfDice = 1
		case 2:
			star.type_ = DEGENERATE
			numberOfDice = 2
		case 3:
			star.type_ = GIANT
			numberOfDice = 3
		default:
			star.type_ = MAIN_SEQUENCE
			numberOfDice = 2
		}

		// Color of star is totally random and influences the number of dice rolled when generating planets.
		// Big stars (blue, blue-white) roll bigger dice. Smaller stars (orange, red) roll smaller dice.
		var planetDiceSize int
		switch rnd(7) {
		case 1:
			star.color = BLUE
			planetDiceSize = 7 + 2 - 1 // RED + 2 - star.color
		case 2:
			star.color = BLUE_WHITE
			planetDiceSize = 7 + 2 - 2 // RED + 2 - star.color
		case 3:
			star.color = WHITE
			planetDiceSize = 7 + 2 - 3 // RED + 2 - star.color
		case 4:
			star.color = YELLOW_WHITE
			planetDiceSize = 7 + 2 - 4 // RED + 2 - star.color
		case 5:
			star.color = YELLOW
			planetDiceSize = 7 + 2 - 5 // RED + 2 - star.color
		case 6:
			star.color = ORANGE
			planetDiceSize = 7 + 2 - 6 // RED + 2 - star.color
		case 7:
			star.color = RED
			planetDiceSize = 7 + 2 - 7 // RED + 2 - star.color
		default:
			panic("galaxy.go: unknown star color")
		}

		/* Size of star is totally random. */
		star.size = rnd(10) - 1

		/* Determine the number of planets in orbit around the star.
		 * The algorithm is something I tweaked until I liked it.
		 * It's weird, but it works. */
		// start at negative 2 and add the rolls
		star.num_planets = -2
		for i := 1; i <= numberOfDice; i++ {
			star.num_planets += rnd(planetDiceSize)
		}
		// trim down if too many
		for star.num_planets > 9 {
			star.num_planets -= rnd(3)
		}
		if star.num_planets < 1 {
			star.num_planets = 1
		}
		fmt.Printf("star %6d: %s %12.4f planets %2d\n", n+1, coords.String(), origin.DistanceTo(coords), star.num_planets)

		_ = star
	}

	return g
}
