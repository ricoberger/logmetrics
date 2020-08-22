# logmetrics

**logmetrics** turns your log messages into useful metrics for Prometheus. logmetrics watches the logs of your Pods in your Kubernetes cluster and parses each line. If a line matches your defined parser a metrics counter will be increased.

For example [Flux](http://github.com/fluxcd/flux) logs the following lines, when it will fail to auto-release workloads:

```
{"caller":"images.go:17","component":"sync-loop","msg":"polling for new images for automated workloads","ts":"2020-08-22T08:16:15.285644598Z"}
{"caller":"warming.go:198","component":"warmer","image":"goharbor/registry-photon","info":"refreshing image","of_which_missing":1,"of_which_refresh":0,"tag_count":59,"to_update":1,"ts":"2020-08-22T08:16:43.189668642Z"}
{"caller":"repocachemanager.go:223","canonical_name":"index.docker.io/goharbor/registry-photon","component":"warmer","impact":"flux will fail to auto-release workloads with matching images, ask the repository administrator to fix the inconsistency","ts":"2020-08-22T08:16:43.42114786Z","warn":"manifest for tag v2.7.1-debug-9186 missing in repository goharbor/registry-photon"}
{"attempted":1,"caller":"warming.go:206","component":"warmer","successful":0,"ts":"2020-08-22T08:16:43.421301345Z","updated":"goharbor/registry-photon"}
{"caller":"images.go:17","component":"sync-loop","msg":"polling for new images for automated workloads","ts":"2020-08-22T08:16:43.423436631Z"}
```

Then you can use the following configuration, to increase a counter named `logmetrics_flux_fail_to_autorelease_total` each time a log line contains the string `fail to auto-release`:

```yaml
- name: flux_fail_to_autorelease_total
  namespace: flux
  selector: app=flux
  parser:
    type: contains
    keywords:
      - fail to auto-release
```

The resulting Prometheus metrics are looking as follow. logmetrics parsed 5 log lines and one line matched our defined parser:

```
# HELP logmetrics_flux_fail_to_autorelease_total
# TYPE logmetrics_flux_fail_to_autorelease_total counter
logmetrics_flux_fail_to_autorelease_total{pod_name="flux-76b4999bc7-cgffd",pod_namespace="flux",type="matched"} 1
logmetrics_flux_fail_to_autorelease_total{pod_name="flux-76b4999bc7-cgffd",pod_namespace="flux",type="parsed"} 5
```

## Installation

You can install logmetrics into your Kubernetes cluster via Kustomize or Helm. To install logmetrics via Kustomize run the following commands:

```sh
kubectl create ns logmetrics
kubectl kustomize github.com/ricoberger/logmetrics --namespace logmetrics
```

Or if you want to install it via Helm run the following:

```sh
helm repo add ricoberger https://ricoberger.github.io/helm-charts
helm repo update

helm upgrade --install logmetrics ricoberger/logmetrics
```

## Configuration

logmetrics is configured via a yaml file. Within this file you define your metrics, the Pods which should be watched and the parser wich should be used. An example configuration file can be found in the [config.yaml](./config.yaml) file within this repository.

| Property | Description |
| -------- | ----------- |
| `name` | Is the name of the metric which should be written for Prometheus. The defined name is prefixed with `logmetrics`. The resulting metric contains a field `type` which indicates the number of parsed log lines (`parsed`), the number of lines which matched your parser (`matched`) and the number of parsing failures (`parsed_with_error`). |
| `namespace` | Is the namespace of the Pod which log messages should be watched. |
| `selector` | Is the selector of the Pod which log messages should be watched. |
| `parser.type` | Is the type of the parser which should be used. Must be `contains` or `regexp` |
| `parser.keywords` | Used when `parser.type` is `contains`. A list of strings which a log line must contain to increase the counter. |
| `parser.regexp` | Used when `parser.type` is `regexp`. The regular expression which must be matched by a log line to increase the counter. |

## Development

logmetrics is written in Go. To build the binary you have to run the following commands:

```sh
git clone git@github.com:ricoberger/logmetrics.git
make build

./bin/logmetrics
```
