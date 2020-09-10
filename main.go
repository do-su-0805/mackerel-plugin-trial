package main

import (
	"flag"
	"fmt"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
	"github.com/mackerelio/go-osstat/uptime"
)

type UptimePlugin struct {
	Prefix string
}

type PluginWithPrefix interface {
	FetchMetrics() (map[string]float64, error)
	GraphDefinition() map[string]mp.Graphs
	MetricKeyPrefix() string
}

func (u UptimePlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(u.MetricKeyPrefix())
	return map[string]mp.Graphs{
		"": {
			Label: labelPrefix,
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "seconds", Label: "Seconds"},
			},
		},
	}
}

func (u UptimePlugin) FetchMetrics() (map[string]float64, error) {
	ut, err := uptime.Get()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch uptime metrics: %s", err)
	}
	return map[string]float64{"seconds": ut.Seconds()}, nil
}

func (u UptimePlugin) MetricKeyPrefix() string {
	if u.Prefix == "" {
		u.Prefix = "uptime"
	}
	return u.Prefix
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "uptime", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	u := UptimePlugin{
		Prefix: *optPrefix,
	}
	plugin := mp.NewMackerelPlugin(u)
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
