package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	// "encoding/json"

	"github.com/catenax-ng/maintenance-dashboard/internal/currentversions"
	"github.com/catenax-ng/maintenance-dashboard/internal/data"
	"github.com/catenax-ng/maintenance-dashboard/internal/helpers"
	"github.com/catenax-ng/maintenance-dashboard/internal/latestversions"
	"github.com/catenax-ng/maintenance-dashboard/internal/metrics"
	health "github.com/hellofresh/health-go/v5"
	log "github.com/sirupsen/logrus"
)

var refreshIntervalInSeconds = helpers.GetEnv("REFRESH_INTERVAL_SECONDS", "60")

func syncAppsVersionInfo() {
	for {
		log.Infoln("New sync started.")
		start := time.Now()
		ctx := context.Background()
		var appsVersionInfo []*data.AppVersionInfo
		appCurrentInfos := currentversions.GetCurrentVersions(ctx)

		for _, appCurrentInfo := range appCurrentInfos {
			appVersionInfo := latestversions.GetForApp(*appCurrentInfo)
			appsVersionInfo = append(appsVersionInfo, appVersionInfo)
		}
		metrics.UpdateMetrics(appsVersionInfo)

		elapsed := time.Since(start)
		log.Infof("Sync finished in %v seconds.\n", elapsed.Seconds())

		refreshSeconds, _ := strconv.ParseFloat(refreshIntervalInSeconds, 64)
		time.Sleep(time.Duration(refreshSeconds * float64(time.Second)))
	}
}

func main() {
	log.Infoln("App startup ongoing.")
	log.SetFormatter(&log.JSONFormatter{})

	go syncAppsVersionInfo()

	prometheusHandler := metrics.CreateMetrics()
	// setup metrics endpoint and start server
	http.Handle("/metrics", prometheusHandler)

	h, _ := health.New()
	http.Handle("/health", h.Handler())

	port := ":2112"
	log.Infof("Listening on %v", port)
	http.ListenAndServe(port, nil)
}
