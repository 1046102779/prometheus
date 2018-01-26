package main

import (
	"fmt"
	"math"

	"github.com/golang/protobuf/proto"
	dto "github.com/prometheus/client_model/go"

	"github.com/prometheus/client_golang/prometheus"
)

// score:times
var stats = map[int]int{}

func statitics() {
	temps := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "pond_temperature_celsius",
		Help:    "The temperature of the frog pond.", // Sorry, we can't measure how badly it smells.
		Buckets: prometheus.LinearBuckets(20, 5, 6),  // 5 buckets, each 5 centigrade wide.
	})
	prometheus.MustRegister(temps)

	// Simulate some observations.
	for i := 0; i < 1000; i++ {
		f := float64(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
		temps.Observe(f)
		if f <= 20.0 {
			stats[20]++
		} else if f <= 25.00 {
			stats[25]++
		} else if f <= 30.0 {
			stats[30]++
		} else if f <= 35.0 {
			stats[35]++
		} else if f <= 40.0 {
			stats[40]++
		} else {
			stats[45]++
		}
	}
	metric := &dto.Metric{}
	temps.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))
	for key, value := range stats {
		fmt.Printf("score: %d, times:%d\n", key, value)
	}
}
func main() {
	statitics()
	fmt.Println("end metrics")
}
