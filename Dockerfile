FROM golang:1.26.4 AS builder

COPY . /go/src/app
WORKDIR /go/src/app

RUN CGO_ENABLED=0 GOOS=linux go build -o bin main.go

FROM alpine:3.24

RUN mkdir -p /opt/module-manager/service /opt/module-manager/data /opt/module-manager/modules /opt/module-manager/deployments
WORKDIR /opt/module-manager
COPY --from=builder /go/src/app/bin bin

ENTRYPOINT ["./bin"]
