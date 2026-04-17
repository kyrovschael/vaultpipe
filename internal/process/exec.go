package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Run executes the given command with the provided environment variables
// injected alongside the current process environment.
func Run(command []string, secrets map[string]string) error {
	if len(command) == 0 {
		return fmt.Errorf("no command provided")
	}

	path, err := exec.LookPath(command[0])
	if err != nil {
		return fmt.Errorf("could not find executable %q: %w", command[0], err)
	}

	env := mergeEnv(os.Environ(), secrets)

	cmd := exec.Command(path, command[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

// mergeEnv combines the base environment with secret key=value pairs.
// Secret values take precedence over existing env vars with the same key.
func mergeEnv(base []string, secrets map[string]string) []string {
	result := make([]string, 0, len(base)+len(secrets))

	// Track keys that will be overridden by secrets
	override := make(map[string]struct{}, len(secrets))
	for k := range secrets {
		override[k] = struct{}{}
	}

	for _, entry := range base {
		key := envKey(entry)
		if _, shadowed := override[key]; !shadowed {
			result = append(result, entry)
		}
	}

	for k, v := range secrets {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}

	return result
}

// envKey extracts the key portion from a "KEY=VALUE" string.
func envKey(entry string) string {
	for i := 0; i < len(entry); i++ {
		if entry[i] == '=' {
			return entry[:i]
		}
	}
	return entry
}
