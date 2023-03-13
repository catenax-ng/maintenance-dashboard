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

// Sync current versions and latest versions periodically
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
		log.Infof("Sync finished in %v seconds.", elapsed.Seconds())

		refreshSeconds, _ := strconv.ParseFloat(refreshIntervalInSeconds, 64)
		time.Sleep(time.Duration(refreshSeconds * float64(time.Second)))
	}
}

func main() {
	// Set logging format to json
	log.SetFormatter(&log.JSONFormatter{})
	log.Infoln("App startup ongoing.")

	go syncAppsVersionInfo()

	prometheusHandler := metrics.CreateMetricsHandler()
	// Setup metrics endpoint and start server
	http.Handle("/metrics", prometheusHandler)

	// Create health http handler
	h, _ := health.New()
	http.Handle("/health", h.Handler())

	// Start webserver on port 2112
	port := ":2112"
	log.Infof("Listening on %v", port)
	http.ListenAndServe(port, nil)
}
