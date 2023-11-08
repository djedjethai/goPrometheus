package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	// "go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
)

const meterName = "github.com/open-telemetry/opentelemetry-go/example/prometheus"

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ctx := context.Background()

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter(meterName)

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics()

	opt := api.WithAttributes(
		attribute.Key("A").String("B"),
		attribute.Key("C").String("D"),
	)

	// This is the equivalent of prometheus.NewCounterVec
	counter, err := meter.Float64Counter("foo", api.WithDescription("a simple counter"))
	if err != nil {
		log.Fatal(err)
	}
	counter.Add(ctx, 5, opt)

	gauge, err := meter.Float64ObservableGauge("bar", api.WithDescription("a fun little gauge"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
		n := -10. + rng.Float64()*(90.) // [-10, 100)
		o.ObserveFloat64(gauge, n, opt)
		return nil
	}, gauge)
	if err != nil {
		log.Fatal(err)
	}

	// // This is the equivalent of prometheus.NewHistogramVec
	// histogram, err := meter.Float64Histogram(
	// 	"baz",
	// 	api.WithDescription("a histogram with custom buckets and rename"),
	// 	api.WithExplicitBucketBoundaries(64, 128, 256, 512, 1024, 2048, 4096),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// histogram.Record(ctx, 136, opt)
	// histogram.Record(ctx, 64, opt)
	// histogram.Record(ctx, 701, opt)
	// histogram.Record(ctx, 830, opt)

	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	<-ctx.Done()
}

func serveMetrics() {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2223", nil) //nolint:gosec // Ignoring G114: Use of net/http serve function that has no support for setting timeouts.
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}

// var meterName = "MyServiceName"
//
// func main() {
// 	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	ctx := context.Background()
//
// 	// The exporter embeds a default OpenTelemetry Reader and
// 	// implements prometheus.Collector, allowing it to be used as
// 	// both a Reader and Collector.
// 	exporter, err := prometheus.New()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// provider := metric.NewMeterProvider(metric.WithReader(exporter))
// 	provider := metric.NewMeterProvider(
// 		metric.WithReader(exporter),
// 		// View to customize histogram buckets and rename a single histogram instrument.
// 		metric.WithView(metric.NewView(
// 			metric.Instrument{
// 				Name:  "baz",
// 				Scope: instrumentation.Scope{Name: meterName},
// 			},
// 			metric.Stream{
// 				Name: "new_baz",
// 				Aggregation: metric.AggregationExplicitBucketHistogram{
// 					Boundaries: []float64{64, 128, 256, 512, 1024, 2048, 4096},
// 				},
// 			},
// 		)),
// 	)
// 	meter := provider.Meter(meterName)
//
// 	// Start the prometheus HTTP server and pass the exporter Collector to it
// 	go serveMetrics()
//
// 	opt := api.WithAttributes(
// 		attribute.Key("A").String("B"),
// 		attribute.Key("C").String("D"),
// 	)
//
// 	// // This is the equivalent of prometheus.NewCounterVec
// 	// counter, err := meter.Float64Counter("foo", api.WithDescription("a simple counter"))
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// counter.Add(ctx, 5, opt)
//
// 	// fail !!!
// 	// Create a new counter metric for the "/" endpoint
// 	// hitsCounter, err := meter.Float64Counter("hits", api.WithDescription("Total hits on /"))
// 	hitsCounter, err := meter.Float64Counter("hits", api.WithDescription("Total hits on /"))
// 	if err != nil {
// 		log.Fatal(err)
// 	} else {
// 		fmt.Println("HitsCounter run.")
// 	}
// 	hitsCounter.Add(ctx, 1, opt)
//
// 	// works !!!
// 	gauge, err := meter.Float64ObservableGauge("bar", api.WithDescription("a fun little gauge"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
// 		n := -10. + rng.Float64()*(90.) // [-10, 100)
// 		o.ObserveFloat64(gauge, n, opt)
// 		return nil
// 	}, gauge)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// // This is the equivalent of prometheus.NewHistogramVec
// 	histogram, err := meter.Float64Histogram("baz", api.WithDescription("a very nice histogram"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	histogram.Record(ctx, 23, opt)
// 	histogram.Record(ctx, 7, opt)
// 	histogram.Record(ctx, 101, opt)
// 	histogram.Record(ctx, 105, opt)
//
// 	// in case no server, keep the program running(common pattern...)
// 	// ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
// 	// <-ctx.Done()
//
// 	var count = 0
//
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		// Increment the hitsCounter metric when the "/" endpoint is accessed
//
// 		// hitsCounter.Add(ctx, 1, opt)
// 		hitsCounter.Add(ctx, 1.0, opt)
//
// 		// this err get any error..??
// 		// if err != nil {
// 		// 	log.Printf("Error incrementing hitsCounter: %v", err)
// 		// }
//
// 		// Write a response to the client.
// 		fmt.Fprintf(w, fmt.Sprintf("Hello, World! - %v\n", count))
// 		count++
// 	})
//
// 	fmt.Println("run on port 8000")
// 	err = http.ListenAndServe(":8000", nil)
// 	if err != nil {
// 		fmt.Printf("error serving http: %v", err)
// 		return
// 	}
// }
//
// func serveMetrics() {
// 	log.Printf("serving metrics at localhost:3000/metrics")
//
// 	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	// 	// Write a response to the client.
// 	// 	fmt.Fprintf(w, "Hello, World!\n")
// 	// })
// 	http.Handle("/metrics", promhttp.Handler())
//
// 	err := http.ListenAndServe(":3000", nil)
// 	if err != nil {
// 		fmt.Printf("error serving http: %v", err)
// 		return
// 	}
// }
