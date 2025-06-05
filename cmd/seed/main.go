package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/EDDYCJY/go-gin-example/pkg/migration"
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

func main() {
	var (
		action      = flag.String("action", "run", "Seeding action: run, rollback, status, list, detailed-status")
		environment = flag.String("env", "", "Environment (development, production, etc.)")
		manual      = flag.Bool("manual", false, "Use manual seeding (direct SQL execution)")
		tracked     = flag.Bool("tracked", true, "Use tracked seeding (default, recommended)")
	)
	flag.Parse()

	// Load configuration
	setting.Setup()

	switch *action {
	case "run":
		if *manual {
			if err := migration.RunSeedsManually(*environment); err != nil {
				log.Fatalf("Failed to run seeds manually: %v", err)
			}
		} else if *tracked {
			if err := migration.RunSeedsWithTracking(*environment); err != nil {
				log.Fatalf("Failed to run tracked seeds: %v", err)
			}
		} else {
			if err := migration.RunSeeds(*environment); err != nil {
				log.Fatalf("Failed to run seeds: %v", err)
			}
		}
	case "rollback":
		if *tracked {
			if err := migration.RollbackLastSeed(*environment); err != nil {
				log.Fatalf("Failed to rollback seed: %v", err)
			}
		} else {
			if err := migration.RollbackSeeds(*environment); err != nil {
				log.Fatalf("Failed to rollback seeds: %v", err)
			}
		}
	case "status":
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
	case "list":
		environments, err := migration.ListAvailableEnvironments()
		if err != nil {
			log.Fatalf("Failed to list environments: %v", err)
		}
		if len(environments) == 0 {
			fmt.Println("No seed environments found")
		} else {
			fmt.Printf("Available seed environments: %s\n", strings.Join(environments, ", "))
		}
		
		// Also show current environment
		currentEnv := migration.GetEnvironment()
		fmt.Printf("Current environment: %s\n", currentEnv)
	case "detailed-status":
		env := *environment
		if env == "" {
			env = migration.GetEnvironment()
		}
		status, err := migration.GetDetailedSeedStatus(env)
		if err != nil {
			log.Fatalf("Failed to get detailed seed status: %v", err)
		}
		
		fmt.Printf("=== Detailed Seed Status for %s ===\n", env)
		fmt.Printf("Total Available Seeds: %v\n", status["total_available"])
		fmt.Printf("Total Applied: %v\n", status["total_applied"])
		fmt.Printf("Total Pending: %v\n", status["total_pending"])
		fmt.Printf("Latest Applied: %v\n", status["latest_applied"])
		
		if appliedSeeds, ok := status["applied_seeds"].([]string); ok && len(appliedSeeds) > 0 {
			fmt.Printf("\nApplied Seeds: %v\n", appliedSeeds)
		}
		
		if pendingSeeds, ok := status["pending_seeds"].([]string); ok && len(pendingSeeds) > 0 {
			fmt.Printf("Pending Seeds: %v\n", pendingSeeds)
		} else {
			fmt.Printf("No pending seeds\n")
		}
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s -action=[run|rollback|status|list|detailed-status] [-env=ENV] [-manual] [-tracked]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nActions:\n")
		fmt.Fprintf(os.Stderr, "  run            - Run seeds for specified environment\n")
		fmt.Fprintf(os.Stderr, "  rollback       - Rollback last seed for specified environment\n")
		fmt.Fprintf(os.Stderr, "  status         - Show current seed status for environment\n")
		fmt.Fprintf(os.Stderr, "  detailed-status- Show detailed seed information\n")
		fmt.Fprintf(os.Stderr, "  list           - List available seed environments\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		fmt.Fprintf(os.Stderr, "  -env      - Specify environment (defaults to current)\n")
		fmt.Fprintf(os.Stderr, "  -manual   - Use manual seeding (direct SQL execution, no tracking)\n")
		fmt.Fprintf(os.Stderr, "  -tracked  - Use tracked seeding with version control (default: true)\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -action=run -env=development           # Tracked seeding (recommended)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=run -env=development -manual   # Manual seeding (no tracking)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=rollback -env=development      # Rollback last tracked seed\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=detailed-status -env=development\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -action=list\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nEnvironment Detection:\n")
		fmt.Fprintf(os.Stderr, "  - Uses APP_ENV environment variable if set\n")
		fmt.Fprintf(os.Stderr, "  - Falls back to RunMode in config (debug=development, release=production)\n")
		os.Exit(1)
	}
}