// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package fhgo

// Types of actions
type combat_action_e int

const (
	DEFENSE_IN_PLACE combat_action_e = iota
	DEEP_SPACE_DEFENSE
	PLANET_DEFENSE
	DEEP_SPACE_FIGHT
	PLANET_ATTACK
	PLANET_BOMBARDMENT
	GERM_WARFARE
	SIEGE
)

// Special types

const NON_COMBATANT = 1

// Types of combatants
type combatant_type_e int

const (
	NONCOMBATANT combatant_type_e = iota
	SHIP
	NAMPLA
	GENOCIDE_NAMPLA
	BESIEGED_NAMPLA
)

// Command codes
type command_code_e int

const (
	UNDEFINED command_code_e = iota
	ALLY
	AMBUSH
	ATTACK
	AUTO
	BASE
	BATTLE
	BUILD
	CONTINUE
	DEEP
	DESTROY
	DEVELOP
	DISBAND
	END
	ENEMY
	ENGAGE
	ESTIMATE
	HAVEN
	HIDE
	HIJACK
	IBUILD
	ICONTINUE
	INSTALL
	INTERCEPT
	JUMP
	LAND
	MESSAGE
	MOVE
	NAME
	NEUTRAL
	ORBIT
	PJUMP
	PRODUCTION
	RECYCLE
	REPAIR
	RESEARCH
	SCAN
	SEND
	SHIPYARD
	START
	SUMMARY
	SURRENDER
	TARGET
	TEACH
	TECH
	TELESCOPE
	TERRAFORM
	TRANSFER
	UNLOAD
	UPGRADE
	VISITED
	WITHDRAW
	WORMHOLE
	ZZZ
)

const NUM_COMMANDS = ZZZ + 1

// Gases in planetary atmospheres
type gas_e int

const (
	GAS_NONE gas_e = iota
	H2             /* Hydrogen */
	CH4            /* Methane */
	HE             /* Helium */
	NH3            /* Ammonia */
	N2             /* Nitrogen */
	CO2            /* Carbon Dioxide */
	O2             /* Oxygen */
	HCL            /* Hydrogen Chloride */
	CL2            /* Chlorine */
	F2             /* Fluorine */
	H2O            /* Steam */
	SO2            /* Sulfur Dioxide */
	H2S            /* Hydrogen Sulfide */
)

// Item IDs
type item_e int

const (
	RM  item_e = iota /* Raw Material Units. */
	PD                /* Planetary Defense Units. */
	SU                /* Starbase Units. */
	DR                /* Damage Repair Units. */
	CU                /* Colonist Units. */
	IU                /* Colonial Mining Units. */
	AU                /* Colonial Manufacturing Units. */
	FS                /* Fail-Safe Jump Units. */
	JP                /* Jump Portal Units. */
	FM                /* Forced Misjump Units. */
	FJ                /* Forced Jump Units. */
	GT                /* Gravitic Telescope Units. */
	FD                /* Field Distortion Units. */
	TP                /* Terraforming Plants. */
	GW                /* Germ Warfare Bombs. */
	SG1               /* Mark-1 Auxiliary Shield Generators. */
	SG2               /* Mark-2. */
	SG3               /* Mark-3. */
	SG4               /* Mark-4. */
	SG5               /* Mark-5. */
	SG6               /* Mark-6. */
	SG7               /* Mark-7. */
	SG8               /* Mark-8. */
	SG9               /* Mark-9. */
	GU1               /* Mark-1 Auxiliary Gun Units. */
	GU2               /* Mark-2. */
	GU3               /* Mark-3. */
	GU4               /* Mark-4. */
	GU5               /* Mark-5. */
	GU6               /* Mark-6. */
	GU7               /* Mark-7. */
	GU8               /* Mark-8. */
	GU9               /* Mark-9. */
	X1                /* Unassigned. */
	X2                /* Unassigned. */
	X3                /* Unassigned. */
	X4                /* Unassigned. */
	X5                /* Unassigned. */
)

// Interspecies transactions
type interspecies_transaction_e int

const (
	INTERSPECIES_TRANSACTION_TYPE_UNKNOWN interspecies_transaction_e = iota
	EU_TRANSFER
	MESSAGE_TO_SPECIES
	BESIEGE_PLANET
	SIEGE_EU_TRANSFER
	TECH_TRANSFER
	DETECTION_DURING_SIEGE
	SHIP_MISHAP
	ASSIMILATION
	INTERSPECIES_CONSTRUCTION
	TELESCOPE_DETECTION
	ALIEN_JUMP_PORTAL_USAGE
	KNOWLEDGE_TRANSFER
	LANDING_REQUEST
	LOOTING_EU_TRANSFER
	ALLIES_ORDER
)

// Status codes for named planets. These are logically ORed together.

const HOME_PLANET = 1 << 0
const COLONY = 1 << 1
const POPULATED = 1 << 3
const MINING_COLONY = 1 << 4
const RESORT_COLONY = 1 << 5
const DISBANDED_COLONY = 1 << 6

// Constants needed for parsing
type parser_token_e int

const (
	UNKNOWN parser_token_e = iota
	TECH_ID
	ITEM_CLASS
	SHIP_CLASS
	PLANET_ID
	SPECIES_ID
)

// planet_status_e might actually be planet_special_e
type planet_status_e int

type planet_special_e int

const (
	NOT_SPECIAL planet_special_e = iota
	IDEAL_HOME_PLANET
	IDEAL_COLONY_PLANET
	RADIOACTIVE_HELLHOLE
)

// Ship types
type ship_e int

const (
	FTL ship_e = iota
	SUB_LIGHT
	STARBASE
)

// Ship classes
type ship_class_e int

const (
	PB ship_class_e = iota /* Picketboat. */
	CT                     /* Corvette. */
	ES                     /* Escort. */
	FF                     /* Frigate. (was FG, was 4) */
	DD                     /* Destroyer. (was 3) */
	CL                     /* Light Cruiser. */
	CS                     /* Strike Cruiser. */
	CA                     /* Heavy Cruiser. */
	CC                     /* Command Cruiser. */
	BC                     /* Battlecruiser. */
	BS                     /* Battleship. */
	DN                     /* Dreadnought. */
	SD                     /* Super Dreadnought. */
	BM                     /* Battlemoon. */
	BW                     /* Battleworld. */
	BR                     /* Battlestar. */
	BA                     /* Starbase. */
	TR                     /* Transport. */
)
const NUM_SHIP_CLASSES = TR + 1

// Ship status codes
type ship_status_e int

const (
	UNDER_CONSTRUCTION ship_status_e = iota
	ON_SURFACE
	IN_ORBIT
	IN_DEEP_SPACE
	JUMPED_IN_COMBAT
	FORCED_JUMP
)

// Types of special targets
type special_target_e int

const (
	TARGET_NORMAL special_target_e = iota
	TARGET_WARSHIPS
	TARGET_TRANSPORTS
	TARGET_STARBASES
	TARGET_PDS
)

// * Star Colors
type star_color_e int

const (
	UNKNOWN_STAR_COLOR star_color_e = iota
	BLUE
	BLUE_WHITE
	WHITE
	YELLOW_WHITE
	YELLOW
	ORANGE
	RED
)

// Star types
type star_type_e byte

const (
	// UNKNOWN_STAR_TYPE TODO: both "unknown" and "main" were ' ' in the original game engine!
	UNKNOWN_STAR_TYPE star_type_e = star_type_e('?')
	// DWARF refers to dwarf stars, which are smaller, cooler stars, typically red dwarfs, known for long, stable lifespans.
	DWARF = star_type_e('d')
	// DEGENERATE points to a "degenerate" star, often meaning a white dwarf, neutron star, or possibly even a black hole. These are remnants of stars that have expended their nuclear fuel and undergone gravitational collapse.
	DEGENERATE = star_type_e('D')
	// MAIN_SEQUENCE indicates a main sequence star, which is a star in the prime of its life, burning hydrogen in its core. Main sequence stars are the most common type, including stars like our Sun.
	MAIN_SEQUENCE = star_type_e(' ')
	// GIANT refers to giant stars, which are significantly larger and more luminous than main sequence stars. Giants often have expanded outer layers and are in later stages of stellar evolution, like red giants.
	GIANT = star_type_e('g')
)

// Tech level ids
type tech_level_e int

const (
	MI tech_level_e = iota /* Mining tech level. */
	MA                     /* Manufacturing tech level. */
	ML                     /* Military tech level. */
	GV                     /* Gravitics tech level. */
	LS                     /* Life Support tech level. */
	BI                     /* Biology tech level. */
)
