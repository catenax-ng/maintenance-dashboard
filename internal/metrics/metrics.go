package metrics

import (
	"net/http"

	"github.com/catenax-ng/maintenance-dashboard/internal/data"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var registry *prometheus.Registry
var metrics = map[string]prometheus.Gauge{}

// Add metrics http handler
func CreateMetricsHandler() http.Handler {
	// Get rid of the default metrics
	registry = prometheus.NewRegistry()
	// Add the metrics for all applications with Gauge type
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	return handler
}

// Generate the metrics for every apps current and latest major, minor and patch versions
func CreateOrUpdateMetrics(appsVersionInfo []*data.AppVersionInfo) {
	for _, appVersionInfo := range appsVersionInfo {
		CreateOrUpdateMetricSingle(appVersionInfo.NewReleasesName, appVersionInfo.ResourceName, "major", "current", float64(appVersionInfo.CurrentVersion.Major()))
		CreateOrUpdateMetricSingle(appVersionInfo.NewReleasesName, appVersionInfo.ResourceName, "minor", "current", float64(appVersionInfo.CurrentVersion.Minor()))
		CreateOrUpdateMetricSingle(appVersionInfo.NewReleasesName, appVersionInfo.ResourceName, "patch", "current", float64(appVersionInfo.CurrentVersion.Patch()))
		CreateOrUpdateMetricSingle(appVersionInfo.NewReleasesName, appVersionInfo.ResourceName, "major", "latest", float64(appVersionInfo.LatestMajorVersion.Major()))
		CreateOrUpdateMetricSingle(appVersionInfo.NewReleasesName, appVersionInfo.ResourceName, "minor", "latest", float64(appVersionInfo.LatestMajorVersion.Minor()))
		CreateOrUpdateMetricSingle(appVersionInfo.NewReleasesName, appVersionInfo.ResourceName, "patch", "latest", float64(appVersionInfo.LatestMajorVersion.Patch()))
	}

	log.Infoln("Metrics created and updated.")
}

// Create a new metric is it doesn't exist then set it's value
func CreateOrUpdateMetricSingle(releaseName string, resourceName string, versionPart string, origin string, value float64) {
	mapKey := resourceName + "-" + origin + "-" + versionPart
	metric := metrics[mapKey]
	if metric == nil {
		metrics[mapKey] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "maintenance",
			Name:      "app_version_info",
			Help:      "Version information about an application",
			ConstLabels: prometheus.Labels{
				"release_name":  releaseName,
				"resource_name": resourceName,
				"version_part":  versionPart,
				"origin":        origin,
			},
		})
		metric = metrics[mapKey]
		registry.MustRegister(metric)
		log.Infof("app_version_info metric created and registered for %v %v version part of app %v.", origin, versionPart, resourceName)
	}
	metric.Set(value)
	log.Infof("app_version_info metric %v value set: %v", mapKey, value)
}
