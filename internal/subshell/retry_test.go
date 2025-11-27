package subshell_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/subshell"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBackendRunner_RetryOnIndexLock(t *testing.T) {
	t.Parallel()

	t.Run("detects index.lock error pattern correctly", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}

		testCases := []struct {
			name        string
			errorMsg    string
			shouldRetry bool
		}{
			{
				name:        "exact match",
				errorMsg:    "fatal: Unable to create '.git/index.lock': File exists.",
				shouldRetry: true,
			},
			{
				name:        "with path variations",
				errorMsg:    "fatal: Unable to create '/path/to/repo/.git/index.lock': File exists.",
				shouldRetry: true,
			},
			{
				name:        "missing 'fatal' prefix but has both patterns",
				errorMsg:    "Unable to create '.git/index.lock': File exists.",
				shouldRetry: false, // doesn't match because "fatal: Unable to create" is required
			},
			{
				name:        "only has 'Unable to create'",
				errorMsg:    "fatal: Unable to create '.git/something'",
				shouldRetry: false,
			},
			{
				name:        "only has 'File exists'",
				errorMsg:    "index.lock': File exists.",
				shouldRetry: false,
			},
			{
				name:        "completely different error",
				errorMsg:    "fatal: not a git repository",
				shouldRetry: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				scriptPath := filepath.Join(tmpDir, fmt.Sprintf("test-%s.sh", tc.name))
				scriptContent := fmt.Sprintf(`#!/bin/bash
>&2 echo "%s"
exit 1
`, tc.errorMsg)
				must.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o600))
				must.NoError(t, os.Chmod(scriptPath, 0o700))

				start := time.Now()
				_, err := runner.Query("bash", scriptPath)
				duration := time.Since(start)

				must.Error(t, err)
				if tc.shouldRetry {
					// Should take at least 4 seconds (exhausted retries)
					must.GreaterEq(t, 4*time.Second, duration)
				} else {
					// Should fail immediately
					must.Less(t, 500*time.Millisecond, duration)
				}
			})
		}
	})

	t.Run("does not retry on non-lock errors", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}

		// Create a script that fails with a different error
		scriptPath := filepath.Join(tmpDir, "other-error.sh")
		scriptContent := `#!/bin/bash
>&2 echo "fatal: Some other error"
exit 1
`
		must.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o600))
		must.NoError(t, os.Chmod(scriptPath, 0o700))

		start := time.Now()
		_, err := runner.Query("bash", scriptPath)
		duration := time.Since(start)

		// Should fail immediately
		must.Error(t, err)
		// Should complete quickly without retries
		must.Less(t, 500*time.Millisecond, duration)
	})

	t.Run("exhausts retries and fails after max attempts", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}

		// Create a script that always fails with lock error
		scriptPath := filepath.Join(tmpDir, "always-fails.sh")
		scriptContent := `#!/bin/bash
>&2 echo "fatal: Unable to create '.git/index.lock': File exists."
exit 1
`
		must.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o600))
		must.NoError(t, os.Chmod(scriptPath, 0o700))

		start := time.Now()
		_, err := runner.Query("bash", scriptPath)
		duration := time.Since(start)

		// Should fail after exhausting retries
		must.Error(t, err)
		// Should take at least 4 seconds (5 attempts with 4 delays between them)
		must.GreaterEq(t, 4*time.Second, duration)
	})

	t.Run("retries and succeeds on transient lock error", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}

		// Create a script that fails twice with lock error, then succeeds
		scriptPath := filepath.Join(tmpDir, "retry-script.sh")
		scriptContent := `#!/bin/bash
COUNTER_FILE="` + tmpDir + `/counter"
if [ ! -f "$COUNTER_FILE" ]; then
    echo "0" > "$COUNTER_FILE"
fi
COUNT=$(cat "$COUNTER_FILE")
echo $((COUNT + 1)) > "$COUNTER_FILE"

if [ "$COUNT" -lt "2" ]; then
    >&2 echo "fatal: Unable to create '.git/index.lock': File exists."
    exit 1
else
    echo "success"
    exit 0
fi
`
		must.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o600))
		must.NoError(t, os.Chmod(scriptPath, 0o700))

		start := time.Now()
		output, err := runner.Query("bash", scriptPath)
		duration := time.Since(start)

		must.NoError(t, err)
		must.EqOp(t, "success\n", output)
		// Should take at least 2 seconds (2 retries * 1 second delay)
		must.GreaterEq(t, 2*time.Second, duration)
	})

	t.Run("retries correct number of times", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}

		// Create a script that counts how many times it's called
		counterFile := filepath.Join(tmpDir, "attempt-counter")
		scriptPath := filepath.Join(tmpDir, "count-attempts.sh")
		scriptContent := `#!/bin/bash
COUNTER_FILE="` + counterFile + `"
if [ ! -f "$COUNTER_FILE" ]; then
    echo "1" > "$COUNTER_FILE"
else
    COUNT=$(cat "$COUNTER_FILE")
    echo $((COUNT + 1)) > "$COUNTER_FILE"
fi
>&2 echo "fatal: Unable to create '.git/index.lock': File exists."
exit 1
`
		must.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o600))
		must.NoError(t, os.Chmod(scriptPath, 0o700))

		_, err := runner.Query("bash", scriptPath)
		must.Error(t, err)

		// Read the counter to verify it was called exactly 5 times
		counterBytes, err := os.ReadFile(counterFile)
		must.NoError(t, err)
		must.EqOp(t, "5\n", string(counterBytes))
	})

	t.Run("succeeds immediately when no lock error", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		runner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}
		start := time.Now()
		output, err := runner.Query("echo", "success")
		duration := time.Since(start)
		must.NoError(t, err)
		must.EqOp(t, "success\n", output)
		// Should complete quickly without any retries
		must.Less(t, 100*time.Millisecond, duration)
	})
}

func TestFrontendRunner_RetryOnIndexLock(t *testing.T) {
	t.Parallel()

	// NOTE: FrontendRunner DOES properly retry because the subprocess is created
	// INSIDE the retry loop (unlike BackendRunner which has a bug).

	t.Run("does not retry on non-lock errors", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		backendRunner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}
		runner := &subshell.FrontendRunner{
			Backend:          backendRunner,
			GetCurrentBranch: nil,
			GetCurrentSHA:    nil,
			PrintBranchNames: false,
			PrintCommands:    false,
			CommandsCounter:  NewMutable(new(gohacks.Counter)),
		}

		// Script that fails with a different error
		scriptPath := filepath.Join(tmpDir, "different-error.sh")
		scriptContent := `#!/bin/bash
>&2 echo "fatal: repository not found"
exit 128
`
		must.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o600))
		must.NoError(t, os.Chmod(scriptPath, 0o700))

		start := time.Now()
		err := runner.Run("bash", scriptPath)
		duration := time.Since(start)

		// Should fail immediately
		must.Error(t, err)
		// Should complete quickly without retries
		must.Less(t, 500*time.Millisecond, duration)
	})

	t.Run("exhausts retries and fails", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		backendRunner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}
		runner := &subshell.FrontendRunner{
			Backend:          backendRunner,
			GetCurrentBranch: nil,
			GetCurrentSHA:    nil,
			PrintBranchNames: false,
			PrintCommands:    false,
			CommandsCounter:  NewMutable(new(gohacks.Counter)),
		}

		// Create a script that always fails with lock error
		scriptPath := filepath.Join(tmpDir, "always-lock.sh")
		scriptContent := `#!/bin/bash
>&2 echo "fatal: Unable to create '.git/index.lock': File exists."
exit 1
`
		must.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o600))
		must.NoError(t, os.Chmod(scriptPath, 0o700))

		start := time.Now()
		err := runner.Run("bash", scriptPath)
		duration := time.Since(start)

		// Should fail after exhausting retries
		must.Error(t, err)
		// Should take at least 4 seconds (5 attempts with 4 delays between them)
		must.GreaterEq(t, 4*time.Second, duration)
	})

	t.Run("retries and succeeds on transient lock error", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		backendRunner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}
		runner := &subshell.FrontendRunner{
			Backend:          backendRunner,
			GetCurrentBranch: nil,
			GetCurrentSHA:    nil,
			PrintBranchNames: false,
			PrintCommands:    false,
			CommandsCounter:  NewMutable(new(gohacks.Counter)),
		}

		// Create a script that fails once with lock error, then succeeds
		scriptPath := filepath.Join(tmpDir, "retry-once.sh")
		scriptContent := `#!/bin/bash
COUNTER_FILE="` + tmpDir + `/counter"
if [ ! -f "$COUNTER_FILE" ]; then
    echo "0" > "$COUNTER_FILE"
fi
COUNT=$(cat "$COUNTER_FILE")
echo $((COUNT + 1)) > "$COUNTER_FILE"

if [ "$COUNT" -lt "1" ]; then
    >&2 echo "fatal: Unable to create '.git/index.lock': File exists."
    exit 1
else
    echo "success"
    exit 0
fi
`
		must.NoError(t, os.WriteFile(scriptPath, []byte(scriptContent), 0o600))
		must.NoError(t, os.Chmod(scriptPath, 0o700))

		start := time.Now()
		err := runner.Run("bash", scriptPath)
		duration := time.Since(start)

		must.NoError(t, err)
		// Should take at least 1 second (1 retry * 1 second delay)
		must.GreaterEq(t, 1*time.Second, duration)
	})

	t.Run("succeeds immediately when no lock error", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		backendRunner := subshell.BackendRunner{Dir: Some(tmpDir), Verbose: false, CommandsCounter: NewMutable(new(gohacks.Counter))}
		runner := &subshell.FrontendRunner{
			Backend:          backendRunner,
			GetCurrentBranch: nil, // not needed for this test
			GetCurrentSHA:    nil, // not needed for this test
			PrintBranchNames: false,
			PrintCommands:    false,
			CommandsCounter:  NewMutable(new(gohacks.Counter)),
		}

		start := time.Now()
		err := runner.Run("echo", "success")
		duration := time.Since(start)

		must.NoError(t, err)
		// Should complete quickly without any retries
		must.Less(t, 100*time.Millisecond, duration)
	})
}
