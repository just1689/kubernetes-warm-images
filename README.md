# kubernetes-warm-images

## Kubernetes Warm Images

The goal of this project is to "keep images warm" by pulling images onto nodes if they may be used in the near future.

## Architecture

The project will consist of three components.

1. The `Controller` is responsible for subscribing to relevant object changes in Kubernetes declarative state.
2. The `Local Agent` will ensure each node pulls images determined relevant by the `Controller`.
3. The communication medium for the project is `nats`. NATs is configurable by operator which means more freedom.

## Installation

1. Create the namespace

```bash
kubectl create ns warm-images
```   

2. Install NATs

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install --namespace warm-images wi-nats bitnami/nats
```

3. Install the app

```bash
# Get the username & password
touch values.yaml
echo "username: $(kubectl get cm --namespace warm-images wi-nats -o jsonpath='{.data.*}' | grep -m 1 user | awk '{print $2}')" >> values.yaml 
echo "password: $(kubectl get cm --namespace warm-images wi-nats -o jsonpath='{.data.*}' | grep -m 1 password | awk '{print $2}')" >> values.yaml

# Allow watching all namespaces
echo "list.spaces: *" >> values.yaml
# TODO: ignore namespaces
# ...

# Install 
helm repo add TBA
helm repo update
helm install --namespace warm-images --values values.yaml wi tba/tba
```   

## Configuration

### Changing the namespaces

Modify the `values.yaml` file for multiple namespaces by separating them with spaces

```yaml
list.spaces: ns1 default ns2
```

or use all namespaces using `*`:

```yaml
list.spaces: *
```

## Roadmap

- Support for ignoring namespaces
- Support for * namespaces
- Move builds to GitHub
- Helm Chart for installation
- Remove the local YAML
- End-to-end test of install guide
- Hosting for Helm chart?
- Redo logging in Controller
- ContainerD client option for pulling images
- Roughly "contains" images or regex ignoring