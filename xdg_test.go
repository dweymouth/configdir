//go:build !windows && !darwin

package configdir_test

import (
	"os"
	"testing"

	"github.com/20after4/configdir"
)

func reset() {
	// Unset environment variables to make the tests deterministic and choose
	// their own default values. The individual test runners may override these
	// variables for their own use case.
	os.Setenv("HOME", "/home/user")
	os.Setenv("XDG_CONFIG_DIRS", "")
	os.Setenv("XDG_CONFIG_HOME", "")
	os.Setenv("XDG_CACHE_HOME", "")
	configdir.Refresh()
}

// TestCase provides common inputs and outputs for the functions being tested.
type TestCase struct {
	Env     string   // Value to put into the relevant environment variable, when Refresh=true
	Refresh bool     // Whether to call the Refresh() function before running
	Paths   []string // Path suffixes to give to the path function
	Values  []string // What we expect the return value(s) to contain
}

// On init, call the reset() function provided by the various OS-specific tests
// to reset environment variables to a known deterministic state.
func init() {
	reset()
}

// Common logic for the local paths, which return single values.
//
// Parameters:
//
//	t (*testing.T)
//	pathType (string): "config" or "cache", controls which environment
//	    variable to play with and which path function to call.
//	defaultPrefix (string): the default path prefix for the kind of path
//	    being tested, e.g. "/home/user/.config" for config paths.
//	customPrefix (string): when a custom value is set for the environment
//	    variable, this is that path prefix instead of the default.
func testLocalCommon(t *testing.T, pathType, defaultPrefix, customPrefix string) {
	reset()

	// Cases to test.
	var tests = []TestCase{
		{
			Paths:  []string{},
			Values: []string{defaultPrefix},
		},
		{
			Paths:  []string{"vendor-name"},
			Values: []string{defaultPrefix + "/vendor-name"},
		},
		{
			Paths:  []string{"vendor-name", "app-name"},
			Values: []string{defaultPrefix + "/vendor-name/app-name"},
		},

		// With custom XDG paths...
		{
			Env:     customPrefix,
			Refresh: true,
			Paths:   []string{},
			Values:  []string{customPrefix},
		},
		{
			Paths:  []string{"vendor-name"},
			Values: []string{customPrefix + "/vendor-name"},
		},
		{
			Paths:  []string{"vendor-name", "app-name"},
			Values: []string{customPrefix + "/vendor-name/app-name"},
		},
	}

	for _, test := range tests {
		if test.Refresh {
			if pathType == "config" {
				os.Setenv("XDG_CONFIG_HOME", test.Env)
			} else {
				os.Setenv("XDG_CACHE_HOME", test.Env)
			}
			configdir.Refresh()
		}

		var result string
		if pathType == "config" {
			result = configdir.LocalConfig(test.Paths...)
		} else {
			result = configdir.LocalCache(test.Paths...)
		}

		if result != test.Values[0] {
			t.Errorf("Got wrong path result. Expected %s, got %s\n",
				test.Values[0],
				result,
			)
		}
	}
}

func TestSystemConfig(t *testing.T) {
	reset()

	// Cases to test.
	var tests = []TestCase{
		{
			Paths:  []string{},
			Values: []string{"/etc/xdg"},
		},
		{
			Paths:  []string{"vendor-name"},
			Values: []string{"/etc/xdg/vendor-name"},
		},
		{
			Paths:  []string{"vendor-name", "app-name"},
			Values: []string{"/etc/xdg/vendor-name/app-name"},
		},

		// With custom XDG paths...
		{
			Env:     "/etc/xdg:/opt/global/conf",
			Refresh: true,
			Paths:   []string{},
			Values:  []string{"/etc/xdg", "/opt/global/conf"},
		},
		{
			Paths:  []string{"vendor-name"},
			Values: []string{"/etc/xdg/vendor-name", "/opt/global/conf/vendor-name"},
		},
		{
			Paths:  []string{"vendor-name", "app-name"},
			Values: []string{"/etc/xdg/vendor-name/app-name", "/opt/global/conf/vendor-name/app-name"},
		},
	}

	for _, test := range tests {
		if test.Refresh {
			os.Setenv("XDG_CONFIG_DIRS", test.Env)
			configdir.Refresh()
		}

		result := configdir.SystemConfig(test.Paths...)

		// Make sure we got the expected result back.
		if len(result) != len(test.Values) {
			t.Errorf("SystemConfig didn't give the expected number of results. "+
				"Expected %d, got %d (env: %s, input paths: %v, result paths: %v)\n",
				len(test.Values),
				len(result),
				test.Env,
				test.Paths,
				result,
			)
			continue
		}

		// Make sure each result is what we expect.
		for i, path := range result {
			if path != test.Values[i] {
				t.Errorf("Got wrong path result on index %d. "+
					"Expected %s, got %s\n",
					i,
					test.Values[i],
					path,
				)
			}
		}
	}
}

func TestLocalConfig(t *testing.T) {
	testLocalCommon(t, "config", "/home/user/.config", "/opt/local")
}

func TestLocalCache(t *testing.T) {
	testLocalCommon(t, "cache", "/home/user/.cache", "/tmp/cache")
}
