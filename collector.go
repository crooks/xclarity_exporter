package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// health is a mapping of textual XClarity health states to a numeric metric
var health = map[string]int{
	"Unknown":         0,
	"Normal":          1,
	"Non-Critical":    2,
	"Warning":         3,
	"Minor-Failure":   4,
	"Major-Failure":   5,
	"Non-Recoverable": 6,
	"Critical":        7,
}

// getJSON expects a URL and a top-level json dict name to scrape.  It returns
// that dict name as a gjson object.
func (s *authClient) getJSON(url, tlj string) (gjson.Result, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("request error: %v", err)
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("node request error: %v", err)
	}
	return gjson.GetBytes(bytes, tlj), nil
}

// nodeParser parses the json output from the XClarity API (https://<xclarity_server>/nodes)
func nodeParser(j gjson.Result) {
	for _, jn := range j.Array() {
		node := strings.ToLower(jn.Get("name").String())
		healthCode, ok := health[jn.Get("overallHealthState").String()]
		if ok {
			nodeHealth.WithLabelValues(node).Set(float64(healthCode))
		}
		nodePower.WithLabelValues(node).Set(jn.Get("powerStatus").Float())
	}
}

// chassisSwitchParser parses a list of switches associated with a specific
// chassis instance
func chassisSwitchParser(j gjson.Result, instance string) {
	for _, sw := range j.Array() {
		swType := sw.Get("type").String()
		if swType != "Switch" {
			log.Printf("Unexpected Switch type: %s", swType)
			continue
		}
		switchName := strings.ToLower(sw.Get("deviceName").String())
		healthCode, ok := health[sw.Get("overallHealthState").String()]
		if ok {
			chassisSwitchHealth.WithLabelValues(instance, switchName).Set(float64(healthCode))
		}
	}
}

// chassisPSUParser parses a list of PSUs associated with a specific
// chassis instance.  PSUs don't have a name so this function infers a name
// from the list item number.
func chassisPSUParser(j gjson.Result, instance string) {
	for n, ps := range j.Array() {
		psType := ps.Get("type").String()
		if psType != "PowerSupply" {
			log.Printf("Unexpected Power Supply type: %s", psType)
			continue
		}
		healthCode, ok := health[ps.Get("overallHealthState").String()]
		if ok {
			chassisPSUHealth.WithLabelValues(instance, strconv.Itoa(n)).Set(float64(healthCode))
		}
	}
}

// chassisParser parses the json output from the XClarity API (https://<xclarity_server>/chassis)
func chassisParser(j gjson.Result) {
	// Iterate through the list of Flex chassis
	for _, jn := range j.Array() {
		// The user-defined chassis name is used to populate the instance label
		// of all metrics associated with this list item.
		chassis := strings.ToLower(jn.Get("userDefinedName").String())
		// cmmHealthState occurs all over the place.  This is the top-level
		// chassis CMM health.
		cmmHealth, ok := health[jn.Get("cmmHealthState").String()]
		if ok {
			chassisCMMHealth.WithLabelValues(chassis).Set(float64(cmmHealth))
		}
		// Power resources are defined at the top-level of each list item
		chassisPowerFree.WithLabelValues(chassis).Set(jn.Get("powerAllocation.remainingOutputPower").Float())
		chassisPowerTotal.WithLabelValues(chassis).Set(jn.Get("powerAllocation.totalOutputPower").Float())
		chassisPowerUsed.WithLabelValues(chassis).Set(jn.Get("powerAllocation.allocatedOutputPower").Float())
		chassisPowerMode.WithLabelValues(chassis).Set(jn.Get("energyPolicies.powerRedundancyMode").Float())
		// switches and PSUs are lists within each chassis instance
		chassisSwitchParser(jn.Get("switches"), chassis)
		chassisPSUParser(jn.Get("powerSupplies"), chassis)
	}
}
