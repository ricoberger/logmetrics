package kube

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"time"

	"github.com/ricoberger/logmetrics/pkg/watchers/parser"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client implements an API client for a Kubernetes cluster.
type Client struct {
	config    *rest.Config
	clientset *kubernetes.Clientset
}

// NewClient creates a new client to interact with the Kubernetes API server.
func NewClient(incluster bool) (*Client, error) {
	var config *rest.Config
	var err error

	if incluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		if _, ok := os.LookupEnv("HOME"); !ok {
			u, err := user.Current()
			if err != nil {
				return nil, fmt.Errorf("could not get current user: %w", err)
			}

			loadingRules.Precedence = append(loadingRules.Precedence, path.Join(u.HomeDir, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName))
		}

		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			loadingRules,
			&clientcmd.ConfigOverrides{},
		).ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:    config,
		clientset: clientset,
	}, nil
}

// WatchPods watches for all pods in the given namespace and the given selector
func (c *Client) WatchPods(namespace, selector string, addChan chan string, logFields log.Fields) {
	for {
		ctx := context.TODO()

		watcher, err := c.clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{
			LabelSelector: selector,
		})
		if err != nil {
			log.WithError(err).WithFields(logFields).Errorf("Could not watch pods")
		}

		ch := watcher.ResultChan()
		for event := range ch {
			log.WithFields(logFields).Infof("Process watch event")

			if event.Type == watch.Error {
				log.WithFields(logFields).Errorf("Error while watching: %#v", event.Object)
				continue
			}

			if event.Type == watch.Added {
				pod := event.Object.(*corev1.Pod)
				log.WithFields(logFields).Infof("Pod %s was added", pod.Name)
				addChan <- pod.Name
			}
		}
	}
}

// ProcessLogs listen to the logs of a Pod and processes each log line with the given parser
func (c *Client) ProcessLogs(namespace, name string, parse parser.Parser, logFields log.Fields) {
	var backoff int
	for {
		ctx := context.TODO()

		readCloser, err := c.clientset.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{
			Follow:    true,
			TailLines: int64Pointer(0),
		}).Stream(ctx)
		if err != nil {
			log.WithError(err).WithFields(logFields).Errorf("Could not stream logs for %s", name)
			backoff = backoff + 1
			time.Sleep(time.Duration(backoff) * 5 * time.Second)
		} else {
			defer readCloser.Close()

			for {
				select {
				default:
					p := make([]byte, 2048)
					n, err := readCloser.Read(p)
					if err != nil {
						if err == io.EOF {
							log.WithFields(logFields).Infof("Stop watching for %s", name)
							return
						}
						log.WithError(err).WithFields(logFields).Errorf("Could not read logs")
						backoff = backoff + 1
						time.Sleep(time.Duration(backoff) * 5 * time.Second)
					}

					if string(p[:n]) != "" {
						backoff = 0
						log.Debugf("Log line received: %s", string(p[:n]))
						_, err := parse.Parse(name, namespace, p[:n])
						if err != nil {
							log.WithError(err).WithFields(logFields).Errorf("Could not parse log line")
						}
					}
				}
			}
		}
	}
}

func int64Pointer(number int64) *int64 {
	pointer := number
	return &pointer
}
