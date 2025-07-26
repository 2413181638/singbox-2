package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/exec"
    "os/signal"
    "path/filepath"
    "syscall"
    "time"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "github.com/yourusername/singxclient/internal/config"
    "github.com/yourusername/singxclient/internal/xboard"
)

var version = "dev"

// main sets up the CLI and executes it.
func main() {
    root := &cobra.Command{
        Use:   "singxclient",
        Short: "XBoard ⇄ sing-box bridge client",
    }

    root.PersistentFlags().StringP("subscription", "s", "", "XBoard subscription URL (required)")
    root.PersistentFlags().StringP("data-dir", "d", ".", "Directory to store generated config & runtime files")
    root.PersistentFlags().DurationP("refresh", "r", time.Hour, "How often to refresh subscription & reload (0 to disable)")
    if err := viper.BindPFlags(root.PersistentFlags()); err != nil {
        log.Fatalf("bind flags: %v", err)
    }

    root.Version = version
    root.SetVersionTemplate("SingXClient version {{.Version}}\n")

    root.RunE = func(cmd *cobra.Command, args []string) error {
        subscription := viper.GetString("subscription")
        if subscription == "" {
            return fmt.Errorf("--subscription is required")
        }
        dataDir := viper.GetString("data-dir")
        refreshInterval := viper.GetDuration("refresh")

        return run(cmd.Context(), subscription, dataDir, refreshInterval)
    }

    if err := root.Execute(); err != nil {
        os.Exit(1)
    }
}

// run orchestrates fetching the subscription, generating sing-box configuration and launching the engine.
func run(ctx context.Context, subURL, dataDir string, refreshInterval time.Duration) error {
    if err := os.MkdirAll(dataDir, 0o755); err != nil {
        return fmt.Errorf("create data dir: %w", err)
    }

    cfgPath := filepath.Join(dataDir, "config.json")

    // Initial bootstrap
    if err := refreshConfig(subURL, cfgPath); err != nil {
        return err
    }

    // Start sing-box process
    procCtx, cancel := context.WithCancel(ctx)
    defer cancel()

    bin, err := exec.LookPath("sing-box")
    if err != nil {
        return fmt.Errorf("sing-box binary not found in PATH: %w", err)
    }

    cmd := exec.CommandContext(procCtx, bin, "run", "-c", cfgPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("launch sing-box: %w", err)
    }

    // Signal handling & optional auto-refresh
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

    var ticker *time.Ticker
    if refreshInterval > 0 {
        ticker = time.NewTicker(refreshInterval)
    }

    for {
        select {
        case sig := <-sigCh:
            if sig == syscall.SIGINT || sig == syscall.SIGTERM {
                cancel()
                _ = cmd.Wait()
                return nil
            }
            if sig == syscall.SIGHUP {
                // manual reload
                if err := refreshConfig(subURL, cfgPath); err != nil {
                    log.Printf("manual reload failed: %v", err)
                }
            }
        case <-func() <-chan time.Time {
            if ticker != nil {
                return ticker.C
            }
            return make(chan time.Time) // nil-safe channel (never fires)
        }():
            if ticker != nil {
                log.Println("Auto refreshing configuration …")
                if err := refreshConfig(subURL, cfgPath); err != nil {
                    log.Printf("refresh failed: %v", err)
                }
            }
        case <-ctx.Done():
            cancel()
            _ = cmd.Wait()
            return nil
        }
    }
}

// refreshConfig downloads subscription & writes sing-box config JSON to the given path.
func refreshConfig(subURL, dst string) error {
    profiles, err := xboard.FetchSubscription(subURL)
    if err != nil {
        return fmt.Errorf("fetch subscription: %w", err)
    }
    cfg, err := config.Generate(profiles)
    if err != nil {
        return fmt.Errorf("generate config: %w", err)
    }
    return os.WriteFile(dst, cfg, 0o644)
}