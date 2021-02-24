package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Node definitions
	nodeHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xclarity_server_health",
			Help: "XClarity server health (0=Good)",
		},
		[]string{"instance"},
	)
	nodePower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xclarity_server_power",
			Help: "XClarity server power status",
		},
		[]string{"instance"},
	)
	// Chassis definitions
	chassisPowerFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xclarity_chassis_power_free",
			Help: "Chassis power unused",
		},
		[]string{"instance"},
	)
	chassisPowerTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xclarity_chassis_power_total",
			Help: "Chassis power total",
		},
		[]string{"instance"},
	)
	chassisPowerUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xclarity_chassis_power_used",
			Help: "Chassis power consumed",
		},
		[]string{"instance"},
	)
	chassisSwitchHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xclarity_chassis_switch_health",
			Help: "Switch overall health (0=Good)",
		},
		[]string{"instance", "switch"},
	)
	chassisPSUHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xclarity_chassis_psu_health",
			Help: "PSU health (0=Good)",
		},
		[]string{"instance", "psu"},
	)
)

func init() {
	// node metrics
	prometheus.MustRegister(nodeHealth)
	prometheus.MustRegister(nodePower)
	// chassis metrics
	prometheus.MustRegister(chassisPowerFree)
	prometheus.MustRegister(chassisPowerTotal)
	prometheus.MustRegister(chassisPowerUsed)
	prometheus.MustRegister(chassisSwitchHealth)
	prometheus.MustRegister(chassisPSUHealth)
}
