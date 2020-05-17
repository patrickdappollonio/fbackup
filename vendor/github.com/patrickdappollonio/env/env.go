package env

import (
	"os"
	"strconv"
	"strings"
)

// GetDefault returns the value set in the
// environment variable or the default value.
// Note that this value will return "defValue"
// if the value is an empty string.
func GetDefault(key, defValue string) string {
	if s := strings.TrimSpace(os.Getenv(key)); s != "" {
		return s
	}

	return defValue
}

// GetBoolean will return the value as a boolean
// if set in the environment. Do note that a non
// valid "true" value will always return false.
// Among the valid values you can use as true or
// false, this function accepts 1, t, T, TRUE, true,
// True, 0, f, F, FALSE, false, False.
func GetBoolean(key string) bool {
	v := GetDefault(key, "")

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}

	return b
}
