package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

)


var (



	qryTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "qry_total",
		Help: "The count of Total Queries",
	})

	qryDbTable = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "qry_db_table",
		Help: "The count of queries by db/table",
	},[]string{"db","table"})





)

func startProm() {

	go http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":"+promport, nil)

}