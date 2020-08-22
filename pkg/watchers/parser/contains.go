package parser

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Contains defines the structure of a contains parser
// The contains parser checks if a log line contains the given keywords
type Contains struct {
	keywords      []string
	metricCounter *prometheus.CounterVec
}

func newContainsParser(name string, config Config) (*Contains, error) {
	return &Contains{
		keywords: config.Keywords,
		metricCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      name,
		}, []string{"type", "pod_namespace", "pod_name"}),
	}, nil
}

// Parse applies the parser to the given log line and increases the corresponding metrics
func (p *Contains) Parse(podName, podNamespace string, logLine []byte) (bool, error) {
	p.metricCounter.With(prometheus.Labels{"type": metricParsedType, "pod_namespace": podNamespace, "pod_name": podName}).Inc()

	for _, keyword := range p.keywords {
		if !strings.Contains(string(logLine), keyword) {
			return false, nil
		}
	}

	p.metricCounter.With(prometheus.Labels{"type": metricMatchedType, "pod_namespace": podNamespace, "pod_name": podName}).Inc()
	return true, nil
}
