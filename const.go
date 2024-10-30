// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fhgo

const (
	STANDARD_NUMBER_OF_SPECIES      = 15 /* A standard game has 15 species. */
	STANDARD_NUMBER_OF_STAR_SYSTEMS = 90 /* A standard game has 90 star-systems. */
	STANDARD_GALACTIC_RADIUS        = 20 /* A standard game has a galaxy with a radius of 20 parsecs. */

	/* Minimum and maximum values for a galaxy. */

	MAX_BATTLES              = 50 /* Maximum number of battle locations for all players. */
	MAX_DIAMETER             = MAX_RADIUS * 2
	MAX_ENGAGE_OPTIONS       = 20 /* Maximum number of engagement options that a player may specify for a single battle. */
	MAX_INTERCEPTS           = 1_000
	MAX_ITEMS                = 38 /* Always bump this up to a multiple of two. Don't forget to make room for zeroth element! */
	MAX_LOCATIONS            = 10_000
	MAX_OBS_LOCS             = 5_000
	MAX_PLANETS              = MAX_STARS * 9
	MAX_SHIPS                = 200   /* Maximum number of ships at a single battle. */
	MAX_TRANSACTIONS         = 1_000 /* Interspecies transactions. */
	MIN_RADIUS, MAX_RADIUS   = 6, 50
	MIN_SPECIES, MAX_SPECIES = 1, 100
	MIN_STARS, MAX_STARS     = 12, 1_000

	HP_AVAILABLE_POP = 1500

	NUM_EXTRA_NAMPLAS = 50  // Additional memory must be allocated for routines that name planets.
	NUM_EXTRA_PLANETS = 100 /* In case gamemaster creates new star systems with Edit program. */
	NUM_EXTRA_SHIPS   = 100 // Additional memory must be allocated for routines that build ships.
	NUM_EXTRA_STARS   = 20  /* In case gamemaster creates new star systems with Edit program. */

	TRUE  = true
	FALSE = false
)
