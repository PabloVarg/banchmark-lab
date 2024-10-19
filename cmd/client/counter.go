package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var serverHits = promauto.NewCounter(prometheus.CounterOpts{
	Name: "client_server_hits",
	Help: "Number of times the server has been called",
})

var serverHitsFailed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "client_server_hits_fails",
	Help: "Number of times the server has been called and received an error",
})

var requestLatency = promauto.NewHistogram(prometheus.HistogramOpts{
	Name: "latency_on_requests",
	Help: "Experienced latency while calling the server",
})
