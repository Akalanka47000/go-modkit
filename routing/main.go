package routing

import (
	"fmt"
)

// Returns a function that formats a given module name with a version number. Useful for generating versioned module prefixes in REST APIs.
// Example usage: routing.VersionablePrefix("analysis")(1) returns "/v1/analysis"
func VersionablePrefix(moduleName string) func(int) string {
	return func(version int) string {
		return fmt.Sprintf("/v%d/%s", version, moduleName)
	}
}
