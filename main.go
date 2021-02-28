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

func parser() {
	chassisURL := fmt.Sprintf("%s/chassis", cfg.API.BaseURL)
	nodeURL := fmt.Sprintf("%s/nodes", cfg.API.BaseURL)
	client := newBasicAuthClient(cfg.Authentication.Username, cfg.Authentication.Password)
	for {
		chassisParser(client.getJSON(chassisURL, "chassisList"))
		nodeParser(client.getJSON(nodeURL, "nodeList"))
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
	go parser()
	http.Handle("/metrics", promhttp.Handler())
	exporter := fmt.Sprintf("%s:%s", cfg.Exporter.Address, cfg.Exporter.Port)
	http.ListenAndServe(exporter, nil)
}
