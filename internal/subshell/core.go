// Package subshell provides facilities to execute CLI commands in subshells.
package subshell

import "time"

// the number of times Git Town should retry when there is another Git process running
const concurrentGitRetries = 5

// the amount of time Git Town should wait between retries when there is another Git process running
const concurrentGitRetryDelay = 1 * time.Second
