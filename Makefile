M = $(shell printf "\033[34;1m▶\033[0m")

export GOBIN = $(CURDIR)/_bin
export PROTOPATH = $(GOPATH)/src
export BINDIR := $(GOBIN)
export PATH := $(GOBIN):$(PATH)

protoc_gen_go := $(GOBIN)/protoc-gen-go
protoc_gen_go_src := vendor/github.com/golang/protobuf/protoc-gen-go

protoc_gen_twirp := $(GOBIN)/protoc-gen-twirp
protoc_gen_twirp_src := vendor/github.com/twitchtv/twirp/protoc-gen-twirp

build: $(info $(M) Building project...)
	@ CGO_ENABLED=0 go build -o /go/bin/Example-TwirpService-Setup ./cmd/example-twirp-service/*.go
.PHONY: build

gen-twirp: $(protoc_gen_go) $(protoc_gen_twirp)
	$(info $(M) Generating twirp files...)
	@protoc --proto_path=$(GOBIN):. --twirp_out=. --go_out=. ./rpc/example-twirp-service/service.proto
.PHONY: gen-twirp

$(protoc_gen_go): $(protoc_gen_go_src)
	@go install ./$^

$(protoc_gen_twirp): $(protoc_gen_twirp_src)
	@go install ./$^

local:
	$(info $(M) Starting Development server...)
	go run ./cmd/example-twirp-service/main.go

docker-image:
	$(info $(M) Building app image...)
	docker build -t example-twirp-service .

docker-container: docker-image
	$(info $(M) Running docker application container...)
	docker run -p 8080:8080 example-twirp-service:latest

