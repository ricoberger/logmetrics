# logmetrics

| Value | Description | Default |
| ----- | ----------- | ------- |
| `replicaCount` | Number of replicas which should be created. | `1` |
| `image.repository` | The repository of the Docker image. | `ricoberger/logmetrics` |
| `image.tag` | The tag of the Docker image which should be used. | `0.1.0` |
| `image.pullPolicy` | The pull policy for the Docker image, | `IfNotPresent` |
| `imagePullSecrets` | Secrets which can be used to pull the Docker image. | `[]` |
| `nameOverride` | Expand the name of the chart. | `""` |
| `fullnameOverride` | Override the name of the app. | `""` |
| `rbac.create` | Create the cluster role and cluster role binding. | `true` |
| `rbac.name` | The name of the cluster role and cluster role binding, which should be created/used by logmetrics. | `logmetrics` |
| `serviceAccount.create` | Create the service account. | `true` |
| `serviceAccount.name` | The name of the service account, which should be created/used by logmetrics. | `logmetrics` |
| `resources` | Set resources for the operator. | `{}` |
| `nodeSelector` | Set a node selector. | `{}` |
| `tolerations` | Set tolerations. | `[]` |
| `affinity` | Set affinity. | `{}` |
| `config` | The content of the configuration file which should be used by logmetrics. | `""` |
| `serviceMonitor.create` | Create a ServiceMonitor for the Prometheus Operator. | `false` |
| `serviceMonitor.labels` | Additional labels which should be set for the ServiceMonitor. | `{}` |
| `serviceMonitor.interval` | Scrape interval. | `10s` |
| `serviceMonitor.scrapeTimeout` | Scrape timeout. | `10s` |
| `serviceMonitor.honorLabels` | Honor labels option. | `true` |
| `serviceMonitor.relabelings` | Additional relabeling config for the ServiceMonitor. | `[]` |
