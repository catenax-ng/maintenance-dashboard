package metrics

import (
	"net/http"
	"time"

	"github.com/catenax-ng/maintenance-dashboard/internal/data"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	info *prometheus.GaugeVec
}

var m *metrics

// Add metrics http handler
func CreateMetrics() http.Handler {
	// get rid of the default metrics
	r := prometheus.NewRegistry()
	// add the metrics for all applications with Gauge type
	m = NewMetrics(r)
	//addAppMetrics(r, appsVersionInfo)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	return handler
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	metr := &metrics{
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "maintenance",
			Name:      "apps_version_info",
			Help:      "Information about current and latest versions for applications",
		}, []string{"app_name", "current_version", "latest_major_version", "latest_minor_version", "latest_patch_version"}),
	}

	reg.MustRegister(metr.info)
	return metr
}

func UpdateMetrics(appsVersionInfo []*data.AppVersionInfo) {
	for _, appVersionInfo := range appsVersionInfo {
		m.info.With(prometheus.Labels{
			"app_name":             appVersionInfo.NewReleasesName,
			"current_version":      appVersionInfo.CurrentVersion.String(),
			"latest_major_version": appVersionInfo.LatestMajorVersion.String(),
			"latest_minor_version": appVersionInfo.LatestMinorVersion.String(),
			"latest_patch_version": appVersionInfo.LatestPatchVersion.String(),
		}).Set(1)
	}
	time.Sleep(time.Duration(10 * float64(time.Second)))
}
