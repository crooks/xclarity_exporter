package main

import (
	"log"
	"net/http"
	"strings"
	"testing"
)

// webserv starts a HTTP listener to test API scraping
func webserv() {
	fs := http.FileServer(http.Dir("./testdata"))
	http.Handle("/", fs)
	// listenHost is a hacky way to get the host:port from the API URL config
	listenHost := strings.Split(cfg.API.BaseURL, "://")[1]
	log.Printf("Opening webserver on %s", listenHost)
	err := http.ListenAndServe(listenHost, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// strEquality compares two strings and fails if they don't match
func strEquality(t *testing.T, wanted, got string) {
	if wanted != got {
		t.Fatalf("Failed string equality test.  Wanted=%s, Got=%s", wanted, got)
	}
}

// intEquality compares two integers and fails if they don't match
func intEquality(t *testing.T, wanted, got int) {
	if wanted != got {
		t.Fatalf("Failed integer equality test.  Wanted=%d, Got=%d", wanted, got)
	}
}

// TestCfg checks the values in the default config file are corrrect.  As
// they're used in other tests, the defaults can't be messed with.
func TestCfg(t *testing.T) {
	var err error
	cfg, err := newConfig("xclarity_exporter.yml")
	if err != nil {
		t.Fatalf("Unable to parse config file: %v", err)
	}
	strEquality(t, "http://127.0.0.1:8080", cfg.API.BaseURL)
	// cfg.API.Interval is a time.Duration but we'll pretend it's an integer
	//for testing equality.
	intEquality(t, 120, int(cfg.API.Interval))
	strEquality(t, "apiauthuser", cfg.Authentication.Username)
	strEquality(t, "apiauthpassword", cfg.Authentication.Password)
	strEquality(t, "/etc/pki/tls/certs/rootcrt.pem", cfg.Authentication.CertFile)
	strEquality(t, "127.0.0.1", cfg.Exporter.Address)
	strEquality(t, "9794", cfg.Exporter.Port)
}

func TestParser(t *testing.T) {
	var err error
	cfg, err = newConfig("xclarity_exporter.yml")
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}
	go webserv()
	parser(false)
}
