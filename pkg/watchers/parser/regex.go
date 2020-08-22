package parser

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Regexp defines the structure of a regexp parser
// The regexp parser checks if a log line matches the given regular expression
type Regexp struct {
	regexp        *regexp.Regexp
	metricCounter *prometheus.CounterVec
}

func newRegexParser(name string, config Config) (*Regexp, error) {
	r, err := regexp.Compile(config.Regexp)
	if err != nil {
		return nil, err
	}

	return &Regexp{
		regexp: r,
		metricCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      name,
		}, []string{"type", "pod_namespace", "pod_name"}),
	}, nil
}

// Parse applies the parser to the given log line and increases the corresponding metrics
func (p *Regexp) Parse(podName, podNamespace string, logLine []byte) (bool, error) {
	p.metricCounter.With(prometheus.Labels{"type": metricParsedType, "pod_namespace": podNamespace, "pod_name": podName}).Inc()

	if p.regexp.Match(logLine) {
		p.metricCounter.With(prometheus.Labels{"type": metricMatchedType, "pod_namespace": podNamespace, "pod_name": podName}).Inc()
		return true, nil
	}

	return false, nil
}
