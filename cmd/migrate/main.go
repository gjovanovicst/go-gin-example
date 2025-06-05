package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/EDDYCJY/go-gin-example/pkg/migration"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

func main() {
	var (
		action      = flag.String("action", "up", "Migration action: up, down, version, migrate, force, seed, seed-rollback, seed-status, list-envs")
		version     = flag.String("version", "", "Migration version (for migrate and force actions)")
		environment = flag.String("env", "", "Environment for seeding (development, production, etc.)")
	)
	flag.Parse()

	// Load configuration
	setting.Setup()

	switch *action {
	case "up":
		if err := migration.RunMigrations(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	case "down":
		if err := migration.RollbackMigrations(); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
	case "version":
		version, dirty, err := migration.GetMigrationVersion()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		if version == 0 {
			fmt.Println("No migrations have been applied")
		} else {
			fmt.Printf("Current migration version: %d", version)
			if dirty {
				fmt.Print(" (dirty)")
			}
			fmt.Println()
		}
	case "migrate":
		if *version == "" {
			log.Fatal("Version is required for migrate action")
		}
		v, err := strconv.ParseUint(*version, 10, 32)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := migration.MigrateToVersion(uint(v)); err != nil {
			log.Fatalf("Failed to migrate to version %d: %v", v, err)
		}
	case "force":
		if *version == "" {
			log.Fatal("Version is required for force action")
		}
		v, err := strconv.ParseUint(*version, 10, 32)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := migration.ForceMigrationVersion(uint(v)); err != nil {
			log.Fatalf("Failed to force migration version %d: %v", v, err)
		}
	case "seed":
		if err := migration.RunSeeds(*environment); err != nil {
			log.Fatalf("Failed to run seeds: %v", err)
		}
	case "seed-rollback":
		if err := migration.RollbackSeeds(*environment); err != nil {
			log.Fatalf("Failed to rollback seeds: %v", err)
		}
	case "seed-status":
		version, dirty, err := migration.GetSeedStatus(*environment)
		if err != nil {
			log.Fatalf("Failed to get seed status: %v", err)
		}
		env := *environment
		if env == "" {
			env = migration.GetEnvironment()
		}
		if version == 0 {
			fmt.Printf("No seeds have been applied for environment: %s\n", env)
		} else {
			fmt.Printf("Current seed version for %s: %d", env, version)
			if dirty {
				fmt.Print(" (dirty)")
			}
			fmt.Println()
		}
	case "list-envs":
		environments, err := migration.ListAvailableEnvironments()
		if err != nil {
			log.Fatalf("Failed to list environments: %v", err)
		}
		if len(environments) == 0 {
			fmt.Println("No seed environments found")
		} else {
			fmt.Printf("Available seed environments: %s\n", strings.Join(environments, ", "))
		}
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s -action=[up|down|version|migrate|force|seed|seed-rollback|seed-status|list-envs] [-version=N] [-env=ENV]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nMigration Actions:\n")
		fmt.Fprintf(os.Stderr, "  up      - Run all pending migrations\n")
		fmt.Fprintf(os.Stderr, "  down    - Rollback the last migration\n")
		fmt.Fprintf(os.Stderr, "  version - Show current migration version\n")
		fmt.Fprintf(os.Stderr, "  migrate - Migrate to a specific version\n")
		fmt.Fprintf(os.Stderr, "  force   - Force set migration version (use with caution)\n")
		fmt.Fprintf(os.Stderr, "\nSeeding Actions:\n")
		fmt.Fprintf(os.Stderr, "  seed          - Run seeds for specified environment\n")
		fmt.Fprintf(os.Stderr, "  seed-rollback - Rollback last seed for specified environment\n")
		fmt.Fprintf(os.Stderr, "  seed-status   - Show current seed status for environment\n")
		fmt.Fprintf(os.Stderr, "  list-envs     - List available seed environments\n")
		fmt.Fprintf(os.Stderr, "\nMigration Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s -action=up\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=down\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=version\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=migrate -version=1\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=force -version=0\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nSeeding Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s -action=seed -env=development\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=seed -env=production\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=seed (uses current environment)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=seed-status -env=development\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=list-envs\n", os.Args[0])
		os.Exit(1)
	}
}