package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func parseFile(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

var health = map[string]int{
	"Normal":        0,
	"Warning":       1,
	"Critical":      2,
	"Major-Failure": 3,
}

func nodeParser() {
	//j := gjson.ParseBytes(parseFile("full.json"))
	j := gjson.GetBytes(parseFile("full.json"), "nodeList")
	for _, jn := range j.Array() {
		instance := strings.ToLower(jn.Get("name").String())
		healthCode, ok := health[jn.Get("overallHealthState").String()]
		if ok {
			nodeHealth.WithLabelValues(instance).Set(float64(healthCode))
		}
		nodePower.WithLabelValues(instance).Set(jn.Get("powerStatus").Float())
	}
}

func testParser() {
	j := gjson.GetBytes(parseFile("chassis_clean.json"), "chassisList.#.userDefinedName")
	for _, jn := range j.Array() {
		fmt.Println(jn.String())
	}
}

func chassisParser() {
	j := gjson.GetBytes(parseFile("chassis_clean.json"), "chassisList")
	for _, jn := range j.Array() {
		instance := strings.ToLower(jn.Get("userDefinedName").String())
		chassisPowerFree.WithLabelValues(instance).Set(jn.Get("powerAllocation.remainingOutputPower").Float())
		chassisPowerTotal.WithLabelValues(instance).Set(jn.Get("powerAllocation.totalOutputPower").Float())
		chassisPowerUsed.WithLabelValues(instance).Set(jn.Get("powerAllocation.allocatedOutputPower").Float())
		for _, sw := range jn.Get("switches").Array() {
			if sw.Get("type").String() != "Switch" {
				fmt.Println(sw.Get("type").String())
				continue
			}
			switchName := strings.ToLower(sw.Get("deviceName").String())
			healthCode, ok := health[sw.Get("overallHealthState").String()]
			if ok {
				chassisSwitchHealth.WithLabelValues(instance, switchName).Set(float64(healthCode))
			}
		}
		for n, ps := range jn.Get("powerSupplies").Array() {
			if ps.Get("type").String() != "PowerSupply" {
				fmt.Println(ps.Get("type").String())
				continue
			}
			healthCode, ok := health[ps.Get("excludedHealthState").String()]
			if ok {
				chassisPSUHealth.WithLabelValues(instance, strconv.Itoa(n)).Set(float64(healthCode))
			}
		}
	}
}
