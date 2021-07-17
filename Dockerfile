FROM reg.captainjustin.space/golang:1.16.6 as packed

COPY . /src
WORKDIR /src
RUN go mod vendor
RUN go mod download

FROM packed as builder
WORKDIR /src
RUN go build bin/app.go

# TODO: multistage build with light final image
FROM reg.captainjustin.space/debian:stable-20210621-slim
RUN mkdir /src
WORKDIR /src
COPY --from=builder /src/app /src/app
ENTRYPOINT ["./app"]
