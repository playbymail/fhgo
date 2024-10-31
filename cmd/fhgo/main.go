// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements the fhgo command.
package main

import (
	"fmt"
	"github.com/mdhender/semver"
	"github.com/playbymail/fhgo"
	"github.com/playbymail/fhgo/prng"
	"github.com/playbymail/fhgo/sqlc/sqlite3"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Lshortfile)

	cmdRoot.PersistentFlags().StringVarP(&argsRoot.db.path, "database", "D", "", "path to the database file")
	cmdRoot.PersistentFlags().BoolP("test", "t", false, "enable test mode")
	cmdRoot.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")

	cmdRoot.AddCommand(
		cmdCombat,
		cmdCreate,
		cmdDb,
		cmdExport,
		cmdFinish,
		cmdImport,
		cmdInspect,
		cmdJump,
		cmdList,
		cmdLocations,
		cmdLogRandom,
		cmdPostArrival,
		cmdPreDeparture,
		cmdProduction,
		cmdReport,
		cmdScan,
		cmdScanNear,
		cmdSexpr,
		cmdShow,
		cmdStats,
		cmdTurn,
		cmdUpdate,
		cmdVersion,
	)

	cmdCreate.AddCommand(cmdCreateGalaxy, cmdCreateHomeSystemTemplates, cmdCreateSpecies)
	cmdCreateGalaxy.Flags().BoolVar(&argsCreateGalaxy.deriveSizes, "derive-sizes", false, "derive radius and number of stars from number of species")
	cmdCreateGalaxy.Flags().BoolVar(&argsCreateGalaxy.lessCrowded, "less-crowded", false, "increases number of stars by 50% for slower-paced games")
	cmdCreateGalaxy.Flags().BoolVar(&argsCreateGalaxy.suggestValues, "suggest-values", false, "display suggested values based on number of species")
	cmdCreateGalaxy.Flags().IntVar(&argsCreateGalaxy.minimumRadiusInParsecs, "radius", 6, "minimum radius of the galaxy in parsecs")
	cmdCreateGalaxy.Flags().IntVar(&argsCreateGalaxy.numberOfSpecies, "species", 1, "defines number of species")
	cmdCreateGalaxy.Flags().IntVar(&argsCreateGalaxy.numberOfStarSystems, "stars", 12, "number of star systems to create")
	cmdCreateGalaxy.Flags().Uint64Var(&argsCreateGalaxy.prngSeed, "seed", 0, "seed for the random number generator")

	cmdDb.AddCommand(cmdDbInit)
	cmdDbInit.Flags().BoolVar(&argsRoot.db.forceCreate, "force", false, "delete database if it exists")
	cmdDbInit.Flags().StringVar(&argsRoot.db.code, "code", "FH", "code to assign to the game")
	cmdDbInit.Flags().StringVar(&argsRoot.db.name, "name", "gamma", "name to assign to the game")
	cmdDbInit.Flags().StringVar(&argsRoot.db.description, "description", "", "description of the game")
	cmdRoot.AddCommand(cmdVersion)

	cmdScan.AddCommand(cmdScanNear)

	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}

var (
	version = semver.Version{Major: 0, Minor: 0, Patch: 1, PreRelease: "dev"}

	argsRoot struct {
		db struct {
			path        string // path to the database directory
			code        string // code for the game
			description string // description of the game
			forceCreate bool   // if true, overwrite existing database
			name        string // name of the game
		}
	}

	cmdRoot = &cobra.Command{
		Use:   "fhgo",
		Short: "parent command for our application",
		Long:  `Not sure yet what this will do.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Hello from root command\n")
		},
	}

	cmdCombat = &cobra.Command{
		Use:   "combat",
		Short: "combat stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'combat' not yet implemented\n")
		},
	}

	cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "create galaxy, species, and home systems stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'create' not yet implemented\n")
		},
	}

	argsCreateGalaxy = struct {
		path                   string // path to the database directory
		deriveSizes            bool   // when set, calculate values for radius and number of stars
		lessCrowded            bool   // when set, increases number of stars by 50% for slower-paced games
		minimumRadiusInParsecs int
		numberOfSpecies        int
		numberOfStarSystems    int
		prngSeed               uint64
		suggestValues          bool // when set, displays suggested values based on number of species and exits
	}{}

	cmdCreateGalaxy = &cobra.Command{
		Use:   "galaxy",
		Short: "create galaxy stub",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if argsCreateGalaxy.minimumRadiusInParsecs < fhgo.MIN_RADIUS || argsCreateGalaxy.minimumRadiusInParsecs > fhgo.MAX_RADIUS {
				return fmt.Errorf("minimum radius must be between %d and %d parsecs", fhgo.MIN_RADIUS, fhgo.MAX_RADIUS)
			} else if argsCreateGalaxy.numberOfSpecies < fhgo.MIN_SPECIES || argsCreateGalaxy.numberOfSpecies > fhgo.MAX_SPECIES {
				return fmt.Errorf("species must be between %d and %d", fhgo.MIN_SPECIES, fhgo.MAX_SPECIES)
			} else if argsCreateGalaxy.numberOfStarSystems < fhgo.MIN_STARS || argsCreateGalaxy.numberOfStarSystems > fhgo.MAX_STARS {
				return fmt.Errorf("star systems must be between %d and %d", fhgo.MIN_STARS, fhgo.MAX_STARS)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if argsCreateGalaxy.suggestValues || argsCreateGalaxy.deriveSizes {
				derivedNumberOfStars := (argsCreateGalaxy.numberOfSpecies * fhgo.STANDARD_NUMBER_OF_STAR_SYSTEMS) / fhgo.STANDARD_NUMBER_OF_SPECIES
				if argsCreateGalaxy.lessCrowded {
					// bump the number of stars by 50% to make it take longer to encounter other species.
					derivedNumberOfStars = 3 * derivedNumberOfStars / 2
				}
				if derivedNumberOfStars > fhgo.MAX_STARS {
					log.Fatalf("error: calculation results in a number greater than %d stars.\n", fhgo.MAX_STARS)
				}

				// radius cubed divided by number of stars. this should be actually compute using stars and density.
				minVolume := derivedNumberOfStars * fhgo.STANDARD_GALACTIC_RADIUS * fhgo.STANDARD_GALACTIC_RADIUS * fhgo.STANDARD_GALACTIC_RADIUS / fhgo.STANDARD_NUMBER_OF_STAR_SYSTEMS
				derivedGalacticRadius := fhgo.MIN_RADIUS
				for derivedGalacticRadius*derivedGalacticRadius*derivedGalacticRadius < minVolume {
					derivedGalacticRadius++
				}
				if derivedGalacticRadius > fhgo.MAX_RADIUS {
					log.Fatalf("error: calculation results in a radius greater than %d parsecs.\n", fhgo.MAX_RADIUS)
				}

				if argsCreateGalaxy.lessCrowded {
					fmt.Printf(" info: for %6d species, %6d systems are needed for a less crowded galaxy.\n", argsCreateGalaxy.numberOfSpecies, derivedNumberOfStars)
				} else {
					fmt.Printf(" info: for %6d species, %6d stars  are needed for a normal density galaxy.\n", argsCreateGalaxy.numberOfSpecies, derivedNumberOfStars)

				}
				fmt.Printf(" info: for %6d stars  , %6d radius should be large enough for the galaxy.\n", derivedNumberOfStars, derivedGalacticRadius)

				if !argsCreateGalaxy.deriveSizes {
					// nothing else to do
					return
				}

				if argsCreateGalaxy.deriveSizes {
					argsCreateGalaxy.minimumRadiusInParsecs = derivedGalacticRadius
					argsCreateGalaxy.numberOfStarSystems = derivedNumberOfStars
				}
			}
			fmt.Printf(" info: creating new system with radius %6d, stars %6d, species %6d\n", argsCreateGalaxy.minimumRadiusInParsecs, argsCreateGalaxy.numberOfStarSystems, argsCreateGalaxy.numberOfSpecies)

			g := fhgo.CreateGalaxy(argsCreateGalaxy.path, argsCreateGalaxy.minimumRadiusInParsecs, argsCreateGalaxy.numberOfStarSystems, argsCreateGalaxy.numberOfSpecies, argsCreateGalaxy.prngSeed)
			log.Printf(" info: created galaxy with type %T\n", g)
		},
	}

	cmdCreateHomeSystemTemplates = &cobra.Command{
		Use:   "home-system-templates",
		Short: "create home-system-templates stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'create home-system-templates' not yet implemented\n")
		},
	}

	cmdCreateSpecies = &cobra.Command{
		Use:   "species",
		Short: "create species stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'create species' not yet implemented\n")
		},
	}

	cmdDb = &cobra.Command{
		Use:   "db",
		Short: "Database management commands",
	}

	cmdDbInit = &cobra.Command{
		Use:   "init",
		Short: "Initialize the database",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if argsRoot.db.path == "" {
				return fmt.Errorf("database: path is required\n")
			} else if path, err := filepath.Abs(argsRoot.db.path); err != nil {
				return fmt.Errorf("database: %v\n", err)
			} else {
				argsRoot.db.path = path
			}
			if argsRoot.db.code == "" {
				argsRoot.db.code = "FH"
			}
			if argsRoot.db.name == "" {
				argsRoot.db.name = "gamma"
			}
			if argsRoot.db.description == "" {
				argsRoot.db.description = fmt.Sprintf("%s (%s)", argsRoot.db.name, argsRoot.db.code)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("db: init: database  %s\n", argsRoot.db.path)

			// create the database
			log.Printf("db: init: creating database in %s\n", argsRoot.db.path)
			err := sqlite3.DatabaseCreate(argsRoot.db.path, argsRoot.db.forceCreate)
			if err != nil {
				log.Fatalf("db: init: %v\n", err)
			}
			log.Printf("db: created %q\n", argsRoot.db.path)

			// initialize a new game in the database
			log.Printf("db: init: creating new game\n")
			log.Printf("db: init: todo: implement the game initialization\n")
			log.Printf("db: created new game\n")
		},
	}

	cmdExport = &cobra.Command{
		Use:   "export",
		Short: "export stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'export' not yet implemented\n")
		},
	}

	cmdFinish = &cobra.Command{
		Use:   "finish",
		Short: "finish stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'finish' not yet implemented\n")
		},
	}

	cmdImport = &cobra.Command{
		Use:   "import",
		Short: "import stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'import' not yet implemented\n")
		},
	}

	cmdInspect = &cobra.Command{
		Use:   "inspect",
		Short: "inspect stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'inspect' not yet implemented\n")
		},
	}

	cmdJump = &cobra.Command{
		Use:   "jump",
		Short: "jump stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'jump' not yet implemented\n")
		},
	}

	cmdList = &cobra.Command{
		Use:   "list",
		Short: "list stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'list' not yet implemented\n")
		},
	}

	cmdLocations = &cobra.Command{
		Use:   "locations",
		Short: "locations stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'locations' not yet implemented\n")
		},
	}

	// logRandomCommand generates random numbers using the historical default seed value.
	cmdLogRandom = &cobra.Command{
		Use:   "log-random",
		Short: "log-random stub",
		Long:  `Generate random numbers using the historical default seed value.`,
		Run: func(cmd *cobra.Command, args []string) {
			// use the historical default seed value
			prng.SetSeed(prng.DefaultHistoricalSeedValue())
			// then print out a nice set of random values
			for i := 0; i < 1_000_000; i++ {
				r := prng.Rand(1024 * 1024)
				if i < 10 {
					fmt.Printf("%9d %9d %s\n", i, r, prng.String())
				} else if 1000 < i && i < 1010 {
					fmt.Printf("%9d %9d %s\n", i, r, prng.String())
				} else if (i % 85713) == 0 {
					fmt.Printf("%9d %9d %s\n", i, r, prng.String())
				}
			}
		},
	}

	cmdPostArrival = &cobra.Command{
		Use:   "post-arrival",
		Short: "post-arrival stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'post-arrival' not yet implemented\n")
		},
	}

	cmdPreDeparture = &cobra.Command{
		Use:   "pre-departure",
		Short: "pre-departure stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'pre-departure' not yet implemented\n")
		},
	}

	cmdProduction = &cobra.Command{
		Use:   "production",
		Short: "production stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'production' not yet implemented\n")
		},
	}

	cmdReport = &cobra.Command{
		Use:   "report",
		Short: "report stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'report' not yet implemented\n")
		},
	}

	cmdScan = &cobra.Command{
		Use:   "scan",
		Short: "scan stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'scan' not yet implemented\n")
		},
	}

	cmdScanNear = &cobra.Command{
		Use:   "near",
		Short: "scan near stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'scan near' not yet implemented\n")
		},
	}

	cmdSexpr = &cobra.Command{
		Use:   "sexpr",
		Short: "sexpr stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'sexpr' not yet implemented\n")
		},
	}

	cmdShow = &cobra.Command{
		Use:   "show",
		Short: "show stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'show' not yet implemented\n")
		},
	}

	cmdStats = &cobra.Command{
		Use:   "stats",
		Short: "stats stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'stats' not yet implemented\n")
		},
	}

	cmdTurn = &cobra.Command{
		Use:   "turn",
		Short: "turn stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'turn' not yet implemented\n")
		},
	}

	cmdUpdate = &cobra.Command{
		Use:   "update",
		Short: "update stub",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Command 'update' not yet implemented\n")
		},
	}

	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of this application",
		Long:  `Version of the server application. This is not the version of the OttoMap application used to create the maps.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", version.String())
		},
	}
)
