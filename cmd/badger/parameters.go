package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/zonedb/zonedb"
)

// Configuration ...
type Configuration struct {
	AllCharacters  bool
	Alpha          bool
	AlphaNumeric   bool
	CustomRange    string
	TLD            []string
	SearchPatterns []string
}

func validateDomain(domain string) bool {
	// This only validates the domain name, not the tld
	// Checks for alphanumeric characters including dash (-)
	// no spaces, min length of 1 and a max length of 63

	// RegEx explanation:
	// - 1st char only alphanumeric
	// - Subsequent chars (min 0, max 61) alphanumeric and dash
	// - Last char only alphanumeric
	r := regexp.MustCompile("(^[a-z0-9])([a-z0-9-]{0,61})(?:[a-z0-9])$")
	return r.MatchString(domain)
}

func validateArguments() error {
	config := getConfigurationFromArguments()

	if config.AllCharacters {
		config.Alpha = false
		config.AlphaNumeric = false
	}

	if config.AlphaNumeric == true {
		config.Alpha = false
	}

	if len(config.CustomRange) > 0 {
		config.Alpha = false
		config.AlphaNumeric = false
		config.AllCharacters = false
	}

	// TODO(ea): check custom range for invalid characters

	// Make sure TLDs have a corresponding nameserver
	for _, tld := range config.TLD {
		zone := zonedb.PublicZone(tld)
		if zone == nil {
			return fmt.Errorf("Invalid TLD specified: %v", tld)
		}
	}

	// TODO(ea): check searchpatterns for wildcard character (underscore)
	// TODO(ea): check searchpatterns for invalid characters

	return nil
}

func getConfigurationFromArguments() Configuration {

	config := Configuration{}

	flag.BoolVar(&config.AllCharacters, "all", true, "Use all possible characters (a-z, 0-9, -)")
	flag.BoolVar(&config.Alpha, "alpha", false, "Use alphabetic range (a-z)")
	flag.BoolVar(&config.AlphaNumeric, "alphanum", false, "Use alphanumeric range (a-z, 0-9)")
	flag.StringVar(&config.CustomRange, "custom", "", "Use a custom character range (ex. abc123)")

	var tlds string
	flag.StringVar(&tlds, "tld", "com", "TLDs to search. Use comma to add multiple (ex. com,org,net)")

	flag.Parse()

	config.TLD = strings.Split(tlds, ",")

	config.SearchPatterns = flag.Args() // Search mask to use (ex. 'se_rchm_sk' to use 2 wildcard ranges)

	return config
}
