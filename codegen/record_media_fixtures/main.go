package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"seall/internal/testutil"
)

func main() {
	var (
		runPattern string
		count      int
		verbose    bool
		dryRun     bool
	)

	flag.StringVar(&runPattern, "run", "", "Optional go test -run pattern for limiting fixture refresh")
	flag.IntVar(&count, "count", 1, "go test count value")
	flag.BoolVar(&verbose, "v", true, "Run go test with -v")
	flag.BoolVar(&dryRun, "dry-run", false, "Print the go test command without executing it")
	flag.Parse()

	configPath := resolveConfigPath()
	if _, err := os.Stat(configPath); err != nil {
		fatalf("test config not found at %s; create it from test/config.example.toml first", configPath)
	}

	cfg := testutil.MustLoadConfig()
	if !cfg.Flags.EnableMediaApiTests {
		fatalf("SIMKL tests are disabled in %s; set flags.enable_simkl_tests=true", configPath)
	}
	if strings.TrimSpace(cfg.Provider.MediaApiJwt) == "" {
		fatalf("provider.simkl_jwt is empty in %s; fixture recording requires an authenticated SIMKL token", configPath)
	}
	if strings.TrimSpace(cfg.Provider.MediaApiUsername) == "" {
		fmt.Fprintf(os.Stderr, "warning: provider.simkl_username is empty in %s; collection-based refresh flows may not cover user-scoped fixtures\n", configPath)
	}

	packages := flag.Args()
	if len(packages) == 0 {
		packages = []string{"./internal/api/simkl"}
	}

	args := []string{"test"}
	if verbose {
		args = append(args, "-v")
	}
	if count > 0 {
		args = append(args, fmt.Sprintf("-count=%d", count))
	}
	if runPattern != "" {
		args = append(args, "-run", runPattern)
	}
	args = append(args, packages...)

	cmd := exec.Command("go", args...)
	cmd.Dir = testutil.ProjectRoot()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), testutil.RecordMediaApiFixturesEnvName+"=true")

	fmt.Fprintf(os.Stderr, "Using test config: %s\n", configPath)
	fmt.Fprintf(os.Stderr, "Recording SIMKL fixtures with %s=true\n", testutil.RecordMediaApiFixturesEnvName)
	fmt.Fprintf(os.Stderr, "Running: go %s\n", strings.Join(args, " "))

	if dryRun {
		return
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fatalf("failed to run go test: %v", err)
	}
}

func resolveConfigPath() string {
	configDir := os.Getenv("TEST_CONFIG_PATH")
	if configDir == "" {
		configDir = filepath.Join(testutil.ProjectRoot(), "test")
	}

	return filepath.Join(configDir, "config.toml")
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
