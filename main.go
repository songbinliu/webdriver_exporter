// Copyright 2016 The Prometheus Authors
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
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/sclevine/agouti"
)

const versionString = "0.0.2"

var (
	driver        = agouti.ChromeDriver()
	listenIP      = flag.String("listenIP", "127.0.0.1", "the ip address to listen on for HTTP requests.")
	listenPort    = flag.Int("port", 9156, "the port to listen on for HTTP requests.")
	showVersion   = flag.Bool("version", false, "Print version information.")
)

var (
	counter = 0
	logFlag = false
)

func init() {
	version.Version = versionString
	prometheus.MustRegister(version.NewCollector("webdriver_exporter"))
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("webdriver_exporter"))
		os.Exit(0)
	}

	log.Infoln("Starting webdriver_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/probe",
		func(w http.ResponseWriter, r *http.Request) {
			probeHandler(w, r)
		})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>WebDriver Exporter</title></head>
            <body>
            <h1>WebDriver Exporter</h1>
	    <p><a href="/probe?target=https://prometheus.io/">Probe prometheus.io</a></p>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})

	log.Infoln("Starting webdriver")
	err := driver.Start()
	if err != nil {
		log.Fatalf("failed to start webdriver: %s", err)
	}
	defer driver.Stop()

	stop := make(chan struct{})
	defer close(stop)
	go guardLog(stop)

	address := fmt.Sprintf("%v:%d", *listenIP, *listenPort)
	log.Infoln("Listening on", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}
}

func guardLog(stop <- chan struct{}) {
	for {

		timer := time.NewTimer(time.Second * 10)
		select {
		case <- stop:
			return
		case <-timer.C:
		}

		counter ++
		logFlag = false
		if counter % 3 == 0 {
			logFlag = true
		}

		if counter > 1e10 {
			counter = 0
		}
	}
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", 400)
		return
	}

	if logFlag {
		log.Infof("begin to probe target: %v", target)
	}

	start := time.Now()
	success := probe(target, w)
	fmt.Fprintf(w, "probe_duration_seconds %f\n", float64(time.Since(start))/1e9)

	if logFlag {
		log.Infof("end of probing target: %v", target)
	}

	if success {
		fmt.Fprintf(w, "probe_success %d\n", 1)
	} else {
		fmt.Fprintf(w, "probe_success %d\n", 0)
	}
}
