# kubernetes-warm-images

The goal of this project is to keep container images on nodes that could be used in the near future to reduce the time needed to spin them up in the case of node downtime, scaling or similar scenarios.

## Architecture

The project will consist of two components. The `Controller` is responsible for subscribing to relevant object changes in Kubernetes declarative state. The `Local Agent` will ensure each node pulls images determined relevant by the `Controller`.

