--  Copyright (c) 2024 Michael D Henderson. All rights reserved.

-- I'll start with the core game objects based on the structures in the GAME_OBJECTS_SOURCES files. Here's a proposed schema for the primary entities:

-- Game state table has a constraint to ensure only one row.
-- It ensures that the database has the state for a single game.
-- You must create separate databases to play multiple games.
CREATE TABLE game_state
(
    id            INTEGER PRIMARY KEY CHECK (id = 1), -- Ensures single row
    code          TEXT    NOT NULL,                   -- unique short identifier for the game
    name          TEXT    NOT NULL,                   -- long name of the game
    description   TEXT    NOT NULL,                   -- description of the game
    turn_number   INTEGER NOT NULL,                   -- current turn number, will be 0 during game setup
    galaxy_radius INTEGER NOT NULL,
    num_species   INTEGER NOT NULL,
    num_stars     INTEGER NOT NULL,
    num_planets   INTEGER NOT NULL,
    prng_seed     INTEGER NOT NULL
);

-- Species table - core player/race data
CREATE TABLE species (
                         id INTEGER PRIMARY KEY,
                         name TEXT NOT NULL,
                         govt_name TEXT NOT NULL,
                         govt_type TEXT NOT NULL,
                         home_system_id INTEGER,
                         tech_level_ml INTEGER NOT NULL,
                         tech_level_gv INTEGER NOT NULL,
                         tech_level_ls INTEGER NOT NULL,
                         tech_level_bi INTEGER NOT NULL
);

-- Stars/Systems table
CREATE TABLE stars (
                       id INTEGER PRIMARY KEY,
                       x INTEGER NOT NULL,
                       y INTEGER NOT NULL,
                       z INTEGER NOT NULL,
                       is_home_system BOOLEAN DEFAULT FALSE,
                       wormhole_destination_id INTEGER,
                       FOREIGN KEY(wormhole_destination_id) REFERENCES stars(id)
);

-- Planets table
CREATE TABLE planets (
                         id INTEGER PRIMARY KEY,
                         star_id INTEGER NOT NULL,
                         orbit_number INTEGER NOT NULL,
                         mining_difficulty INTEGER NOT NULL,
                         mineral_content INTEGER NOT NULL,
                         is_homeworld BOOLEAN DEFAULT FALSE,
                         FOREIGN KEY(star_id) REFERENCES stars(id)
);

-- Ships table
CREATE TABLE ships (
                       id INTEGER PRIMARY KEY,
                       species_id INTEGER NOT NULL,
                       name TEXT NOT NULL,
                       class TEXT NOT NULL,
                       current_star_id INTEGER,
                       status TEXT NOT NULL,
                       FOREIGN KEY(species_id) REFERENCES species(id),
                       FOREIGN KEY(current_star_id) REFERENCES stars(id)
);

-- Would you like to see the supporting tables next (like combat, production, resources)?
-- Here are the supporting tables that handle game mechanics and resources:

-- Production/Resources table
CREATE TABLE resources (
                           planet_id INTEGER NOT NULL,
                           species_id INTEGER NOT NULL,
                           economic_units INTEGER DEFAULT 0,
                           mining_units INTEGER DEFAULT 0,
                           manufacturing_units INTEGER DEFAULT 0,
                           auto_IUs INTEGER DEFAULT 0,
                           AUs_needed INTEGER DEFAULT 0,
                           PRIMARY KEY(planet_id, species_id),
                           FOREIGN KEY(planet_id) REFERENCES planets(id),
                           FOREIGN KEY(species_id) REFERENCES species(id)
);

-- Combat/Military table
CREATE TABLE military (
                          ship_id INTEGER NOT NULL,
                          combat_status TEXT,
                          strike_target_id INTEGER,
                          unloading_point INTEGER DEFAULT 0,
                          PRIMARY KEY(ship_id),
                          FOREIGN KEY(ship_id) REFERENCES ships(id),
                          FOREIGN KEY(strike_target_id) REFERENCES ships(id)
);

-- Would you like to see the tables for tracking turns and game state next?
-- Here are the tables for managing game state and turns:

-- Turn events table
CREATE TABLE turn_events (
                             id INTEGER PRIMARY KEY,
                             turn_number INTEGER NOT NULL,
                             event_type TEXT NOT NULL,
                             species_id INTEGER NOT NULL,
                             location_x INTEGER,
                             location_y INTEGER,
                             location_z INTEGER,
                             description TEXT NOT NULL,
                             timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
                             FOREIGN KEY(species_id) REFERENCES species(id)
);

-- Turn orders table
CREATE TABLE orders (
                        id INTEGER PRIMARY KEY,
                        turn_number INTEGER NOT NULL,
                        species_id INTEGER NOT NULL,
                        order_type TEXT NOT NULL,
                        order_sequence INTEGER NOT NULL,
                        parameters TEXT NOT NULL,
                        FOREIGN KEY(species_id) REFERENCES species(id)
);


-- Would you like to see indexes and constraints next?
-- Here are the key indexes and constraints to optimize queries and maintain data integrity:

-- Indexes for spatial queries
CREATE INDEX idx_stars_location ON stars(x, y, z);
CREATE INDEX idx_stars_home ON stars(is_home_system);

-- Indexes for turn processing
CREATE INDEX idx_orders_turn_seq ON orders(turn_number, species_id, order_sequence);
CREATE INDEX idx_events_turn ON turn_events(turn_number, species_id);

-- Indexes for common lookups
CREATE INDEX idx_planets_star ON planets(star_id);
CREATE INDEX idx_ships_species ON ships(species_id);
CREATE INDEX idx_ships_location ON ships(current_star_id);

-- Unique constraints
CREATE UNIQUE INDEX idx_species_name ON species(name);
CREATE UNIQUE INDEX idx_ships_name ON ships(species_id, name);


-- Would you like to see the views that would be useful for reporting next?
-- Here are useful views for reporting game state and statistics:

-- Species status view
CREATE VIEW v_species_status AS
SELECT
    s.name,
    COUNT(DISTINCT p.id) as num_planets,
    COUNT(DISTINCT sh.id) as num_ships,
    SUM(r.economic_units) as total_economic_units,
    SUM(r.mining_units) as total_mining_units
FROM species s
         LEFT JOIN resources r ON s.id = r.species_id
         LEFT JOIN planets p ON r.planet_id = p.id
         LEFT JOIN ships sh ON s.id = sh.species_id
GROUP BY s.id;

-- System activity view
CREATE VIEW v_system_activity AS
SELECT
    st.id as star_id,
    st.x, st.y, st.z,
    COUNT(DISTINCT p.id) as num_planets,
    COUNT(DISTINCT sh.id) as ships_present,
    COUNT(DISTINCT r.species_id) as species_present
FROM stars st
         LEFT JOIN planets p ON st.id = p.star_id
         LEFT JOIN ships sh ON st.id = sh.current_star_id
         LEFT JOIN resources r ON p.id = r.planet_id
GROUP BY st.id;


-- Would you like to see additional views for combat or production reporting?
-- Here are views for combat and production reporting:

-- Combat status view
CREATE VIEW v_combat_status AS
SELECT
    s.name as species_name,
    sh.name as ship_name,
    sh.class,
    st.x, st.y, st.z,
    m.combat_status,
    m.strike_target_id
FROM ships sh
         JOIN species s ON sh.species_id = s.id
         JOIN stars st ON sh.current_star_id = st.id
         JOIN military m ON sh.id = m.ship_id
WHERE m.combat_status IS NOT NULL;

-- Production overview view
CREATE VIEW v_production_status AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    p.orbit_number,
    r.economic_units,
    r.mining_units,
    r.manufacturing_units,
    r.auto_IUs,
    r.AUs_needed
FROM resources r
         JOIN species s ON r.species_id = s.id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id
ORDER BY s.name, st.x, st.y, st.z, p.orbit_number;


-- Would you like to see additional views for turn processing or game mechanics?
-- Here are views for turn processing and game mechanics:

-- Turn orders summary view
CREATE VIEW v_turn_orders_summary AS
SELECT
    o.turn_number,
    s.name as species_name,
    o.order_type,
    COUNT(*) as order_count
FROM orders o
         JOIN species s ON o.species_id = s.id
GROUP BY o.turn_number, s.id, o.order_type;

-- Movement and jump status view
CREATE VIEW v_movement_status AS
SELECT
    s.name as species_name,
    sh.name as ship_name,
    sh.class,
    src.x as from_x, src.y as from_y, src.z as from_z,
    dest.x as to_x, dest.y as to_y, dest.z as to_z,
    m.unloading_point
FROM ships sh
         JOIN species s ON sh.species_id = s.id
         JOIN stars src ON sh.current_star_id = src.id
         JOIN military m ON sh.id = m.ship_id
         LEFT JOIN stars dest ON m.unloading_point = dest.id
WHERE m.unloading_point > 0;

-- Would you like to see views for tech levels and development tracking?
-- Here are views for tracking technology and development:

-- Technology levels view
CREATE VIEW v_tech_levels AS
SELECT
    s.name as species_name,
    s.tech_level_ml as mining_tech,
    s.tech_level_gv as gravitics_tech,
    s.tech_level_ls as life_support_tech,
    s.tech_level_bi as biology_tech,
    (s.tech_level_ml + s.tech_level_gv + s.tech_level_ls + s.tech_level_bi) as total_tech
FROM species s;

-- Development tracking view
CREATE VIEW v_development_status AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    p.orbit_number,
    r.economic_units,
    r.mining_units + r.manufacturing_units as industrial_units,
    r.auto_IUs as automated_industry,
    CASE
        WHEN p.is_homeworld THEN 'Homeworld'
        ELSE 'Colony'
        END as planet_type
FROM resources r
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id
         JOIN species s ON r.species_id = s.id
ORDER BY s.name, planet_type DESC;


-- Would you like to see views for spatial analysis and exploration tracking?
-- Here are views for spatial analysis and exploration tracking:

-- Spatial density view
CREATE VIEW v_spatial_density AS
SELECT
    CAST(x/10 AS INTEGER) as sector_x,
    CAST(y/10 AS INTEGER) as sector_y,
    CAST(z/10 AS INTEGER) as sector_z,
    COUNT(id) as star_count,
    COUNT(DISTINCT wormhole_destination_id) as wormhole_count
FROM stars
GROUP BY sector_x, sector_y, sector_z;

-- Species exploration view
CREATE VIEW v_exploration_status AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    COUNT(DISTINCT p.id) as planets_in_system,
    COUNT(DISTINCT sh.id) as ships_in_system,
    MAX(te.timestamp) as last_visited
FROM species s
         JOIN ships sh ON s.id = sh.species_id
         JOIN stars st ON sh.current_star_id = st.id
         JOIN planets p ON st.id = p.star_id
         LEFT JOIN turn_events te ON s.id = te.species_id
    AND st.x = te.location_x
    AND st.y = te.location_y
    AND st.z = te.location_z
GROUP BY s.id, st.id;


-- Would you like to see views for economic analysis and resource distribution?
-- Here are views for economic analysis and resource distribution:

-- Economic distribution view
CREATE VIEW v_economic_distribution AS
SELECT
    s.name as species_name,
    COUNT(p.id) as total_planets,
    SUM(r.economic_units) as total_economy,
    SUM(r.mining_units) as total_mining,
    SUM(r.manufacturing_units) as total_manufacturing,
    ROUND(AVG(r.economic_units), 2) as avg_economy_per_planet,
    MAX(r.economic_units) as largest_economy
FROM species s
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
GROUP BY s.id;

-- Resource concentration view
CREATE VIEW v_resource_concentration AS
SELECT
    st.x, st.y, st.z,
    COUNT(p.id) as planets_in_system,
    SUM(p.mineral_content) as total_minerals,
    MIN(p.mining_difficulty) as best_mining_difficulty,
    COUNT(DISTINCT r.species_id) as species_present
FROM stars st
         JOIN planets p ON st.id = p.star_id
         LEFT JOIN resources r ON p.id = r.planet_id
GROUP BY st.id;


-- Would you like to see views for strategic analysis and relative power metrics?
-- Here are views for strategic analysis and relative power metrics:

-- Strategic power view
CREATE VIEW v_strategic_power AS
SELECT
    s.name as species_name,
    COUNT(DISTINCT sh.id) as total_ships,
    COUNT(DISTINCT p.id) as controlled_planets,
    SUM(r.economic_units + r.mining_units + r.manufacturing_units) as industrial_power,
    s.tech_level_ml + s.tech_level_gv + s.tech_level_ls + s.tech_level_bi as tech_power,
    COUNT(DISTINCT st.id) as systems_present
FROM species s
         LEFT JOIN ships sh ON s.id = sh.species_id
         LEFT JOIN resources r ON s.id = r.species_id
         LEFT JOIN planets p ON r.planet_id = p.id
         LEFT JOIN stars st ON p.star_id = st.id
GROUP BY s.id;

-- Relative position analysis
CREATE VIEW v_relative_positions AS
SELECT
    s1.name as species_one,
    s2.name as species_two,
    MIN(
            SQRT(
                    POWER(st1.x - st2.x, 2) +
                    POWER(st1.y - st2.y, 2) +
                    POWER(st1.z - st2.z, 2)
            )
    ) as minimum_distance,
    COUNT(DISTINCT st1.id) as shared_systems
FROM species s1
         CROSS JOIN species s2
         JOIN ships sh1 ON s1.id = sh1.species_id
         JOIN ships sh2 ON s2.id = sh2.species_id
         JOIN stars st1 ON sh1.current_star_id = st1.id
         JOIN stars st2 ON sh2.current_star_id = st2.id
WHERE s1.id < s2.id
GROUP BY s1.id, s2.id;


-- Would you like to see views for tracking game events and historical trends?
-- Here are views for tracking game events and historical trends:

-- Historical events analysis
CREATE VIEW v_event_history AS
SELECT
    turn_number,
    species_id,
    event_type,
    COUNT(*) as event_count,
    MIN(timestamp) as first_occurrence,
    MAX(timestamp) as last_occurrence
FROM turn_events
GROUP BY turn_number, species_id, event_type;

-- Species growth trends
CREATE VIEW v_growth_trends AS
SELECT
    te.turn_number,
    s.name as species_name,
    COUNT(DISTINCT p.id) as planets_controlled,
    COUNT(DISTINCT sh.id) as total_ships,
    SUM(r.economic_units) as total_economy
FROM turn_events te
         JOIN species s ON te.species_id = s.id
         LEFT JOIN planets p ON te.location_x = p.id
         LEFT JOIN ships sh ON te.species_id = sh.species_id
         LEFT JOIN resources r ON te.species_id = r.species_id
GROUP BY te.turn_number, s.id
ORDER BY te.turn_number, s.name;


-- Would you like to see views for analyzing game balance and competition metrics?
-- Here are views for analyzing game balance and competition metrics:

-- Competition intensity view
CREATE VIEW v_competition_metrics AS
SELECT
    st.x, st.y, st.z,
    COUNT(DISTINCT s.id) as competing_species,
    COUNT(DISTINCT sh.id) as total_ships,
    SUM(r.economic_units) as total_economic_power,
    COUNT(DISTINCT p.id) as contested_planets
FROM stars st
         JOIN planets p ON st.id = p.star_id
         JOIN resources r ON p.id = r.planet_id
         JOIN species s ON r.species_id = s.id
         JOIN ships sh ON s.id = sh.species_id
WHERE sh.current_star_id = st.id
GROUP BY st.id
HAVING competing_species > 1;

-- Species balance metrics
CREATE VIEW v_balance_metrics AS
SELECT
    s.name,
    COUNT(DISTINCT p.id) as controlled_planets,
    SUM(r.economic_units) as economic_power,
    SUM(r.mining_units) as mining_power,
    COUNT(DISTINCT sh.id) as fleet_size,
    ROUND(
            (s.tech_level_ml + s.tech_level_gv + s.tech_level_ls + s.tech_level_bi) /
            (SELECT AVG(tech_level_ml + tech_level_gv + tech_level_ls + tech_level_bi) FROM species)
        , 2) as relative_tech_level
FROM species s
         LEFT JOIN resources r ON s.id = r.species_id
         LEFT JOIN planets p ON r.planet_id = p.id
         LEFT JOIN ships sh ON s.id = sh.species_id
GROUP BY s.id;

-- Would you like to see views for analyzing spatial expansion patterns?
-- Here are views for analyzing spatial expansion patterns:

-- Expansion patterns view
CREATE VIEW v_expansion_patterns AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    COUNT(DISTINCT p.id) as planets_in_system,
    SQRT(POWER(st.x - home.x, 2) + POWER(st.y - home.y, 2) + POWER(st.z - home.z, 2)) as distance_from_home
FROM species s
         JOIN stars home ON s.home_system_id = home.id
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id
GROUP BY s.id, st.id;

-- Territory analysis view
CREATE VIEW v_territory_analysis AS
SELECT
    s.name as species_name,
    COUNT(DISTINCT st.id) as systems_controlled,
    MAX(SQRT(POWER(st.x - home.x, 2) + POWER(st.y - home.y, 2) + POWER(st.z - home.z, 2))) as max_expansion_distance,
    AVG(SQRT(POWER(st.x - home.x, 2) + POWER(st.y - home.y, 2) + POWER(st.z - home.z, 2))) as avg_expansion_distance
FROM species s
         JOIN stars home ON s.home_system_id = home.id
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id
GROUP BY s.id;

-- Would you like to see views for analyzing wormhole usage and strategic corridors?
-- Here are views for analyzing wormhole usage and strategic corridors:

-- Wormhole network analysis
CREATE VIEW v_wormhole_network AS
SELECT
    src.x as source_x, src.y as source_y, src.z as source_z,
    dest.x as dest_x, dest.y as dest_y, dest.z as dest_z,
    COUNT(DISTINCT sh.id) as ships_using_route,
    COUNT(DISTINCT sh.species_id) as species_using_route
FROM stars src
         JOIN stars dest ON src.wormhole_destination_id = dest.id
         LEFT JOIN ships sh ON sh.current_star_id = src.id
GROUP BY src.id, dest.id;

-- Strategic corridors view
CREATE VIEW v_strategic_corridors AS
SELECT
    st.x, st.y, st.z,
    COUNT(DISTINCT sh.id) as ships_present,
    COUNT(DISTINCT sh.species_id) as species_present,
    CASE
        WHEN st.wormhole_destination_id IS NOT NULL THEN 'Wormhole'
        ELSE 'Regular'
        END as system_type,
    COUNT(DISTINCT p.id) as planets_in_system
FROM stars st
         LEFT JOIN ships sh ON st.id = sh.current_star_id
         LEFT JOIN planets p ON st.id = p.star_id
GROUP BY st.id;


-- Would you like to see views for analyzing movement patterns and fleet distributions?
-- Here are views for analyzing movement patterns and fleet distributions:

-- Fleet movement patterns
CREATE VIEW v_fleet_movements AS
SELECT
    s.name as species_name,
    sh.class as ship_type,
    st.x, st.y, st.z,
    COUNT(*) as fleet_size,
    m.unloading_point,
    CASE
        WHEN m.combat_status IS NOT NULL THEN 'Combat'
        WHEN m.unloading_point > 0 THEN 'Transport'
        ELSE 'Stationed'
        END as mission_type
FROM ships sh
         JOIN species s ON sh.species_id = s.id
         JOIN stars st ON sh.current_star_id = st.id
         JOIN military m ON sh.id = m.ship_id
GROUP BY s.id, sh.class, st.id, mission_type;

-- Fleet distribution analysis
CREATE VIEW v_fleet_distribution AS
SELECT
    s.name as species_name,
    COUNT(DISTINCT st.id) as systems_with_ships,
    MAX(system_fleet_size) as largest_fleet_size,
    AVG(system_fleet_size) as avg_fleet_size
FROM species s
         JOIN (
    SELECT species_id, current_star_id, COUNT(*) as system_fleet_size
    FROM ships
    GROUP BY species_id, current_star_id
) fleet ON s.id = fleet.species_id
         JOIN stars st ON fleet.current_star_id = st.id
GROUP BY s.id;


-- Would you like to see views for analyzing combat effectiveness and military positioning?
-- Here are views for analyzing combat effectiveness and military positioning:

-- Combat effectiveness analysis
CREATE VIEW v_combat_effectiveness AS
SELECT
    s.name as species_name,
    sh.class as ship_type,
    COUNT(*) as total_ships,
    COUNT(CASE WHEN m.combat_status IS NOT NULL THEN 1 END) as ships_in_combat,
    COUNT(CASE WHEN m.strike_target_id IS NOT NULL THEN 1 END) as strike_missions,
    COUNT(DISTINCT st.id) as combat_zones
FROM ships sh
         JOIN species s ON sh.species_id = s.id
         JOIN stars st ON sh.current_star_id = st.id
         JOIN military m ON sh.id = m.ship_id
GROUP BY s.id, sh.class;

-- Military positioning view
CREATE VIEW v_military_positions AS
SELECT
    st.x, st.y, st.z,
    s.name as species_name,
    COUNT(sh.id) as fleet_strength,
    COUNT(DISTINCT sh.class) as ship_types,
    MIN(SQRT(POWER(st.x - home.x, 2) + POWER(st.y - home.y, 2) + POWER(st.z - home.z, 2))) as distance_from_home
FROM ships sh
         JOIN species s ON sh.species_id = s.id
         JOIN stars st ON sh.current_star_id = st.id
         JOIN stars home ON s.home_system_id = home.id
WHERE sh.status != 'Destroyed'
GROUP BY st.id, s.id;


-- Would you like to see views for analyzing defensive positions and strategic vulnerabilities?
-- Here are views for analyzing defensive positions and strategic vulnerabilities:

-- Defensive positions analysis
CREATE VIEW v_defensive_positions AS
SELECT
    st.x, st.y, st.z,
    s.name as species_name,
    COUNT(DISTINCT p.id) as defended_planets,
    SUM(r.economic_units) as economic_value,
    COUNT(DISTINCT sh.id) as defending_ships,
    COUNT(DISTINCT n.id) as neighboring_systems
FROM stars st
         JOIN planets p ON st.id = p.star_id
         JOIN resources r ON p.id = r.planet_id
         JOIN species s ON r.species_id = s.id
         LEFT JOIN ships sh ON sh.current_star_id = st.id AND sh.species_id = s.id
         LEFT JOIN stars n ON (
    POWER(st.x - n.x, 2) + POWER(st.y - n.y, 2) + POWER(st.z - n.z, 2) <= 100
    )
GROUP BY st.id, s.id;

-- Strategic vulnerability index
CREATE VIEW v_vulnerability_index AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    r.economic_units + r.mining_units + r.manufacturing_units as industrial_value,
    COUNT(DISTINCT sh.id) as defensive_forces,
    COUNT(DISTINCT enemy.id) as potential_threats
FROM species s
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id
         LEFT JOIN ships sh ON sh.current_star_id = st.id AND sh.species_id = s.id
         LEFT JOIN ships enemy ON enemy.current_star_id = st.id AND enemy.species_id != s.id
GROUP BY st.id, s.id;


-- Would you like to see views for analyzing strategic chokepoints and system control?
-- Here are views for analyzing strategic chokepoints and system control:

-- Strategic chokepoints analysis
CREATE VIEW v_chokepoints AS
SELECT
    st.x, st.y, st.z,
    COUNT(DISTINCT p.id) as planets_in_system,
    COUNT(DISTINCT s.id) as competing_species,
    COUNT(DISTINCT sh.id) as total_ships,
    CASE
        WHEN st.wormhole_destination_id IS NOT NULL THEN 'Wormhole'
        ELSE 'Regular'
        END as system_type,
    COUNT(DISTINCT n.id) as connected_systems
FROM stars st
         LEFT JOIN planets p ON st.id = p.star_id
         LEFT JOIN resources r ON p.id = r.planet_id
         LEFT JOIN species s ON r.species_id = s.id
         LEFT JOIN ships sh ON sh.current_star_id = st.id
         LEFT JOIN stars n ON (
    POWER(st.x - n.x, 2) + POWER(st.y - n.y, 2) + POWER(st.z - n.z, 2) <= 100
    )
GROUP BY st.id;

-- System control analysis
CREATE VIEW v_system_control AS
SELECT
    st.x, st.y, st.z,
    s.name as controlling_species,
    COUNT(DISTINCT p.id) as controlled_planets,
    SUM(r.economic_units) as economic_power,
    COUNT(DISTINCT sh.id) as military_presence,
    COUNT(DISTINCT other_sh.id) as foreign_ships
FROM stars st
         JOIN planets p ON st.id = p.star_id
         JOIN resources r ON p.id = r.planet_id
         JOIN species s ON r.species_id = s.id
         LEFT JOIN ships sh ON sh.current_star_id = st.id AND sh.species_id = s.id
         LEFT JOIN ships other_sh ON other_sh.current_star_id = st.id AND other_sh.species_id != s.id
GROUP BY st.id, s.id;


-- Would you like to see views for analyzing expansion opportunities and contested zones?
-- Here are views for analyzing expansion opportunities and contested zones:

-- Expansion opportunities analysis
CREATE VIEW v_expansion_opportunities AS
SELECT
    st.x, st.y, st.z,
    COUNT(p.id) as available_planets,
    MIN(p.mining_difficulty) as best_mining_potential,
    MAX(p.mineral_content) as max_minerals,
    COUNT(DISTINCT n.id) as nearby_systems,
    MIN(
            SELECT SQRT(POWER(st.x - home.x, 2) + POWER(st.y - home.y, 2) + POWER(st.z - home.z, 2))
            FROM species s
            JOIN stars home ON s.home_system_id = home.id
    ) as distance_to_nearest_empire
FROM stars st
         LEFT JOIN planets p ON st.id = p.star_id
         LEFT JOIN resources r ON p.id = r.planet_id
         LEFT JOIN stars n ON (
    POWER(st.x - n.x, 2) + POWER(st.y - n.y, 2) + POWER(st.z - n.z, 2) <= 100
    )
WHERE r.species_id IS NULL
GROUP BY st.id;

-- Contested zones analysis
CREATE VIEW v_contested_zones AS
SELECT
    st.x, st.y, st.z,
    COUNT(DISTINCT s.id) as active_species,
    GROUP_CONCAT(DISTINCT s.name) as present_species,
    COUNT(DISTINCT p.id) as contested_planets,
    SUM(r.economic_units) as total_economic_value,
    COUNT(DISTINCT sh.id) as military_units
FROM stars st
         JOIN planets p ON st.id = p.star_id
         JOIN resources r ON p.id = r.planet_id
         JOIN species s ON r.species_id = s.id
         LEFT JOIN ships sh ON sh.current_star_id = st.id
GROUP BY st.id
HAVING active_species > 1;


-- Would you like to see views for analyzing technological advantages and resource efficiency?
-- Here are views for analyzing technological advantages and resource efficiency:

-- Technology advantage analysis
CREATE VIEW v_tech_advantages AS
SELECT
    s.name as species_name,
    s.tech_level_ml as mining_tech,
    s.tech_level_gv as gravitics_tech,
    s.tech_level_ls as life_support_tech,
    s.tech_level_bi as biology_tech,
    ROUND(s.tech_level_ml / avg_ml.avg_level, 2) as mining_advantage,
    ROUND(s.tech_level_gv / avg_gv.avg_level, 2) as gravitics_advantage,
    ROUND(s.tech_level_ls / avg_ls.avg_level, 2) as life_support_advantage,
    ROUND(s.tech_level_bi / avg_bi.avg_level, 2) as biology_advantage
FROM species s
         CROSS JOIN (SELECT AVG(tech_level_ml) as avg_level FROM species) avg_ml
         CROSS JOIN (SELECT AVG(tech_level_gv) as avg_level FROM species) avg_gv
         CROSS JOIN (SELECT AVG(tech_level_ls) as avg_level FROM species) avg_ls
         CROSS JOIN (SELECT AVG(tech_level_bi) as avg_level FROM species) avg_bi;

-- Resource efficiency metrics
CREATE VIEW v_resource_efficiency AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    r.economic_units,
    r.mining_units,
    r.manufacturing_units,
    r.auto_IUs,
    ROUND(r.mining_units * s.tech_level_ml / p.mining_difficulty, 2) as mining_efficiency,
    ROUND(r.economic_units / NULLIF(r.mining_units + r.manufacturing_units, 0), 2) as economic_efficiency
FROM species s
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id;


-- Would you like to see views for analyzing production capacity and industrial development?
-- Here are views for analyzing production capacity and industrial development:

-- Production capacity analysis
CREATE VIEW v_production_capacity AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    SUM(r.manufacturing_units) as total_manufacturing,
    SUM(r.mining_units) as total_mining,
    SUM(r.auto_IUs) as automated_industry,
    COUNT(DISTINCT p.id) as production_centers,
    ROUND(SUM(r.manufacturing_units) / NULLIF(COUNT(p.id), 0), 2) as avg_manufacturing_per_planet
FROM species s
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id
GROUP BY s.id, st.id;

-- Industrial development trends
CREATE VIEW v_industrial_development AS
SELECT
    s.name as species_name,
    COUNT(p.id) as developed_planets,
    SUM(r.manufacturing_units + r.mining_units) as total_industrial_base,
    SUM(r.auto_IUs) as automated_capacity,
    SUM(r.AUs_needed) as automation_needs,
    ROUND(SUM(r.auto_IUs) / NULLIF(SUM(r.manufacturing_units + r.mining_units), 0) * 100, 2) as automation_percentage
FROM species s
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
GROUP BY s.id;


-- Would you like to see views for analyzing industrial efficiency and automation trends?
-- Here are views for analyzing industrial efficiency and automation trends:

-- Industrial efficiency analysis
CREATE VIEW v_industrial_efficiency AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    p.mining_difficulty,
    r.mining_units,
    r.manufacturing_units,
    r.auto_IUs,
    ROUND(r.mining_units * s.tech_level_ml / p.mining_difficulty, 2) as mining_yield,
    ROUND(r.manufacturing_units * s.tech_level_gv, 2) as manufacturing_output,
    ROUND(r.auto_IUs / NULLIF(r.manufacturing_units, 0) * 100, 2) as automation_ratio
FROM species s
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id;

-- Automation deployment patterns
CREATE VIEW v_automation_patterns AS
SELECT
    s.name as species_name,
    COUNT(CASE WHEN r.auto_IUs > 0 THEN 1 END) as automated_planets,
    SUM(r.auto_IUs) as total_automation,
    MAX(r.auto_IUs) as max_automation_level,
    AVG(CASE WHEN r.auto_IUs > 0 THEN r.auto_IUs END) as avg_automation_level,
    SUM(r.AUs_needed) as pending_automation_needs
FROM species s
         JOIN resources r ON s.id = r.species_id
GROUP BY s.id;


-- Would you like to see views for analyzing resource distribution and industrial specialization?
-- Here are views for analyzing resource distribution and industrial specialization:

-- Resource distribution analysis
CREATE VIEW v_resource_distribution AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    p.mineral_content,
    p.mining_difficulty,
    r.mining_units,
    r.manufacturing_units,
    r.economic_units,
    CASE
        WHEN r.mining_units > r.manufacturing_units THEN 'Mining'
        WHEN r.manufacturing_units > r.mining_units THEN 'Manufacturing'
        ELSE 'Balanced'
        END as specialization
FROM species s
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id;

-- Industrial specialization patterns
CREATE VIEW v_industrial_specialization AS
SELECT
    s.name as species_name,
    COUNT(CASE WHEN r.mining_units > r.manufacturing_units THEN 1 END) as mining_focused_planets,
    COUNT(CASE WHEN r.manufacturing_units > r.mining_units THEN 1 END) as manufacturing_focused_planets,
    COUNT(CASE WHEN r.mining_units = r.manufacturing_units THEN 1 END) as balanced_planets,
    ROUND(AVG(r.mining_units), 2) as avg_mining_units,
    ROUND(AVG(r.manufacturing_units), 2) as avg_manufacturing_units
FROM species s
         JOIN resources r ON s.id = r.species_id
GROUP BY s.id;


-- Would you like to see views for analyzing industrial clusters and production networks?
-- Here are views for analyzing industrial clusters and production networks:

-- Industrial clusters analysis
CREATE VIEW v_industrial_clusters AS
SELECT
    st.x, st.y, st.z,
    COUNT(DISTINCT p.id) as planets_in_cluster,
    SUM(r.manufacturing_units) as cluster_manufacturing,
    SUM(r.mining_units) as cluster_mining,
    COUNT(DISTINCT s.id) as species_present,
    GROUP_CONCAT(DISTINCT s.name) as participating_species
FROM stars st
         JOIN planets p ON st.id = p.star_id
         JOIN resources r ON p.id = r.planet_id
         JOIN species s ON r.species_id = s.id
GROUP BY st.id
HAVING cluster_manufacturing + cluster_mining > 100;

-- Production network efficiency
CREATE VIEW v_production_networks AS
SELECT
    s.name as species_name,
    COUNT(DISTINCT st.id) as production_systems,
    SUM(r.manufacturing_units) as total_manufacturing,
    SUM(r.mining_units) as total_mining,
    MAX(SQRT(POWER(st.x - home.x, 2) + POWER(st.y - home.y, 2) + POWER(st.z - home.z, 2))) as network_radius
FROM species s
         JOIN stars home ON s.home_system_id = home.id
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id
GROUP BY s.id;


-- Would you like to see views for analyzing production chains and resource flows?
-- Here are views for analyzing production chains and resource flows:

-- Production chain analysis
CREATE VIEW v_production_chains AS
SELECT
    s.name as species_name,
    st.x, st.y, st.z,
    COUNT(DISTINCT p.id) as production_nodes,
    SUM(r.mining_units) as resource_extraction,
    SUM(r.manufacturing_units) as processing_capacity,
    SUM(r.auto_IUs) as automation_level,
    COUNT(DISTINCT sh.id) as transport_capacity
FROM species s
         JOIN resources r ON s.id = r.species_id
         JOIN planets p ON r.planet_id = p.id
         JOIN stars st ON p.star_id = st.id
         LEFT JOIN ships sh ON sh.species_id = s.id
    AND sh.current_star_id = st.id
GROUP BY s.id, st.id;

-- Resource flow metrics
CREATE VIEW v_resource_flows AS
SELECT
    src.x as source_x, src.y as source_y, src.z as source_z,
    dest.x as dest_x, dest.y as dest_y, dest.z as dest_z,
    s.name as species_name,
    COUNT(DISTINCT sh.id) as transport_ships,
    SUM(r_src.mining_units) as source_production,
    SUM(r_dest.manufacturing_units) as destination_processing
FROM ships sh
         JOIN species s ON sh.species_id = s.id
         JOIN stars src ON sh.current_star_id = src.id
         JOIN stars dest ON m.unloading_point = dest.id
         JOIN military m ON sh.id = m.ship_id
         JOIN resources r_src ON s.id = r_src.species_id
         JOIN resources r_dest ON s.id = r_dest.species_id
WHERE m.unloading_point > 0
GROUP BY s.id, src.id, dest.id;

--  Copyright (c) 2024 Michael D Henderson. All rights reserved.

-- Let me propose a schema for storing orders:
-- This design:
--
-- Tracks orders by turn and species
-- Maintains order sequence within sections
-- Stores command parameters as JSON for flexibility
-- Includes processing status and timestamp
-- Uses indexes for efficient order execution

-- Orders table for storing all player commands
CREATE TABLE orders (
                        id INTEGER PRIMARY KEY,
                        turn_number INTEGER NOT NULL,
                        species_id INTEGER NOT NULL,
                        section TEXT NOT NULL,  -- PRE-DEPARTURE, JUMP, PRODUCTION, etc
                        sequence INTEGER NOT NULL,  -- Order within section
                        command TEXT NOT NULL,
                        parameters TEXT NOT NULL,  -- JSON encoded parameters
                        status TEXT DEFAULT 'pending',
                        processed_at DATETIME,
                        FOREIGN KEY(species_id) REFERENCES species(id)
);

-- Index for efficient order processing
CREATE INDEX idx_orders_processing ON orders(
                                             turn_number,
                                             section,
                                             sequence,
                                             status
    );

-- Order validation rules
CREATE TABLE order_rules (
                             id INTEGER PRIMARY KEY,
                             section TEXT NOT NULL,
                             command TEXT NOT NULL,
                             required_parameters TEXT NOT NULL,  -- JSON schema for validation
                             tech_requirements TEXT,  -- JSON encoded tech level requirements
                             validation_rules TEXT    -- JSON encoded rules
);

-- Order results tracking
CREATE TABLE order_results (
                               id INTEGER PRIMARY KEY,
                               order_id INTEGER NOT NULL,
                               turn_number INTEGER NOT NULL,
                               result_type TEXT NOT NULL,
                               result_data TEXT NOT NULL,  -- JSON encoded results
                               timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
                               FOREIGN KEY(order_id) REFERENCES orders(id)
);

-- Order execution log
CREATE TABLE order_execution_log (
                                     id INTEGER PRIMARY KEY,
                                     order_id INTEGER NOT NULL,
                                     execution_phase TEXT NOT NULL,
                                     success BOOLEAN NOT NULL,
                                     message TEXT,
                                     timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
                                     FOREIGN KEY(order_id) REFERENCES orders(id)
);
