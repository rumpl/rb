package root

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rumpl/rb/pkg/environment"
	"github.com/rumpl/rb/pkg/paths"
)

type rootFlags struct {
	enableOtel  bool
	debugMode   bool
	logFilePath string
	logFile     *os.File
}

func NewRootCmd() *cobra.Command {
	var flags rootFlags

	cmd := &cobra.Command{
		Use:   "rb",
		Short: "rb - AI agent runner",
		Long:  "rb is a command-line tool for running AI agents",
		Example: `  rb run ./agent.yaml
  rb run agentcatalog/pirate`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Initialize logging before anything else so logs don't break TUI
			if err := flags.setupLogging(); err != nil {
				// If logging setup fails, fall back to stderr so we still get logs
				slog.SetDefault(slog.New(slog.NewTextHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{
					Level: func() slog.Level {
						if flags.debugMode {
							return slog.LevelDebug
						}
						return slog.LevelInfo
					}(),
				})))
			}

			if flags.enableOtel {
				if err := initOTelSDK(cmd.Context()); err != nil {
					slog.Warn("Failed to initialize OpenTelemetry SDK", "error", err)
				} else {
					slog.Debug("OpenTelemetry SDK initialized successfully")
				}
			}

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if flags.logFile != nil {
				_ = flags.logFile.Close()
			}
			return nil
		},
		// If no subcommand is specified, show help
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Add persistent debug flag available to all commands
	cmd.PersistentFlags().BoolVarP(&flags.debugMode, "debug", "d", false, "Enable debug logging")
	cmd.PersistentFlags().BoolVarP(&flags.enableOtel, "otel", "o", false, "Enable OpenTelemetry tracing")
	cmd.PersistentFlags().StringVar(&flags.logFilePath, "log-file", "", "Path to debug log file (default: ~/.rb/rb.debug.log; only used with --debug)")

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newRunCmd())
	cmd.AddCommand(newExecCmd())
	cmd.AddCommand(newNewCmd())
	cmd.AddCommand(newAPICmd())
	cmd.AddCommand(newACPCmd())
	cmd.AddCommand(newMCPCmd())
	cmd.AddCommand(newA2ACmd())
	cmd.AddCommand(newEvalCmd())
	cmd.AddCommand(newPushCmd())
	cmd.AddCommand(newPullCmd())

	// Define groups
	cmd.AddGroup(&cobra.Group{ID: "core", Title: "Core Commands:"})
	cmd.AddGroup(&cobra.Group{ID: "advanced", Title: "Advanced Commands:"})
	cmd.AddGroup(&cobra.Group{ID: "server", Title: "Server Commands:"})

	return cmd
}

func Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, args ...string) error {
	rootCmd := NewRootCmd()
	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		envErr := &environment.RequiredEnvError{}
		runtimeErr := RuntimeError{}

		switch {
		case ctx.Err() != nil:
			return ctx.Err()
		case errors.As(err, &envErr):
			fmt.Fprintln(stderr, "The following environment variables must be set:")
			for _, v := range envErr.Missing {
				fmt.Fprintf(stderr, " - %s\n", v)
			}
			fmt.Fprintln(stderr, "\nEither:\n - Set those environment variables before running rb\n - Run rb with --env-from-file\n - Store those secrets using one of the built-in environment variable providers.")
		case errors.As(err, &runtimeErr):
			// Runtime errors have already been printed by the command itself
			// Don't print them again or show usage
		default:
			// Command line usage errors - show the error and usage
			fmt.Fprintln(stderr, err)
			fmt.Fprintln(stderr)
			if strings.HasPrefix(err.Error(), "unknown command ") || strings.HasPrefix(err.Error(), "accepts ") {
				_ = rootCmd.Usage()
			}
		}

		return err
	}

	return nil
}

// setupLogging configures slog logging behavior.
// When --debug is enabled, logs are written to a single file <dataDir>/rb.debug.log (append mode),
// or to the file specified by --log-file.
func (f *rootFlags) setupLogging() error {
	if !f.debugMode {
		slog.SetDefault(slog.New(slog.DiscardHandler))
		return nil
	}

	path := strings.TrimSpace(f.logFilePath)
	if path == "" {
		path = filepath.Join(paths.GetDataDir(), "rb.debug.log")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	f.logFile = logFile

	slog.SetDefault(slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug})))

	return nil
}

// RuntimeError wraps runtime errors to distinguish them from usage errors
type RuntimeError struct {
	Err error
}

func (e RuntimeError) Error() string {
	return e.Err.Error()
}

func (e RuntimeError) Unwrap() error {
	return e.Err
}

// isFirstRun checks if this is the first time rb is being run
// It creates a marker file in the user's config directory
func isFirstRun() bool {
	configDir := paths.GetConfigDir()
	markerFile := filepath.Join(configDir, ".rb_first_run")

	// Check if marker file exists
	if _, err := os.Stat(markerFile); err == nil {
		return false // File exists, not first run
	}

	// Create marker file to indicate this run has happened
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return false // Can't create config dir, assume not first run
	}

	if err := os.WriteFile(markerFile, []byte(""), 0o644); err != nil {
		return false // Can't create marker file, assume not first run
	}

	return true // Successfully created marker, this is first run
}
