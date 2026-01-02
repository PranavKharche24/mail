package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pranavKharche24/mail/cli"
	"github.com/pranavKharche24/mail/config"
	"github.com/pranavKharche24/mail/mailer"
	"github.com/pranavKharche24/mail/web"
)

const version = "1.0.0"

func main() {
	// Load configuration
	cfg := config.Load()

	// Create mailer instance
	m := mailer.New()
	m.SetCredentials(cfg.EmailFrom, cfg.EmailPassword)

	// Check command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "cli", "-c", "--cli":
			runCLI(m)
		case "web", "-w", "--web":
			runWeb(cfg, m)
		case "version", "-v", "--version":
			fmt.Printf("Gomail v%s\n", version)
		case "help", "-h", "--help":
			printHelp()
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			printHelp()
			os.Exit(1)
		}
	} else {
		// Default: launch both web server and CLI
		runBoth(cfg, m)
	}
}

func runBoth(cfg *config.Config, m *mailer.Mailer) {
	printBanner()

	// Start web server in background
	go func() {
		server := web.New(cfg.Port, m)
		if err := server.Start(); err != nil {
			log.Printf("Web server error: %v", err)
		}
	}()

	fmt.Println()
	fmt.Println("  Services started:")
	fmt.Printf("  - Web Interface: http://localhost:%s\n", cfg.Port)
	fmt.Printf("  - Admin Panel:   http://localhost:%s/admin\n", cfg.Port)
	fmt.Println("  - CLI Interface: Active below")
	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))

	// Run CLI in foreground
	c := cli.New()
	email, pass := m.GetCredentials()
	if email != "" && pass != "" {
		c.SetMailer(m)
	}
	c.Run()
}

func runCLI(m *mailer.Mailer) {
	printBanner()
	c := cli.New()
	email, pass := m.GetCredentials()
	if email != "" && pass != "" {
		c.SetMailer(m)
	}
	c.Run()
}

func runWeb(cfg *config.Config, m *mailer.Mailer) {
	printBanner()
	server := web.New(cfg.Port, m)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func printBanner() {
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════════════════════╗")
	fmt.Println("  ║                                                      ║")
	fmt.Println("  ║    GOMAIL - Professional Email Utility               ║")
	fmt.Printf("  ║    Version %s                                       ║\n", version)
	fmt.Println("  ║                                                      ║")
	fmt.Println("  ╚══════════════════════════════════════════════════════╝")
}

func printHelp() {
	fmt.Println()
	fmt.Println("GOMAIL - Professional Email Utility")
	fmt.Println()
	fmt.Println("Usage: gomail [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  (none)             Start both Web and CLI interfaces")
	fmt.Println("  cli, -c, --cli     Start CLI interface only")
	fmt.Println("  web, -w, --web     Start Web interface only")
	fmt.Println("  version, -v        Show version")
	fmt.Println("  help, -h, --help   Show this help")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  Create a .env file with:")
	fmt.Println("    EMAIL_FROM=your-email@gmail.com")
	fmt.Println("    EMAIL_PASSWORD=your-app-password")
	fmt.Println("    PORT=8080")
	fmt.Println()
	fmt.Println("Documentation: https://github.com/pranavKharche24/mail")
	fmt.Println()
}
