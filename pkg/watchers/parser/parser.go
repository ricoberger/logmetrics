package parser

import "fmt"

const namespace = "logmetrics"
const metricParsedType = "parsed"
const metricParsedWithErrorType = "parsed_with_error"
const metricMatchedType = "matched"

// Config defines the structure for the parser configuration
type Config struct {
	ParserType string   `yaml:"type"`
	Keywords   []string `yaml:"keywords,omitempty"`
	Regexp     string   `yaml:"regexp,omitempty"`
}

// Parser is the interface each parser must implement
type Parser interface {
	Parse(podName, podNamespace string, logLine []byte) (bool, error)
}

// New returns a new parser or if the initialization of the parser fails an error
func New(name string, config Config) (Parser, error) {
	switch config.ParserType {
	case "contains":
		return newContainsParser(name, config)
	case "regexp":
		return newRegexParser(name, config)
	default:
		return nil, fmt.Errorf("Parser type %s is invalide", config.ParserType)
	}
}
