package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
)

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

func TestNode(t *testing.T) {
	var err error
	cfg, err = newConfig("xclarity_exporter.yml")
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}
	go webserv()
	url := fmt.Sprintf("%s/chassis", cfg.API.BaseURL)
	client := newBasicAuthClient(cfg.Authentication.Username, cfg.Authentication.Password)
	j := client.getJSON(url, "chassisList")
	nodeParser(j)
}

func TestChassis(t *testing.T) {
	var err error
	cfg, err = newConfig("xclarity_exporter.yml")
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}
	go webserv()
	url := fmt.Sprintf("%s/chassis", cfg.API.BaseURL)
	client := newBasicAuthClient(cfg.Authentication.Username, cfg.Authentication.Password)
	j := client.getJSON(url, "chassisList")
	chassisParser(j)
}
