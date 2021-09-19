package config

import "github.com/prometheus/client_golang/prometheus"

type BuildInfo struct {
	Version           string
	Commit            string
	Date              string
	GoVersion         string
	PythonVersion     string
	GoreleaserVersion string

	desc *prometheus.Desc
}

func NewBuildInfo(version, commit, date, goVersion, pythonVersion, goreleaserVersion string) *BuildInfo {
	return &BuildInfo{
		Version:           version,
		Commit:            commit,
		Date:              date,
		GoVersion:         goVersion,
		PythonVersion:     pythonVersion,
		GoreleaserVersion: goreleaserVersion,

		desc: prometheus.NewDesc(
			"reckon_build_info",
			"Information about reckon build.",
			[]string{"version", "commit", "date", "go_version", "python_version", "goreleaser_version"},
			nil,
		),
	}
}

func (b *BuildInfo) Describe(ch chan<- *prometheus.Desc) {
	ch <- b.desc
}

func (b *BuildInfo) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		b.desc,
		prometheus.GaugeValue,
		1,                                                                              // metric value
		b.Version, b.Commit, b.Date, b.GoVersion, b.PythonVersion, b.GoreleaserVersion, // metric labels
	)
}
