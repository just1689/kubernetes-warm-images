# kubernetes-warm-images

The goal of this project is to pull images onto nodes if they may be used there in the near future. Another way of
putting this is "Keeping images warm".

## Architecture

The project will consist of three components.

1. The `Controller` is responsible for subscribing to relevant object changes in Kubernetes declarative state.
2. The `Local Agent` will ensure each node pulls images determined relevant by the `Controller`.
3. The communication medium for the project is NATs. NATs is configurable by operator which means more freedom.



