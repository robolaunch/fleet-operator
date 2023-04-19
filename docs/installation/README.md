# Installation

## Prerequisites

For using robolaunch Fleet Operator, these prerequisites should be satisfied:

|     Tool     |       Version      |
|:------------:|:------------------:|
|  Kubernetes  |  `v1.21` and above |
| Cert-Manager | `v1.8.x` and above |
|    OpenEBS   | `v3.x.x` and above |
|    Robot Operator   | You can check the version map [here](). |

### Labeling Node

Select an active node from your cluster and add these labels:

```bash
kubectl label <NODE> robolaunch.io/organization=robolaunch
kubectl label <NODE> robolaunch.io/team=robotics
kubectl label <NODE> robolaunch.io/region=europe-east
kubectl label <NODE> robolaunch.io/cloud-instance=cluster
kubectl label <NODE> robolaunch.io/cloud-instance-alias=cluster-alias
```

## Installing Fleet Operator

### via Helm

Add robolaunch Helm repository and update:

```bash
helm repo add robolaunch https://robolaunch.github.io/charts/
helm repo update
```

Install latest version of Fleet Operator (remove `--devel` for getting latest stable version):

```bash
helm upgrade -i fleet-operator robolaunch/fleet-operator  \
--namespace fleet-system \
--create-namespace \
--devel
```

Or you can specify a version (remove the `v` letter at the beginning of the release or tag name):

```bash
VERSION="0.1.6-alpha.5"
helm upgrade -i fleet-operator robolaunch/fleet-operator  \
--namespace fleet-system \
--create-namespace \
--version $VERSION
```

To uninstall Fleet Operator installed with Helm, run the following commands:

```bash
helm delete fleet-operator -n fleet-system
kubectl delete ns fleet-system
```

### via Manifest

Deploy Fleet Operator one-file YAML using the command below:

```bash
# select a tag
TAG="v0.1.6-alpha.5"
kubectl apply -f https://raw.githubusercontent.com/robolaunch/fleet-operator/$TAG/hack/deploy/manifests/fleet_operator.yaml
```

To uninstall Fleet Operator installed with one-file YAML, run the following commands:
```bash
# find the tag you installed
TAG="v0.1.6-alpha.5"
kubectl delete -f https://raw.githubusercontent.com/robolaunch/fleet-operator/$TAG/hack/deploy/manifests/fleet_operator.yaml
```