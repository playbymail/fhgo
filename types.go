// Copyright (C) 2022 Michael D Henderson. All rights reserved.

package fhgo

type action_data_t struct {
	num_units_fighting     int
	fighting_species_index [MAX_SHIPS]species_id_t
	num_shots              [MAX_SHIPS]int
	shots_left             [MAX_SHIPS]int
	weapon_damage          [MAX_SHIPS]int64
	shield_strength        [MAX_SHIPS]int64
	shield_strength_left   [MAX_SHIPS]int64
	original_age_or_PDs    [MAX_SHIPS]int64
	bomb_damage            [MAX_SHIPS]int64
	surprised              [MAX_SHIPS]byte
	unit_type              [MAX_SHIPS]byte
	fighting_unit          [MAX_SHIPS]string
}
type action_data = action_data_t

type battle_data_t struct {
	x, y, z, pn               int
	num_species_here          byte
	spec_num                  [MAX_SPECIES]species_id_t
	summary_only              [MAX_SPECIES]species_id_t
	transport_withdraw_age    [MAX_SPECIES]species_id_t
	warship_withdraw_age      [MAX_SPECIES]species_id_t
	fleet_withdraw_percentage [MAX_SPECIES]species_id_t
	haven_x                   [MAX_SPECIES]species_id_t
	haven_y                   [MAX_SPECIES]species_id_t
	haven_z                   [MAX_SPECIES]species_id_t
	special_target            [MAX_SPECIES]species_id_t
	hijacker                  [MAX_SPECIES]species_id_t
	can_be_surprised          [MAX_SPECIES]species_id_t
	enemy_mine                [MAX_SPECIES][MAX_SPECIES]species_id_t
	num_engage_options        [MAX_SPECIES]species_id_t
	engage_option             [MAX_SPECIES][MAX_ENGAGE_OPTIONS]species_id_t
	engage_planet             [MAX_SPECIES][MAX_ENGAGE_OPTIONS]species_id_t
	ambush_amount             [MAX_SPECIES]int
}
type battle_data = battle_data_t

type double = float64
type long = int

type intercept_t struct {
	x, y, z      int
	amount_spent int
}

type message_id_t int

type nampla_data_t struct {
	id             nampla_id_t     // unique identifier for this named planet
	name           string          // Name of planet
	x, y, z, pn    int             // Coordinates
	status         planet_status_e // Status of planet
	hiding         bool            // HIDE order given
	hidden         bool            // Colony is hidden
	planet_index   int             // Index (starting at zero) into the file "planets.dat" of this planet
	siege_eff      int             // Siege effectiveness - a percentage between 0 and 99
	shipyards      int             // Number of shipyards on planet
	IUs_needed     int             // Incoming ship with only CUs on board
	AUs_needed     int             // Incoming ship with only CUs on board
	auto_IUs       int             // Number of IUs to be automatically installed
	auto_AUs       int             // Number of AUs to be automatically installed
	IUs_to_install int             // Colonial mining units to be installed
	AUs_to_install int             // Colonial manufacturing units to be installed
	mi_base        int             // Mining base times 10
	ma_base        int             // Manufacturing base times 10
	pop_units      int             // Number of available population units
	item_quantity  [MAX_ITEMS]int  // Quantity of each item available
	use_on_ambush  int             // Amount to use on ambush
	message        message_id_t    // Message associated with this planet, if any
	special        int             // Different for each application
	star           *star_data_t    // pointer to system the colony is in
	planet         *planet_data_t  // pointer to planet the colony is on
}

type nampla_id_t int

type planet_data_t struct {
	id                planet_id_t  // unique identifier for this planet
	index             int          // index of this planet into the planet_base array
	temperature_class int          // Temperature class, 1-30
	pressure_class    int          // Pressure class, 0-29
	special           int          // 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole
	gas               [4]gas_e     // Gas in atmosphere. Zero if none
	gas_percent       [4]int       // Percentage of gas in atmosphere
	diameter          int          // Diameter in thousands of kilometers
	gravity           int          // Surface gravity. Multiple of Earth gravity times 100
	mining_difficulty int          // Mining difficulty times 100
	econ_efficiency   int          // Economic efficiency. Always 100 for a home planet
	md_increase       int          // Increase in mining difficulty
	message           int          // Message associated with this planet, if any
	isValid           bool         // FALSE if the record is invalid
	star              *star_data_t // pointer to the star the planet is orbiting
	orbit             int          // orbit of planet in the system
}
type planet_id_t int

type scan_system_t struct {
	star     *star_data_t
	distance float64
}

type ship_data_t struct {
	id                     ship_id_t      // unique identifier for this ship
	name                   string         // Name of ship
	x, y, z, pn            int            // Current coordinates
	status                 ship_status_e  // Current status of ship
	type_                  ship_type_e    // Ship type
	dest_x, dest_y, dest_z int            // Destination if ship was forced to jump from combat. Also used by TELESCOPE command
	just_jumped            bool           // Set if ship jumped this turn
	arrived_via_wormhole   bool           // Ship arrived via wormhole in the PREVIOUS turn
	class                  ship_class_e   // Ship class
	tonnage                int            // Ship tonnage divided by 10,000
	item_quantity          [MAX_ITEMS]int // Quantity of each item carried
	age                    int            // Ship age
	remaining_cost         int            // The cost needed to complete the ship if still under construction
	loading_point          nampla_id_t    // Nampla index for planet where ship was last loaded with CUs. Zero = none. Use 9999 for home planet
	unloading_point        nampla_id_t    // Nampla index for planet that ship should be given orders to jump to where it will unload. Zero = none. Use 9999 for home planet
	special                int            // Different for each application
}

type ship_id_t int

type ship_type_e int

type sp_loc_data_t struct {
	s       species_id_t /* Species number */
	x, y, z int
}

type species_cfg_t struct {
	email        string
	govtname     string
	govttype     string
	homeworld    string
	name         string
	ml           int
	gv           int
	ls           int
	bi           int
	experimental struct {
		econ_units   int
		make_bridges int
		ma_base      int
		mi_base      int
		ship_yards   int
		tech_bi      int
		tech_gv      int
		tech_ls      int
		tech_ma      int
		tech_mi      int
		tech_ml      int
	}
}

type species_data_t struct {
	id                 species_id_t          // unique identifier for this species
	index              int                   // index of this species in spec_data array
	name               string                // Name of species
	govt_name          string                // Name of government
	govt_type          string                // Type of government
	x, y, z, pn        int                   // Coordinates of home planet
	required_gas       gas_e                 // Gas required by species
	required_gas_min   int                   // Minimum needed percentage
	required_gas_max   int                   // Maximum allowed percentage
	neutral_gas        [6]gas_e              // Gases neutral to species
	poison_gas         [6]gas_e              // Gases poisonous to species
	auto_orders        bool                  // AUTO command was issued
	tech_level         [6]int                // Actual tech levels
	init_tech_level    [6]int                // Tech levels at start of turn
	tech_knowledge     [6]int                // Unapplied tech level knowledge
	num_namplas        int                   // Number of named planets, including home planet and colonies
	num_ships          int                   // Number of ships
	tech_eps           [6]int                // Experience points for tech levels
	hp_original_base   int                   // If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing
	econ_units         int                   // Number of economic units
	fleet_cost         int                   // Total fleet maintenance cost
	fleet_percent_cost int                   // Fleet maintenance cost as a percentage times one hundred
	contact            map[species_id_t]bool // A bit is set if corresponding species has been met
	ally               map[species_id_t]bool // A bit is set if corresponding species is considered an ally
	enemy              map[species_id_t]bool // A bit is set if corresponding species is considered an enemy
	home               struct {
		star   *star_data_t   // pointer to the star containing the planet containing the colony
		planet *planet_data_t // pointer to the planet containing the colony
		nampla *nampla_data_t // pointer to the nampla defining the colony
	}
}

type species_id_t int

type star_data_t struct {
	id           star_id_t    // unique identifier for this system
	index        int          // index of this system in star_base array
	x, y, z      int          // Coordinates
	type_        star_type_e  // Dwarf, degenerate, main sequence or giant
	color        star_color_e // Star color. Blue, blue-white, etc
	size         int          // Star size, from 0 through 9 inclusive
	num_planets  int          // Number of usable planets in star system
	home_system  bool         // TRUE if this is a good potential home system
	worm_here    bool         // TRUE if wormhole entry/exit
	worm_x       int          // Coordinates of wormhole's exit
	worm_y       int
	worm_z       int
	wormholeExit *star_data_t
	planet_index int                   // Index (starting at zero) into the file "planets.dat" of the first planet in the star system
	message      int                   // Message associated with this star system, if any
	visited_by   map[species_id_t]bool // A bit is set if corresponding species has been here
	planets      [10]*planet_data_t    // planets in this star system
}

type star_id_t int

type trans_data_t struct {
	type_       interspecies_transaction_e // Transaction type
	donor       species_id_t
	recipient   species_id_t
	value       int // Value of transaction
	x, y, z, pn int // Location associated with transaction
	number1     int // Other items associated with transaction
	name1       string
	number2     int
	name2       string
	number3     int
	name3       string
}

type wormhole_data_t struct {
	from_star_x int
	from_star_y int
	from_star_z int
	to_star_x   int
	to_star_y   int
	to_star_z   int
}
