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
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// OpenhabStatsCollector implements the Collector interface
type OpenhabStatsCollector struct {
	logger log.Logger
}

// Descriptors
var (
	metricDescriptions = prometheus.NewDesc(
		"openhab_item_state_current",
		"Openhab items current state",
		[]string{"item", "label", "type", "tags", "groupnames"}, nil,
	)
)

// Describe implements DescribeByCollect
func (cc OpenhabStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

// Collect collects data from items
func (cc OpenhabStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats, err := getRestItems()

	if err != nil {
		level.Error(cc.logger).Log("msg", "Error getting items", "err", err)
	}

	for _, item := range stats {
		process := true
		switch item.Type {
		case "Number":
		case "Dimmer":
		case "Switch":
		case "Contact":
		default:
			process = false
		}
		if process {
			ch <- prometheus.MustNewConstMetric(
				metricDescriptions,
				prometheus.GaugeValue,
				func() float64 {
					f := 0.0
					switch item.State {
					case "CLOSED":
						f = 0
					case "OFF":
						f = 0
					case "ON":
						f = 1
					case "OPEN":
						f = 1
					case "NULL":
						break
					default:
						f, _ = strconv.ParseFloat(item.State, 64)
					}
					return f
				}(),
				item.Name,
				item.Label,
				item.Type,
				strings.Join(item.Tags, ";"),
				strings.Join(item.GroupNames, ";"),
			)
		} else {
			level.Debug(cc.logger).Log("msg", "Skipped item", "name", item.Name, "type", item.Type, "state", item.State)
		}
	}
}

func handleCollector(logger log.Logger) {
	reg := prometheus.NewPedanticRegistry()

	cc := OpenhabStatsCollector{logger: logger}
	reg.MustRegister(cc)

	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
}
