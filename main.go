package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

// Config contains the XClarity Exporter configuration
type Config struct {
	API struct {
		Interval time.Duration `yaml:"interval"`
		BaseURL  string        `yaml:"base_url"`
	}
	Authentication struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		CertFile string `yaml:"certfile"`
	}
	Exporter struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	}
}

var (
	cfg            *Config
	flagConfigFile string
)

func newConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	config := &Config{}
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func parseFlags() {
	flag.StringVar(
		&flagConfigFile,
		"config",
		"xclarity_exporter.yml",
		"Path to xclarity_exporter configuration file",
	)
	flag.Parse()
	return
}

// parser is the main loop that endlessly fetches URLs and parses them into
// Prometheus metrics
func parser(loop bool) {
	chassisURL := fmt.Sprintf("%s/chassis", cfg.API.BaseURL)
	nodeURL := fmt.Sprintf("%s/nodes", cfg.API.BaseURL)
	client := newBasicAuthClient(cfg.Authentication.Username, cfg.Authentication.Password)
	for {
		j, err := client.getJSON(chassisURL, "chassisList")
		// Failing to fetch a URL shouldn't be fatal.  Skip this parsing cycle and carry on.
		if err != nil {
			log.Printf("Parsing %s returned: %v", chassisURL, err)
		} else {
			chassisParser(j)
		}
		j, err = client.getJSON(nodeURL, "nodeList")
		if err != nil {
			log.Printf("Parsing %s returned: %v", nodeURL, err)
		} else {
			nodeParser(j)
		}
		// Bail out if not configured to loop.  This is useful for testing.
		if !loop {
			break
		}
		// Sleep for a configured duration
		time.Sleep(cfg.API.Interval * time.Second)
	} // Endless loop
}

func main() {
	var err error
	parseFlags()
	cfg, err = newConfig(flagConfigFile)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}
	// If loop is true, the parser will loop endlessly
	loop := true
	go parser(loop)
	http.Handle("/metrics", promhttp.Handler())
	exporter := fmt.Sprintf("%s:%s", cfg.Exporter.Address, cfg.Exporter.Port)
	http.ListenAndServe(exporter, nil)
}
