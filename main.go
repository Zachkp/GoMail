package main

import (
	"fmt"
	"os"

	"github.com/Zachkp/GoMail/config"
	"github.com/Zachkp/GoMail/models"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check for command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			handleConfigCommand()
			return
		case "help", "-h", "--help":
			printHelp()
			return
		case "version", "-v", "--version":
			fmt.Println("GoMail v1.0.0")
			return
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			printHelp()
			os.Exit(1)
		}
	}

	// Try to load configuration and start the TUI
	if _, err := config.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nTip: Run 'GoMail config' to manage your configuration\n")
		os.Exit(1)
	}

	// Start the TUI
	if _, err := tea.NewProgram(models.CreateTable(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func handleConfigCommand() {
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "path":
			configPath, err := config.GetConfigPath()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting config path: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(configPath)
		case "init":
			if err := config.InitConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
				os.Exit(1)
			}
		case "validate":
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Configuration validation failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Configuration is valid!")
			fmt.Printf("Username: %s\n", cfg.EmailUsername)
			fmt.Printf("IMAP Host: %s:%s\n", cfg.EmailImapHost, cfg.EmailImapPort)
			fmt.Printf("SMTP Host: %s:%s\n", cfg.EmailSmtpHost, cfg.EmailSmtpPort)
		default:
			fmt.Printf("Unknown config command: %s\n", os.Args[2])
			printConfigHelp()
			os.Exit(1)
		}
	} else {
		printConfigHelp()
	}
}

func printHelp() {
	fmt.Println(`GoMail - A terminal-based email viewer with fuzzy search

Usage:
  GoMail            Start the email client
  GoMail config     Manage configuration
  GoMail help       Show this help message
  GoMail version    Show version information

For configuration management, use:
  GoMail config init      Initialize configuration
  GoMail config path      Show configuration file path
  GoMail config validate  Validate current configuration`)
}

func printConfigHelp() {
	fmt.Println(`Configuration management commands:

  config init      Create a new configuration file with default values
  config path      Show the path to the configuration file
  config validate  Check if the current configuration is valid

Configuration file location:
  The configuration file is stored at ~/.config/GoMail/.env

Example configuration:
  EMAIL_USERNAME=your.email@gmail.com
  EMAIL_PASSWORD=your-app-password
  EMAIL_IMAP_HOST=imap.gmail.com
  EMAIL_IMAP_PORT=993
  EMAIL_SMTP_HOST=smtp.gmail.com
  EMAIL_SMTP_PORT=587

Note: For Gmail, you'll need to use an "App Password" instead of your regular password.`)
}
