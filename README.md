# kubernetes-warm-images

## Kubernetes Warm Images

The goal of this project is to "keep images warm" by pulling images onto nodes in-case they may be used in the near
future.

## Architecture

The project will consist of three components.

1. The `Controller` is responsible for subscribing to relevant object changes in Kubernetes declarative state.
2. The `Agent` runs on each node and pulls images to the node.
3. The communication medium for the project is `nats`. NATs is configurable by the operator which means more freedom for
   those running the project.

## Installation

1. Create the namespace

```bash
kubectl create ns warm-images
```   

2. Install NATs

I suggest using the Bitnami NATs Helm chart as at time of writing it is maintained, up-to-date and fairly
configurable: https://github.com/bitnami/charts/tree/master/bitnami/nats

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install --namespace warm-images wi-nats bitnami/nats
```

3. Install the app

```bash
# Get the username & password from NATs
touch values.yaml
echo "nats:" >> values.yaml
echo "  url: \"nats://wi-nats-client:4222\"" >> values.yaml
echo "  username: $(kubectl get cm --namespace warm-images wi-nats -o jsonpath='{.data.*}' | grep -m 1 user | awk '{print $2}')" >> values.yaml 
echo "  password: $(kubectl get cm --namespace warm-images wi-nats -o jsonpath='{.data.*}' | grep -m 1 password | awk '{print $2}')" >> values.yaml

# Allow watching some namespaces & ignoring others
echo "list: \"*\"" >> values.yaml
echo "ignore: \"\"" >> values.yaml

# Install 
helm repo add captains-charts https://storage.googleapis.com/captains-charts
helm repo update
helm install --namespace warm-images --values values.yaml wi captains-charts/warm-images
```   

## Configuration

### Changing the namespaces

Modify the `values.yaml` file for multiple namespaces by separating them with spaces in the field `list.spaces`

```yaml
list: "ns1 default ns2"
```

or use all namespaces using `*`:

```yaml
list: "*"
```

To ignore some number of namespaces modify the `ignore` field in the `values.yaml`. It accepts a spaced separated list
of namespaces.

## Monitoring
TBA

## Roadmap v0.9.0 - Core Functionality

- Remove temp password from `values.yaml`

## Roadmap v1.0.0 - Stable

- Config: Exclude images that "contain".
- Integrate health check for Controller.
- Integrate health check for Agent.
- Test different scenarios.
- End-to-end test of install guide.
- v1

## Roadmap v1.1.0 - Find all the bugs
- Tests - Go.
- Tests - Helm.
- Monitoring first pass - logs.

## Roadmap v1.2.0

- Support for ContainerD.
- Clean shutdown - Controller.
- Clean shutdown - Agent.
- Consider support for non-NATs streaming.

## Roadmap - Future
- Push Helm package to online repo as part of GitHub Action.
- Export Prometheus endpoint. Config for Helm.
