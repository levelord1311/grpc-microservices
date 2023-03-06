API_DIR=$(CURDIR)/api
SERVICE_NAME=user-service-api
SERVICE_PATH=levelord1311/grpc-microservices/grpc-user-service

BUF_VERSION:="v1.15.0"

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)

#-----------------------------------------------------------

.PHONY: build
build: gen
	go build -o bin/app$(shell go env GOEXE) cmd/grpc-server/main.go

.PHONY: run
run:
	go run cmd/grpc-server/main.go

#-----------------------------------------------------------

.PHONY: gen
gen: clean .generate-install-buf .generate-go .generate-finalize-go

.generate-install-buf:
	@ command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
    		curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) --create-dirs -o "$(BUF_EXE)" && \
    		chmod +x "$(BUF_EXE)")

.generate-go:
	$(BUF_EXE) generate
.generate-finalize-go:
	mv pkg/$(SERVICE_NAME)/github.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME)/* pkg/$(SERVICE_NAME)
	rm -rf pkg/$(SERVICE_NAME)/github.com/
	cd pkg/$(SERVICE_NAME) && ls go.mod || (go mod init github.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME) && go mod tidy)

#-----------------------------------------------------------

.PHONY: deps
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

#-----------------------------------------------------------

.PHONY: lint
lint:
	cd $(API_DIR) && buf lint --config buf.yaml

.PHONY: format
format:
	cd $(API_DIR) && buf format

.PHONY: clean
clean:
	rm -rf pkg || true
	rm -rf bin || true


