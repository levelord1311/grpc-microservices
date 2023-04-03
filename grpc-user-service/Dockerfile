# Builder
FROM golang:1.20-alpine as builder
RUN apk add --update make git protoc protobuf protobuf-dev curl gcc musl-dev

ARG GITHUB_PATH=github.com/levelord1311/grpc-microservices/grpc-user-service

WORKDIR /home/${GITHUB_PATH}

COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum
RUN make deps
COPY . .

RUN make build


# gRPC Server
FROM alpine:latest as server
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

ARG GITHUB_PATH=github.com/levelord1311/grpc-microservices/grpc-user-service

COPY --from=builder /home/${GITHUB_PATH}/bin/app .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .
#COPY --from=builder /home/${GITHUB_PATH}/swagger ./swagger
COPY --from=builder /home/${GITHUB_PATH}/migrations/ ./migragtions

RUN chown root:root app

EXPOSE 50051
EXPOSE 6001
EXPOSE 9100

CMD ["./app"]