# kubernetes-warm-images

## Kubernetes Warm Images

The goal of this project is to pull images onto nodes if they may be used there in the near future. Another way of
putting this is "Keeping images warm".

## Architecture

The project will consist of three components.

1. The `Controller` is responsible for subscribing to relevant object changes in Kubernetes declarative state.
2. The `Local Agent` will ensure each node pulls images determined relevant by the `Controller`.
3. The communication medium for the project is `nats`. NATs is configurable by operator which means more freedom.

## Installation

1. Create the namespace
2. Install NATs
3. Install Warm Images

## Configuration

### Changing the namespaces

Create values.yaml file with the following entry:
```yaml
  list.spaces: ns1 default ns2
```
or this:
```yaml
  list.spaces: *
```
