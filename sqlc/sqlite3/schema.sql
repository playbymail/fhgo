--  Copyright (c) 2024 Michael D Henderson. All rights reserved.

-----------------------------------------------------------------------
-- drop tables to clear out old data and ready for initialization

-- foreign keys must be disabled to drop tables with foreign keys
PRAGMA foreign_keys = OFF;

DROP TABLE IF EXISTS star_color_e;

-- foreign keys must be enabled with every database connection
PRAGMA foreign_keys = ON;

-----------------------------------------------------------------------
-- lookup tables

-- planet_status_e might actually be planet_status_e.
CREATE TABLE planet_status
(
    code        TEXT PRIMARY KEY UNIQUE,
    value       TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
);

INSERT INTO planet_status (code, value, description)
VALUES ('0', 'NOT_SPECIAL', 'not special');
INSERT INTO planet_status (code, value, description)
VALUES ('1', 'IDEAL_HOME_PLANET', 'ideal home planet');
INSERT INTO planet_status (code, value, description)
VALUES ('2', 'IDEAL_COLONY_PLANET', 'ideal colony planet');
INSERT INTO planet_status (code, value, description)
VALUES ('3', 'RADIOACTIVE_HELLHOLE', 'radioactive hellhole');

-- star_color_e stores star_color_e values.
CREATE TABLE star_color_e
(
    code        TEXT PRIMARY KEY UNIQUE,
    value       TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
);

INSERT INTO star_color_e (code, value, description)
VALUES ('0', 'UNKNOWN_STAR_COLOR', 'unknown');
INSERT INTO star_color_e (code, value, description)
VALUES ('1', 'BLUE', 'blue');
INSERT INTO star_color_e (code, value, description)
VALUES ('2', 'BLUE_WHITE', 'blue white');
INSERT INTO star_color_e (code, value, description)
VALUES ('3', 'WHITE', 'white');
INSERT INTO star_color_e (code, value, description)
VALUES ('4', 'YELLOW_WHITE', 'yellow white');
INSERT INTO star_color_e (code, value, description)
VALUES ('5', 'YELLOW', 'yellow');
INSERT INTO star_color_e (code, value, description)
VALUES ('6', 'ORANGE', 'orange');
INSERT INTO star_color_e (code, value, description)
VALUES ('7', 'RED', 'red');

-- star_type_e stores star_type_e values.
CREATE TABLE star_type_e
(
    code        TEXT PRIMARY KEY UNIQUE,
    value       TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT ''
);

INSERT INTO star_type_e (code, value, description)
VALUES ('?', 'UNKNOWN_STAR_TYPE',
        'Both "unknown" and "main" were blanks in the original game engine!');
INSERT INTO star_type_e (code, value, description)
VALUES ('d', 'DWARF',
        'Smaller, cooler stars, typically red dwarfs, known for long, stable lifespans.');
INSERT INTO star_type_e (code, value, description)
VALUES ('D', 'DEGENERATE',
        'Often meaning a white dwarf, neutron star, or possibly even a black hole. These are remnants of stars that have expended their nuclear fuel and undergone gravitational collapse.');
INSERT INTO star_type_e (code, value, description)
VALUES (' ', 'MAIN_SEQUENCE sequence',
        'Star in the prime of its life, burning hydrogen in its core. Main sequence stars are the most common type, including stars like our Sun.');
INSERT INTO star_type_e (code, value, description)
VALUES ('g', 'GIANT',
        'Significantly larger and more luminous than main sequence stars. Giants often have expanded outer layers and are in later stages of stellar evolution, like red giants.');


-----------------------------------------------------------------------
-- core tables

-- game_state has a constraint to ensure only one row.
-- It ensures that the database has the state for a single game.
-- You must create separate databases to play multiple games.
CREATE TABLE game_state
(
    id INTEGER PRIMARY KEY CHECK (id = 1) -- Ensures single row
);

-- galaxy_data stores galaxy_data_t.
CREATE TABLE galaxy_data
(
    num_species INTEGER NOT NULL,
    radius      INTEGER NOT NULL,
    turn_number INTEGER NOT NULL, -- current turn number, will be 0 during game setup
    prng_seed   INTEGER NOT NULL
);

-- message_data stores message data.
CREATE TABLE message_data
(
    id      INTEGER PRIMARY KEY,
    message TEXT NOT NULL
);

-- need to understand where nampla.status is used.
-- item_quantity moved to nampla_inventory table.
CREATE TABLE nampla_data
(
    id             INTEGER PRIMARY KEY,        -- unique identifier for this system
    planet_id      INTEGER NOT NULL,           -- pointer to planet the colony is on
    name           TEXT    NOT NULL,           -- Name of planet
    AUs_needed     INTEGER,                    -- Incoming ship with only CUs on board
    AUs_to_install INTEGER,                    -- Colonial manufacturing units to be installed
    IUs_needed     INTEGER,                    -- Incoming ship with only CUs on board
    IUs_to_install INTEGER,                    -- Colonial mining units to be installed
    auto_AUs       INTEGER,                    -- Number of AUs to be automatically installed
    auto_IUs       INTEGER,                    -- Number of IUs to be automatically installed
    hidden         INTEGER NOT NULL DEFAULT 0, -- Colony is hidden
    hiding         INTEGER NOT NULL DEFAULT 0, -- HIDE order given
    ma_base        INTEGER,                    -- Manufacturing base times 10
    message        INTEGER,                    -- Message associated with this planet, if any
    mi_base        INTEGER,                    -- Mining base times 10
    pop_units      INTEGER,                    -- Number of available population units
    shipyards      INTEGER,                    -- Number of shipyards on planet
    siege_eff      INTEGER,                    -- Siege effectiveness - a percentage between 0 and 99
    special        INTEGER,                    -- Different for each application
    status         INTEGER,                    -- Status of planet
    use_on_ambush  INTEGER                     -- Amount to use on ambush
);

-- nampla_inventory stores inventory for a named planet (eg colony).
-- 	item_quantity  [MAX_ITEMS]int  -- Quantity of each item available
CREATE TABLE nampla_inventory
(
    nampla_id INTEGER NOT NULL,
    item_id   INTEGER NOT NULL,
    quantity  INTEGER NOT NULL,
    PRIMARY KEY (nampla_id, item_id)
);

-- planet_data stores planet_data_t.
-- gas data moved to planet_atmosphere_data table.
CREATE TABLE planet_data
(
    id                INTEGER PRIMARY KEY, -- unique identifier for this planet
    star_id           INTEGER NOT NULL,    -- pointer to the star the planet is orbiting
    pn                INTEGER NOT NULL,    -- orbital position of planet in the system
    diameter          INTEGER NOT NULL,    -- Diameter in thousands of kilometers
    econ_efficiency   INTEGER NOT NULL,    -- Economic efficiency. Always 100 for a home planet
    gravity           INTEGER NOT NULL,    -- Surface gravity. Multiple of Earth gravity times 100
    md_increase       INTEGER NOT NULL,    -- Increase in mining difficulty
    message           INTEGER NOT NULL,    -- Message associated with this planet, if any
    mining_difficulty INTEGER NOT NULL,    -- Mining difficulty times 100
    orbit             INTEGER NOT NULL,    -- orbit of planet in the system
    pressure_class    INTEGER NOT NULL,    -- Pressure class, 0-29
    special           INTEGER NOT NULL,    -- 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole
    temperature_class INTEGER NOT NULL     -- Temperature class, 1-30
);

CREATE TABLE planet_atmosphere_data
(
    planet_id INTEGER NOT NULL,
    gas_id    INTEGER NOT NULL, -- Gas in atmosphere
    percent   INTEGER NOT NULL  -- Percentage of gas in atmosphere
);

-- planet_inventory stores inventory for a planet.
-- 	item_quantity  [MAX_ITEMS]int  -- Quantity of each item available
CREATE TABLE planet_inventory
(
    planet_id INTEGER NOT NULL,
    item_id   INTEGER NOT NULL,
    quantity  INTEGER NOT NULL,
    PRIMARY KEY (planet_id, item_id)
);

-- ship_data stores ship_data_t.
-- item_quantity moved to nampla_inventory table.
CREATE TABLE ship_data
(
    id                   INTEGER PRIMARY KEY,        -- unique identifier for this system
    name                 TEXT    NOT NULL,           -- Name of ship
    x                    INTEGER NOT NULL,           -- Coordinates
    y                    INTEGER NOT NULL,           -- Coordinates
    z                    INTEGER NOT NULL,           -- Coordinates
    age                  INTEGER NOT NULL,           -- Ship age
    arrived_via_wormhole INTEGER NOT NULL DEFAULT 0, -- Ship arrived via wormhole in the PREVIOUS turn
    class                INTEGER NOT NULL,           -- Ship class
    dest_x               INTEGER NOT NULL,           -- Destination if ship was forced to jump from combat. Also used by TELESCOPE command
    dest_y               INTEGER NOT NULL,
    dest_z               INTEGER NOT NULL,
    just_jumped          INTEGER NOT NULL DEFAULT 0, -- Set if ship jumped this turn
    loading_point        INTEGER NOT NULL,           -- Nampla index for planet where ship was last loaded with CUs. Zero = none. Use 9999 for home planet
    pn                   INTEGER NOT NULL,           -- Current coordinates
    remaining_cost       INTEGER NOT NULL,           -- The cost needed to complete the ship if still under construction
    special              INTEGER,                    -- Different for each application
    status               INTEGER NOT NULL,           -- Current status of ship
    tonnage              INTEGER NOT NULL,           -- Ship tonnage divided by 10,000
    type_                INTEGER NOT NULL,           -- Ship type
    unloading_point      INTEGER NOT NULL            -- Nampla index for planet that ship should be given orders to jump to where it will unload. Zero = none. Use 9999 for home planet
);

-- ship_inventory stores inventory for a ship.
-- 	item_quantity  [MAX_ITEMS]int  -- Quantity of each item available
CREATE TABLE ship_inventory
(
    ship_id  INTEGER NOT NULL,
    item_id  INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    PRIMARY KEY (ship_id, item_id)
);

CREATE TABLE species_cfg
(
    email          TEXT    NOT NULL,
    name           TEXT    NOT NULL UNIQUE,
    govt_name      TEXT    NOT NULL,
    govt_type      TEXT    NOT NULL,
    homeworld_name TEXT    NOT NULL,
    bi             INTEGER NOT NULL,
    gv             INTEGER NOT NULL,
    ls             INTEGER NOT NULL,
    ml             INTEGER NOT NULL,
    PRIMARY KEY (email)
);

CREATE TABLE species_data
(
    id                 INTEGER PRIMARY KEY,
    auto_orders        INTEGER NOT NULL DEFAULT 0, -- AUTO command was issued
    econ_units         INTEGER NOT NULL,           -- Number of economic units
    fleet_cost         INTEGER NOT NULL,           -- Total fleet maintenance cost
    fleet_percent_cost INTEGER NOT NULL,           -- Fleet maintenance cost as a percentage times one hundred
    govt_name          TEXT    NOT NULL,           -- Name of government
    govt_type          TEXT    NOT NULL            -- Type of government
);

CREATE TABLE species_atmospheric_gases
(
    species_id     INTEGER NOT NULL,
    gas_id         INTEGER NOT NULL,
    poison         INTEGER NOT NULL DEFAULT 0,
    required       INTEGER NOT NULL DEFAULT 0,
    min_percentage INTEGER, -- minimum needed percentage, set only for required gases
    max_percentage INTEGER  -- maximum allowed percentage, set only for required gases
);

CREATE TABLE species_contacts
(
    species_id INTEGER NOT NULL,
    alien_id   INTEGER NOT NULL,
    contact    INTEGER NOT NULL DEFAULT 0,
    ally       INTEGER NOT NULL DEFAULT 0,
    enemy      INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (species_id, alien_id)
);

CREATE TABLE species_home_planet
(
    species_id       INTEGER NOT NULL,
    planet_id        INTEGER NOT NULL,
    hp_original_base INTEGER, -- If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing.
    PRIMARY KEY (species_id)
);

CREATE TABLE species_tech_levels
(
    species_id   INTEGER NOT NULL,
    bi           INTEGER NOT NULL DEFAULT 0, -- Biology tech level
    bi_exp       INTEGER NOT NULL DEFAULT 0, -- experience points for Biology tech level
    bi_unapplied INTEGER NOT NULL DEFAULT 0, -- un-applied Biology tech level
    gv           INTEGER NOT NULL DEFAULT 0, -- Gravitics tech level
    gv_exp       INTEGER NOT NULL DEFAULT 0, -- experience points for Gravitics tech level
    gv_unapplied INTEGER NOT NULL DEFAULT 0, -- un-applied Gravitics tech level
    ls           INTEGER NOT NULL DEFAULT 0, -- Life Support tech level
    ls_exp       INTEGER NOT NULL DEFAULT 0, -- experience points for Life Support tech level
    ls_unapplied INTEGER NOT NULL DEFAULT 0, -- un-applied Life Support tech level
    ma           INTEGER NOT NULL DEFAULT 0, -- Manufacturing tech level
    ma_exp       INTEGER NOT NULL DEFAULT 0, -- experience points for Manufacturing tech level
    ma_unapplied INTEGER NOT NULL DEFAULT 0, -- un-applied Manufacturing tech level
    mi           INTEGER NOT NULL DEFAULT 0, -- Mining tech level
    mi_exp       INTEGER NOT NULL DEFAULT 0, -- experience points for Mining tech level
    mi_unapplied INTEGER NOT NULL DEFAULT 0, -- un-applied Mining tech level
    ml           INTEGER NOT NULL DEFAULT 0, -- Military tech level
    ml_exp       INTEGER NOT NULL DEFAULT 0, -- experience points for Military tech level
    ml_unapplied INTEGER NOT NULL DEFAULT 0  -- un-applied Military tech level
);

-- star_data stores star_data_t.
-- note that wormhole data has been moved to the wormhole_data table.
-- note that visited_by has been moved to the star_visitied_by_table.
CREATE TABLE star_data
(
    id          INTEGER PRIMARY KEY, -- unique identifier for this system
    x           INTEGER NOT NULL,    -- Coordinates
    y           INTEGER NOT NULL,    -- Coordinates
    z           INTEGER NOT NULL,    -- Coordinates
    color       TEXT    NOT NULL,    -- Star color, e.g., Blue, blue-white
    home_system INTEGER NOT NULL DEFAULT 0,
    message_id  INTEGER,
    size        INTEGER NOT NULL,
    type_       TEXT    NOT NULL     -- Dwarf, degenerate, main sequence, or giant

--     -- check constraints
--     CONSTRAINT size_check CHECK (size BETWEEN 0 AND 9),         -- Star size, from 0 through 9 inclusive
--     CONSTRAINT home_system_check CHECK (home_system IN (0, 1)), -- TRUE if this is a good potential home system
--     -- foreign key constraints
--     CONSTRAINT color_fk FOREIGN KEY (color) REFERENCES star_color_e (code)
);

-- star_visited_by captures the species that have visited a star system.
CREATE TABLE star_visited_by
(
    star_id     INTEGER NOT NULL REFERENCES star_data (id),
    species_id  INTEGER NOT NULL REFERENCES species_data (id),
    turn_number INTEGER NOT NULL -- last turn visited by this species
);
--planets [10]*planet_data_t                 -- planets in this star system

CREATE TABLE wormhole_data
(
    from_star_x INTEGER NOT NULL,
    from_star_y INTEGER NOT NULL,
    from_star_z INTEGER NOT NULL,
    to_star_x   INTEGER NOT NULL,
    to_star_y   INTEGER NOT NULL,
    to_star_z   INTEGER NOT NULL
--     from_star_id INTEGER NOT NULL REFERENCES star_data (id),
--     to_star_id   INTEGER NOT NULL REFERENCES star_data (id),
--     CONSTRAINT in_system_check CHECK (from_star_id != to_star_id) -- wormhole must be between two different stars
);


-----------------------------------------------------------------------
-- lookup tables

-----------------------------------------------------------------------
-- add foreign key constraints


-----------------------------------------------------------------------
-- initialize lookup tables

