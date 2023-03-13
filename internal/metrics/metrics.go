package metrics

import (
	"net/http"

	"github.com/catenax-ng/maintenance-dashboard/internal/data"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type metrics struct {
	info *prometheus.GaugeVec
}

var m *metrics

// Add metrics http handler
func CreateMetricsHandler() http.Handler {
	// Get rid of the default metrics
	r := prometheus.NewRegistry()
	// Add the metrics for all applications with Gauge type
	m = CreateMetrics(r)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	return handler
}

// Create new metric type with labels and add it to the registry.
func CreateMetrics(reg prometheus.Registerer) *metrics {
	metr := &metrics{
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "maintenance",
			Name:      "apps_version_info",
			Help:      "Information about current and latest versions for applications",
		}, []string{"app_name", "current_version", "latest_major_version", "latest_minor_version", "latest_patch_version"}),
	}

	reg.MustRegister(metr.info)
	log.Infoln("apps_version_info metric created and registered.")
	return metr
}

// Update metrics with the latest app version infos.
func UpdateMetrics(appsVersionInfo []*data.AppVersionInfo) {
	log.Infoln("Updating the metrics with the latest results.")
	for _, appVersionInfo := range appsVersionInfo {
		m.info.With(prometheus.Labels{
			"app_name":             appVersionInfo.NewReleasesName,
			"current_version":      appVersionInfo.CurrentVersion.String(),
			"latest_major_version": appVersionInfo.LatestMajorVersion.String(),
			"latest_minor_version": appVersionInfo.LatestMinorVersion.String(),
			"latest_patch_version": appVersionInfo.LatestPatchVersion.String(),
		}).Set(1)
	}
}
