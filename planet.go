// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fhgo

type planet_data_t struct {
	id                planet_id_t      // unique identifier for this planet
	index             int              // index of this planet into the planet_base array
	temperature_class int              // Temperature class, 1-30
	pressure_class    int              // Pressure class, 0-29
	special           planet_special_e // 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole
	gas               [4]gas_e         // Gas in atmosphere. Zero if none
	gas_percent       [4]int           // Percentage of gas in atmosphere
	diameter          int              // Diameter in thousands of kilometers
	gravity           int              // Surface gravity. Multiple of Earth gravity times 100
	mining_difficulty int              // Mining difficulty times 100
	econ_efficiency   int              // Economic efficiency. Always 100 for a home planet
	md_increase       int              // Increase in mining difficulty
	message           int              // Message associated with this planet, if any
	isValid           bool             // FALSE if the record is invalid
	star              *star_data_t     // pointer to the star the planet is orbiting
	orbit             int              // orbit of planet in the system
}
type planet_id_t int

// generate_planets creates planets and inserts them into the planet_data array.
// returns a slice of planet_data_t pointers and a flag indicating if the planet is a potential home system.
//
// note that the potential home system is always false if the caller set earth_like to false.
func generate_planets(star *star_data_t, num_planets int, earth_like, makeMiningEasier bool) ([10]*planet_data_t, bool) {
	/* Values for the planets of Earth's solar system will be used as starting values.
	 * Diameters are in thousands of kilometers.
	 * The zeroth element of each array is a placeholder and is not used.
	 * The fifth element corresponds to the asteroid belt, and is pure fantasy on my part.
	 * I omitted Pluto because it is probably a captured planet, rather than an original member of our solar system. */
	//start_diameter := []int{0, 5, 12, 13, 7, 20, 143, 121, 51, 49}
	//start_temp_class := []int{0, 29, 27, 11, 9, 8, 6, 5, 5, 3}
	start := [10]struct {
		diameter, temp_class int
	}{
		{}, // ignored
		{diameter: 5, temp_class: 29},
		{diameter: 12, temp_class: 27},
		{diameter: 13, temp_class: 11},
		{diameter: 7, temp_class: 9},
		{diameter: 20, temp_class: 8}, // fantasy asteroid belt
		{diameter: 143, temp_class: 6},
		{diameter: 121, temp_class: 5},
		{diameter: 51, temp_class: 5},
		{diameter: 49, temp_class: 3},
	}

	// we need a temporary place to store values for each planet because the logic
	// will make updates to neighboring planets as it progresses through the system.
	var planet_values [10]struct {
		diameter, temperature_class int
		gas_giant                   bool
		density                     int
		gravity                     int
		atmosphere                  [5]struct {
			gas     gas_e
			percent int
		}
		mining_difficulty, pressure_class int
		special                           planet_special_e
	}

	/* Set flag to indicate if this star system requires an earth-like planet.
	 * If so, we will zero this flag after we use it. */
	make_earth := earth_like

	planets := [10]*planet_data_t{}

	/* Main loop. Generate one planet at a time. */
	for planet_number := 1; planet_number <= num_planets; planet_number++ {
		pv, ppv := &planet_values[planet_number], &planet_values[planet_number-1]
		if planet_number == 1 {
			ppv = nil
		}

		/* Start with diameters, temperature classes and pressure classes based on the planets in Earth's solar system. */

		// nudge the planet starting point towards the Earth-like zone
		basePlanetTemplate := (9 * planet_number) / num_planets
		if num_planets <= 3 {
			basePlanetTemplate = 2*planet_number + 1
		}

		pv.diameter, pv.temperature_class = start[basePlanetTemplate].diameter, start[basePlanetTemplate].temp_class

		// Randomize the diameter.
		// Minimum allowable diameter is 3,000 km.
		// Note that the maximum diameter we can generate is 283,000 km.
		die_size := pv.diameter / 4
		if die_size < 2 {
			die_size = 2
		}
		for i := 1; i <= 4; i++ {
			roll := rnd(die_size)
			if rnd(100) > 50 {
				pv.diameter += roll
			} else {
				pv.diameter -= roll
			}
		}
		for pv.diameter < 3 {
			pv.diameter += rnd(4)
		}

		// if diameter is greater than 40,000 km, assume the planet is a gas giant
		pv.gas_giant = pv.diameter > 40

		/* Density will depend on whether the planet is a gas giant.
		 * Again ignoring Pluto, densities range from 0.7 to 1.6 times the
		 * density of water for the gas giants, and from 3.9 to 5.5 for the
		 * others. We will expand this range slightly and use 100 times the
		 * actual density so that we can use integer arithmetic. */
		if pv.gas_giant {
			/* Final values from 0.60 through 1.70 (scaled to 60 through 170). */
			pv.density = 58 + rnd(56) + rnd(56)
		} else {
			/* Final values from 3.70 through 5.70 (scaled to 370 through 570). */
			pv.density = 368 + rnd(101) + rnd(101)
		}

		/* Gravitational acceleration is proportional to the mass divided by the radius-squared.
		 * The radius is proportional to the diameter, and the mass is proportional to the
		 * density times the radius-cubed. The net result is that "g" is proportional to
		 * the density times the diameter. Our value for "g" will be a multiple of Earth
		 * gravity, and will be further multiplied by 100 to allow us to use integer arithmetic.
		 *
		 * The factor 72 ensures that "g" will be 100 for Earth (density=550 and diameter=13). */
		pv.gravity = (pv.density * pv.diameter) / 72

		/* Randomize the temperature class obtained earlier. */
		die_size = pv.temperature_class / 4
		if die_size < 2 {
			die_size = 2
		}
		n_rolls := rnd(3) + rnd(3) + rnd(3)
		for i := 1; i <= n_rolls; i++ {
			roll := rnd(die_size)
			if rnd(100) > 50 {
				pv.temperature_class += roll
			} else {
				pv.temperature_class -= roll
			}
		}
		// adjust the temperature class for gas giants and small planets
		if pv.gas_giant {
			// nudge the temperature class towards the gas giant zone, 3 through 7
			for pv.temperature_class < 3 {
				pv.temperature_class += rnd(2)
			}
			for pv.temperature_class > 7 {
				pv.temperature_class -= rnd(2)
			}
		} else {
			// nudge the temperature class towards the small planet zone, 1 through 30
			for pv.temperature_class < 1 {
				pv.temperature_class += rnd(3)
			}
			for pv.temperature_class > 30 {
				pv.temperature_class -= rnd(3)
			}
		}

		/* Sometimes, planets close to the sun in star systems with less than four planets are too cold.
		 * Warm them up a little. */
		if num_planets < 4 && planet_number < 3 {
			for pv.temperature_class < 12 {
				pv.temperature_class += rnd(4)
			}
		}
		/* Make sure that planets farther from the sun are not warmer than planets closer to the sun. */
		if ppv != nil && ppv.temperature_class < pv.temperature_class {
			pv.temperature_class = ppv.temperature_class
		}

		/* Check if this planet should be earth-like.
		 * If so, replace all the above with earth-like characteristics. */
		// BUG: it doesn't actually replace all the above; that check for warmer planets is missing
		if make_earth && (pv.temperature_class <= 11) {
			make_earth = false // do this for only one planet per system

			// make attributes earth-like
			pv.diameter = 11 + rnd(3)
			pv.gravity = 93 + rnd(11) + rnd(11) + rnd(5)
			pv.temperature_class = 9 + rnd(3)
			pv.pressure_class = 8 + rnd(3)
			pv.mining_difficulty = 208 + rnd(11) + rnd(11)
			pv.special = IDEAL_HOME_PLANET /* Maybe ideal home planet. */

			// make some earth-like atmospheric gases
			i, total_percent := 0, 0
			// 33% chance that it has up to 30% ammonia
			if rnd(3) == 1 {
				pct := rnd(30)
				pv.atmosphere[i].gas = NH3
				pv.atmosphere[i].percent = pct
				i, total_percent = i+1, total_percent+pct
			}
			// always at least 10 % nitrogen
			nitro := i
			pct := 10
			pv.atmosphere[i].gas = N2
			pv.atmosphere[i].percent = pct
			i, total_percent = i+1, total_percent+pct
			i = i + 1
			// 33% chance that it has up to 30% carbon dioxide
			if rnd(3) == 1 {
				pct := rnd(30)
				pv.atmosphere[i].gas = CO2
				pv.atmosphere[i].percent = pct
				i, total_percent = i+1, total_percent+pct
			}
			// always 10% to 30% oxygen
			pct = rnd(20) + 10
			pv.atmosphere[i].gas = O2
			pv.atmosphere[i].percent = pct
			total_percent += pct
			// the remainder must be allocated to nitrogen
			pv.atmosphere[nitro].percent += 100 - total_percent

			// this planet is now earth-like; move on to the next planet
			continue
		}

		/* Pressure class depends primarily on gravity.
		 * Calculate an approximate value and randomize it. */
		pv.pressure_class = pv.gravity / 10
		die_size = pv.pressure_class / 4
		if die_size < 2 {
			die_size = 2
		}
		n_rolls = rnd(3) + rnd(3) + rnd(3)
		for i := 1; i <= n_rolls; i++ {
			roll := rnd(die_size)
			if rnd(100) > 50 {
				pv.pressure_class += roll
			} else {
				pv.pressure_class -= roll
			}
		}
		if pv.gas_giant {
			for pv.pressure_class < 11 {
				pv.pressure_class += rnd(3)
			}
			for pv.pressure_class > 29 {
				pv.pressure_class -= rnd(3)
			}
		} else {
			for pv.pressure_class < 0 {
				pv.pressure_class += rnd(3)
			}
			for pv.pressure_class > 12 {
				pv.pressure_class -= rnd(3)
			}
		}
		if pv.gravity < 10 {
			/* Planet's gravity is too low to retain an atmosphere. */
			pv.pressure_class = 0
		} else if pv.temperature_class < 2 || pv.temperature_class > 27 {
			/* Planets outside this temperature range have no atmosphere. */
			pv.pressure_class = 0
		}

		/* Generate gases, if any, in the atmosphere. */
		if pv.pressure_class != 0 {
			/* Convert planet's temperature class to a value between 1 and 9.
			 * We will use it as the start index into the list of 13 potential gases. */
			var atmospheric_gases []gas_e
			if first_gas := 100 * pv.temperature_class / 225; first_gas < 1 {
				atmospheric_gases = all_gases[1:]
			} else if first_gas < 9 {
				atmospheric_gases = all_gases[first_gas:]
			} else {
				atmospheric_gases = all_gases[9:]
			}

			// the following algorithm is something I tweaked until it worked well.
			// number of gases to generate is 2d4 divided by 2, rounded down.
			num_gases_remaining, num_gases_found := (rnd(4)+rnd(4))/2, 0
			total_gas_quantity := 0

			// it is important to limit the maximum number of atmospheric gases to 4.
			// if we don't, it will make the life support numbers for alien planets unrealistically high.
			for _, gas := range atmospheric_gases[:4] {
				if num_gases_remaining == 0 {
					break
				}
				percent := 0
				if gas == HE { // treat Helium specially
					// just a third of the very coldest planets will actually have He
					if pv.temperature_class > 5 || rnd(3) != 1 {
						continue
					}
					percent = rnd(20)
				} else { // all other gases
					// a third of the remaining gases will be silently ignored, helps with LSN calculations
					if rnd(3) == 3 {
						continue
					}
					percent = rnd(100)
					if gas == O2 {
						// Oxygen is self-limiting
						percent = (percent + 1) / 2
					}
				}
				if percent == 0 {
					continue
				}
				num_gases_found++
				num_gases_remaining--
				pv.atmosphere[num_gases_found].gas = gas
				pv.atmosphere[num_gases_found].percent = percent
				total_gas_quantity += percent
			}
			// convert gas quantities to percentages.
			// it is okay for us to go through all the slots because Go initializes them to 0
			if total_gas_quantity > 0 {
				total_percent := 0
				for i := range pv.atmosphere {
					pv.atmosphere[i].percent = 100 * pv.atmosphere[i].percent / total_gas_quantity
					total_percent += pv.atmosphere[i].percent
				}
				// give leftover to first gas
				pv.atmosphere[1].percent += 100 - total_percent
			}
		}

		// Get mining difficulty.
		// Mining difficulty is proportional to planetary diameter with randomization and an occasional big surprise.
		// Default mining difficulty values will range between 0.80 and 10.00.
		// The actual value will be scaled by 100 to allow use of integer arithmetic.
		minMiningDifficulty, maxMiningDifficulty, surpriseFactor := 40, 500, 30
		if makeMiningEasier {
			// Tweak values to move the range to between 0.30 and 10.00.
			minMiningDifficulty, maxMiningDifficulty, surpriseFactor = 30, 1000, 20
		}
		for pv.mining_difficulty < minMiningDifficulty || pv.mining_difficulty > maxMiningDifficulty {
			pv.mining_difficulty = (rnd(3)+rnd(3)+rnd(3)-rnd(4))*rnd(pv.diameter) + rnd(surpriseFactor) + rnd(surpriseFactor)
		}
		if !makeMiningEasier {
			pv.mining_difficulty = (pv.mining_difficulty * 11) / 5 // fudge factor to make things harder
		}
	}

	/* Copy planet data to structure. */
	potentialHomeSystem := false
	var home_planet *planet_data_t
	for i := 1; i <= num_planets; i++ {
		pv := &planet_values[i]
		current_planet := &planet_data_t{
			diameter:          pv.diameter,
			gravity:           pv.gravity,
			mining_difficulty: pv.mining_difficulty,
			temperature_class: pv.temperature_class,
			pressure_class:    pv.pressure_class,
			special:           pv.special,
		}
		for n := 0; n < 4; n++ {
			current_planet.gas[n] = pv.atmosphere[n+1].gas
			current_planet.gas_percent[n] = pv.atmosphere[n+1].percent
		}

		if current_planet.special == IDEAL_HOME_PLANET {
			potentialHomeSystem = true
			home_planet = current_planet
		}

		planets[i] = current_planet
	}

	// note that potentialHomeSystem and home_planet are set by the loop above.
	// it turns out that they're only set if the caller set the earth_like flag to true.
	if potentialHomeSystem && home_planet != nil {
		// the planets in potential home systems must also have a viable atmosphere and economic prospects.
		// we have a couple of magic tests that will (should?) ensure that. if the tests fail, we will unset
		// the potentialHomeSystem flag and the home_planet pointer.
		//
		// check the mining potential of all planets in the system relative to the designated home planet.
		// calculate the score by:
		//
		//  1. taking each planet in the system
		//  2. computing its LSN (Life Support Needed) distance from the home planet
		//  3. fetch the mining difficulty of each planet
		//  4. combining these into a formula: 20000 / ((3 + LSN) * (50 + mining_difficulty))
		//  5. summing up this value for all planets (saving into potential)
		//
		// the system is considered viable as a home system only if the total score falls between 53 and 57.
		// this ensures that:
		//
		//  1. the planets are close enough in terms of life support needs to be reasonably colonizable
		//  2. the mining difficulties are balanced to provide good resource gathering opportunities
		//  3. the overall system has enough accessible resources to support a starting civilization
		//
		// this calculation helps ensure new species start in systems that are both habitable and economically viable.
		//
		// note that we add O2 to the desired gases list because it is required for life and must always be present.
		potential, desiredGases := 0, append(home_planet.gas[:], O2)
		for i := 1; i <= num_planets; i++ {
			approximateLSN := planets[i].approximateLSN(home_planet.temperature_class, home_planet.pressure_class, desiredGases)
			potential += 20_000 / ((3 + approximateLSN) * (50 + planets[i].mining_difficulty))
		}

		// the system is viable only if the potential is 54, 55, or 56.
		// these values are magical and I have no idea why they work.
		potentialHomeSystem = 53 < potential && potential < 57
		if !potentialHomeSystem {
			home_planet = nil
		}
	}

	return planets, potentialHomeSystem
}

// approximateLSN is a helper function provides an approximate LSN for a planet.
// This is an approximation because it doesn't account for the species' list of neutral and poison gases.
func (p *planet_data_t) approximateLSN(temperature_class, pressure_class int, atomosphere []gas_e) int {
	// base the life support needed on differences in temperature and pressure classes
	delta_temperature := p.temperature_class - temperature_class
	if delta_temperature < 0 {
		delta_temperature = -delta_temperature
	}
	delta_pressure := p.pressure_class - pressure_class
	if delta_pressure < 0 {
		delta_pressure = -delta_pressure
	}
	ls_needed := 2*delta_temperature + 2*delta_pressure

	// compare the atmospheres of this planet against the input.
	// bump the life support needed for every gas that is in this
	// planet's atmosphere that ia missing in the other's atmosphere.
	for _, gas := range p.gas {
		// skip if no gas in this slot
		if gas == GAS_NONE {
			continue
		}
		// check if the gas is in the given atmosphere
		hasGas := false
		for _, atmo_gas := range atomosphere {
			if gas == atmo_gas {
				hasGas = true
				break
			}
		}
		if !hasGas {
			ls_needed += 2
		}
	}

	return ls_needed
}

var potential_home_system = false
