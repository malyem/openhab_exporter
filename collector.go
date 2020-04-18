// Copyright 2014 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// OpenhabStats is an example for a system that might have been built without
// Prometheus in mind. It models a central manager of jobs running in a
// cluster. Thus, we implement a custom Collector called
// OpenhabStatsCollector, which collects information from a OpenhabStats
// using its provided methods and turns them into Prometheus Metrics for
// collection.
//
// An additional challenge is that multiple instances of the OpenhabStats are
// run within the same binary, each in charge of a different zone. We need to
// make use of wrapping Registerers to be able to register each
// OpenhabStatsCollector instance with Prometheus.
type OpenhabStats struct {
	Zone string
	// Contains many more fields not listed in this example.
}

// ReallyExpensiveAssessmentOfTheSystemState is a mock for the data gathering a
// real cluster manager would have to do. Since it may actually be really
// expensive, it must only be called once per collection. This implementation,
// obviously, only returns some made-up data.
func (c *OpenhabStats) ReallyExpensiveAssessmentOfTheSystemState() []Item {
	return doAuthRequest("http://192.168.0.116:8080/rest/items")

}

// OpenhabStatsCollector implements the Collector interface.
type OpenhabStatsCollector struct {
	OpenhabStats *OpenhabStats
}

// Descriptors used by the OpenhabStatsCollector below.
var (
	oomCountDesc = prometheus.NewDesc(
		"openhab_item_state_current",
		"Openhab items current state",
		[]string{"item", "label", "type", "tags", "groupnames"}, nil,
	)
)

// Describe is implemented with DescribeByCollect. That's possible because the
// Collect method will always return the same two metrics with the same two
// descriptors.
func (cc OpenhabStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

// Collect first triggers the ReallyExpensiveAssessmentOfTheSystemState. Then it
// creates constant metrics for each host on the fly based on the returned data.
//
// Note that Collect could be called concurrently, so we depend on
// ReallyExpensiveAssessmentOfTheSystemState to be concurrency-safe.
func (cc OpenhabStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := cc.OpenhabStats.ReallyExpensiveAssessmentOfTheSystemState()

	// rand.Seed(time.Now().UnixNano())

	for _, item := range stats {
		ch <- prometheus.MustNewConstMetric(
			oomCountDesc,
			prometheus.GaugeValue,
			func() float64 {
				f, _ := strconv.ParseFloat(item.State, 64)
				return f
			}(),
			item.Name,
			item.Label,
			item.Type,
			strings.Join(item.Tags, ";"),
			strings.Join(item.GroupNames, ";"),
		)
	}

}

// NewOpenhabStats first creates a Prometheus-ignorant OpenhabStats
// instance. Then, it creates a OpenhabStatsCollector for the just created
// OpenhabStats. Finally, it registers the OpenhabStatsCollector with a
// wrapping Registerer that adds the zone as a label. In this way, the metrics
// collected by different OpenhabStatsCollectors do not collide.
func NewOpenhabStats(zone string, reg prometheus.Registerer) *OpenhabStats {
	c := &OpenhabStats{
		Zone: zone,
	}
	cc := OpenhabStatsCollector{OpenhabStats: c}
	prometheus.WrapRegistererWith(prometheus.Labels{"zone": zone}, reg).MustRegister(cc)
	return c
}

//HandleCollector HandleCollector
func HandleCollector() {
	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()

	// Construct cluster managers. In real code, we would assign them to
	// variables to then do something with them.
	NewOpenhabStats("db", reg)

	// Add the standard process and Go metrics to the custom registry.
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
