# kubernetes-warm-images

![status](https://img.shields.io/badge/Status-Beta-informational)
![version](https://img.shields.io/docker/v/just1689/warmimages)
![version](https://img.shields.io/badge/Helm-0.9.0-blue)
[![Docker](https://github.com/just1689/kubernetes-warm-images/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/just1689/kubernetes-warm-images/actions/workflows/docker-publish.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/just1689/kubernetes-warm-images)](https://goreportcard.com/report/github.com/just1689/kubernetes-warm-images)
[![codebeat badge](https://codebeat.co/badges/2aff7ff0-8af7-43ee-95dc-72bbbd098c4f)](https://codebeat.co/projects/github-com-just1689-kubernetes-warm-images-main)
[![Maintainability](https://api.codeclimate.com/v1/badges/a1f55c3e1e1518fdcaa5/maintainability)](https://codeclimate.com/github/just1689/kubernetes-warm-images/maintainability)

## Kubernetes Warm Images

The goal of this project is to "keep images warm" by pulling images onto nodes in-case they may be used in the near
future.

## Use Cases

- "Warm serverless images" - This could be useful to you if you're running serverless workloads on Kubernetes where the
  overhead for pulling images each time is consequential.
- "Warm critical images" - You may want your nodes to have images cached for critical workloads before they're actually
  needed.
- "I don't trust my container registry" - Your image server might not be as HA as your K8s cluster.

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

3. Generate your `values.yaml` file for using with the Helm chart

```bash

# Create a values.yaml file
echo "nats:" > values.yaml
echo "  url: \"nats://wi-nats-client:4222\"" >> values.yaml
# Get the username & password from NATs
echo "  username: $(kubectl get cm --namespace warm-images wi-nats -o jsonpath='{.data.*}' | grep -m 1 user | awk '{print $2}')" >> values.yaml 
echo "  password: $(kubectl get cm --namespace warm-images wi-nats -o jsonpath='{.data.*}' | grep -m 1 password | awk '{print $2}')" >> values.yaml
# Allow watching some namespaces & ignoring others
echo "list: \"*\"" >> values.yaml
echo "ignore: \"kube-system\"" >> values.yaml

```   

4. Install Warm Images

```bash

# Install 
helm repo add captains-charts https://storage.googleapis.com/captains-charts
helm repo update
helm install --namespace warm-images --values values.yaml wi captains-charts/warm-images

```

## Usage Guide

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

## Roadmap v1.0.0 - Stable

- Config: Exclude images that "contain".
- Integrate health check for Controller.
- Integrate health check for Agent.
- v1

## Roadmap v1.1.0 - Post Experientia

- Philosophical: Figure out exactly which resources could be watched (DaemonSets, Deployments, etc)
- Philosophical: Figure out how to `select`
- Philosophical: Figure out how to `skip`
- LabelSelectors?
- Lua for Controller-side custom logic?
- Lua for Agent-side custom logic?
- Tests - Go.
- Tests - Helm.
- Export Prometheus endpoint. Config for Helm.
- Grafana Dashboard
- Monitoring first pass - logs.

## Roadmap v1.2.0 - 🌟🌟🌟

- Support for ContainerD.
- Clean shutdown - Controller.
- Clean shutdown - Agent.
- Consider support for non-NATs streaming.

## Roadmap - 💭💭💭

- Test different scenarios (at scale, low availability, etc).
- Push Helm package to online repo as part of GitHub Action.
- K8s documentation (diagram, scaling etc)
- Helm documentation - all values tables & values example files.

