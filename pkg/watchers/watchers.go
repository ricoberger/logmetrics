package watchers

import (
	"io/ioutil"

	"github.com/ricoberger/logmetrics/pkg/kube"
	"github.com/ricoberger/logmetrics/pkg/watchers/parser"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Watcher defines the parser type and selector for pods to watch
type Watcher struct {
	Name      string        `yaml:"name"`
	Namespace string        `yaml:"namespace"`
	Selector  string        `yaml:"selector"`
	Parser    parser.Config `yaml:"parser"`
}

// Run starts a specific watcher
func (w *Watcher) Run(kubeClient *kube.Client) {
	logFields := log.Fields{
		"watcher": w.Name,
	}
	log.WithFields(logFields).Infof("Start watcher %s", w.Name)

	p, err := parser.New(w.Name, w.Parser)
	if err != nil {
		log.WithFields(logFields).WithError(err).Errorf("Could not create parser")
		return
	}

	addChan := make(chan string)
	go kubeClient.WatchPods(w.Namespace, w.Selector, addChan, logFields)

	for {
		select {
		case podName := <-addChan:
			log.WithFields(logFields).Infof("Start watching %s", podName)
			go kubeClient.ProcessLogs(w.Namespace, podName, p, logFields)
		}
	}
}

// ParseConfig parses the given configuration file for the watchers
func ParseConfig(file string) ([]Watcher, error) {
	var watchers []Watcher
	config, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(config, &watchers)
	if err != nil {
		return nil, err
	}

	return watchers, nil
}

// Run starts all given watchers
func Run(watchers []Watcher, kubeClient *kube.Client) {
	for _, watcher := range watchers {
		w := watcher
		go w.Run(kubeClient)
	}
}
