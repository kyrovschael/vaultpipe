// Command vaultpipe injects secrets from HashiCorp Vault into a child
// process environment without exposing them through the shell.
//
// Usage:
//
//	vaultpipe [flags] -- <command> [args...]
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/vaultpipe/internal/audit"
	"github.com/yourorg/vaultpipe/internal/config"
	"github.com/yourorg/vaultpipe/internal/env"
	"github.com/yourorg/vaultpipe/internal/health"
	"github.com/yourorg/vaultpipe/internal/process"
	"github.com/yourorg/vaultpipe/internal/vault"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Locate the "--" separator that divides vaultpipe flags from the
	// child command. Everything after "--" is the command to execute.
	cmdArgs, err := splitArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "vaultpipe: %v\n", err)
		return 2
	}

	// Load configuration from file / environment.
	cfgPath := os.Getenv("VAULTPIPE_CONFIG")
	if cfgPath == "" {
		cfgPath = "vaultpipe.yaml"
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "vaultpipe: config: %v\n", err)
		return 1
	}

	// Set up structured audit logger.
	var logger *audit.Logger
	if cfg.AuditLog != "" {
		fl, err := audit.NewFileLogger(cfg.AuditLog)
		if err != nil {
			fmt.Fprintf(os.Stderr, "vaultpipe: audit log: %v\n", err)
			return 1
		}
		defer fl.Close()
		logger = audit.NewLogger(fl)
	} else {
		logger = audit.NewLogger(os.Stderr)
	}

	// Wait for Vault to become healthy before proceeding.
	checker := health.NewChecker(cfg.Vault.Address)
	if err := health.WaitUntilHealthy(ctx, checker, health.DefaultRetryConfig()); err != nil {
		fmt.Fprintf(os.Stderr, "vaultpipe: vault health: %v\n", err)
		return 1
	}

	// Authenticate with Vault and retrieve secrets.
	vc, err := vault.NewClient(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "vaultpipe: vault client: %v\n", err)
		return 1
	}

	secrets, err := vault.SecretsToEnv(ctx, vc, cfg.Secrets)
	if err != nil {
		fmt.Fprintf(os.Stderr, "vaultpipe: fetch secrets: %v\n", err)
		return 1
	}

	logger.Info("secrets fetched", map[string]string{
		"count": fmt.Sprintf("%d", len(secrets)),
	})

	// Build the child process environment: start from a filtered OS
	// snapshot, then overlay vault secrets.
	base := env.DefaultEnv()
	merged := vault.MergeEnv(base, secrets)

	// Execute the child command with the enriched environment.
	exitCode, err := process.Run(ctx, cmdArgs, merged)
	if err != nil {
		fmt.Fprintf(os.Stderr, "vaultpipe: exec: %v\n", err)
		return 1
	}

	return exitCode
}

// splitArgs separates vaultpipe's own flags from the child command.
// It looks for the first "--" argument and returns everything after it.
// If no separator is found, an error is returned.
func splitArgs(args []string) ([]string, error) {
	for i, a := range args {
		if a == "--" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("no command specified after '--'")
			}
			return args[i+1:], nil
		}
	}
	// Allow omitting "--" when no vaultpipe flags are needed.
	if len(args) == 0 {
		return nil, fmt.Errorf("usage: vaultpipe [flags] -- <command> [args...]")
	}
	return args, nil
}
