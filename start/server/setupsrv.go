package main

// import (
// 	"context"
// 	"fmt"
// 	"go.opentelemetry.io/otel/exporters/metric/prometheus"
// 	// "go.opentelemetry.io/otel/exporters/stdout"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/label"
// 	"go.opentelemetry.io/otel/metric"
//
// 	"log"
// 	"net/http"
// 	"runtime"
// )
//
// var (
// 	requests    metric.Int64Counter
// 	appName     string = "prometheus"
// 	serviceName string = "myservice"
// )
//
// var labels = []label.KeyValue{
// 	label.Key("application").String(appName),
// 	label.Key("service").String(serviceName),
// 	label.Key("container_id").String("1234"),
// }
//
// func configPrometheus() {
// 	prometheusExporter, err := prometheus.NewExportPipeline(prometheus.Config{})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	// Get the meter provider from the exporter.
// 	mp := prometheusExporter.MeterProvider()
//
// 	// Set it as the global meter provider.
// 	otel.SetMeterProvider(mp)
//
// 	// // Register the exporter as the handler for the "/metrics" pattern.
// 	// http.Handle("/metrics", prometheusExporter)
// 	// // Start the HTTP server listening on port 3000.
// 	// log.Fatal(http.ListenAndServe(":3000", nil))
// 	go runPrometheusEndPoint(prometheusExporter)
//
// 	// meter := otel.GetMeterProvider().Meter("golru")
//
// 	err = buildRequestsCounter()
// 	if err != nil {
// 		log.Println("Error from build request counter: ", err)
// 	}
//
// 	buildRuntimeObservers()
// }
//
// func runPrometheusEndPoint(prometheusExporter *prometheus.Exporter) {
// 	// Register the exporter as the handler for the "/metrics" pattern.
// 	http.Handle("/metrics", prometheusExporter)
// 	// Start the HTTP server listening on port 3000.
// 	log.Fatal(http.ListenAndServe(":3000", nil))
// }
//
// func buildRequestsCounter() error {
// 	var err error
// 	// Retrieve the meter from the meter provider.
// 	meter := otel.GetMeterProvider().Meter(serviceName)
// 	// Get an Int64Counter for a metric called "prometheus_requests_total".
// 	requests, err = meter.NewInt64Counter("prometheus_requests_total",
// 		metric.WithDescription("Total number of golru requests."),
// 	)
// 	return err
// }
//
// // the NewInt64UpDownSumObserver accepts the name of the metric as a
// // string, something called a Int64ObserverFunc, and zero or more instrument
// // options (such as the metric description)
// func buildRuntimeObservers() {
// 	meter := otel.GetMeterProvider().Meter(serviceName)
// 	m := runtime.MemStats{}
// 	meter.NewInt64UpDownSumObserver("memory_usage_bytes",
// 		func(_ context.Context, result metric.Int64ObserverResult) {
// 			runtime.ReadMemStats(&m)
// 			result.Observe(int64(m.Sys), labels...)
// 		},
// 		metric.WithDescription("Amount of memory used."),
// 	)
// 	meter.NewInt64UpDownSumObserver("num_goroutines",
// 		func(_ context.Context, result metric.Int64ObserverResult) {
// 			result.Observe(int64(runtime.NumGoroutine()), labels...)
// 		},
// 		metric.WithDescription("Number of running goroutines."),
// 	)
// }
