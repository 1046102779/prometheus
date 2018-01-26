package main

import (
	"fmt"
	"math"

	"github.com/golang/protobuf/proto"

	"github.com/prometheus/client_golang/prometheus"
)

func statistic() {
	temps := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "pond_temperature_celsius",
			Help:       "The temperature of the frog pond.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"species"},
	)

	// Simulate some observations.
	for i := 0; i < 1000; i++ {
		temps.WithLabelValues("litoria-caerulea").Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
		temps.WithLabelValues("lithobates-catesbeianus").Observe(32 + math.Floor(100*math.Cos(float64(i)*0.11))/10)
	}

	// Create a Summary without any observations.
	temps.WithLabelValues("leiopelma-hochstetteri")

	// Just for demonstration, let's check the state of the summary vector
	// by registering it with a custom registry and then let it collect the
	// metrics.
	reg := prometheus.NewRegistry()
	reg.MustRegister(temps)

	metricFamilies, err := reg.Gather()
	if err != nil || len(metricFamilies) != 1 {
		panic("unexpected behavior of custom test registry")
	}
	fmt.Println(proto.MarshalTextString(metricFamilies[0]))
}

func main() {
	statistic()
	fmt.Println("end metrics")
}
