APP := reckon
PYTHON_VERSION := 3.7.12
GO_VERSION := 1.16.10
GORELEASER_VERSION := 0.174.0
OSX_SDK_VERSION := 11.3
PORT := 8080:8080
PROM_CLIENT_URL ?= http://prometheus.rpi.topbass.studio:9090
PROM_EXPORTER_ADDR ?= :8080
# comma separated list or inline yaml
# WATCH_LIST ?= sensehat_temperature,sensehat_humidity
WATCH_LIST ?= {sensehat_temperature: [Prophet, Tangram], sensehat_humidity: [Prophet, Tangram]}
SCHEDULE ?= @every 10m
GRPC_SERVER_ADDRESS ?= localhost:18443
GRPC_ROOT_CA := $$(cat cert/gRPC_Root_CA.crt)
GRPC_SERVER_CERT := $$(cat cert/localhost.crt)
GRPC_SERVER_KEY := $$(cat cert/localhost.key)
GRPC_CLIENT_CERT := $$(cat cert/grpc_client.crt)
GRPC_CLIENT_KEY := $$(cat cert/grpc_client.key)

.PHONY: all
all: build

~/.virtualenvs/$(APP):
	( \
		source $$VIRTUALENVWRAPPER_SCRIPT; \
		mkvirtualenv $(APP); \
		pip install -r ./model/requirements.txt; \
		pip install pystan==2.19.1.1; \
		pip install prophet==1.0.1; \
	)

.PHONY: venv
venv: ~/.virtualenvs/$(APP)

.PHONY: build
build: cert
	go build

.PHONY: run
run: venv build
	PROM_CLIENT_URL=$(PROM_CLIENT_URL) \
	PROM_EXPORTER_ADDR=$(PROM_EXPORTER_ADDR) \
	WATCH_LIST="$(WATCH_LIST)" \
	SCHEDULE="$(SCHEDULE)" \
	GRPC_SERVER_ADDRESS=$(GRPC_SERVER_ADDRESS) \
	GRPC_ROOT_CA=$(GRPC_ROOT_CA) \
	GRPC_SERVER_CERT=$(GRPC_SERVER_CERT) \
	GRPC_SERVER_KEY=$(GRPC_SERVER_KEY) \
	GRPC_CLIENT_CERT=$(GRPC_CLIENT_CERT) \
	GRPC_CLIENT_KEY=$(GRPC_CLIENT_KEY) \
		./$(APP)

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: mock
mock:
	mockgen -package=mocks -mock_names=Logger=Logger \
		-destination=mocks/logger.go github.com/waltzofpearls/reckon/logs Logger

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		model/api/forecast.proto
	python -m grpc_tools.protoc -I. --python_out=. \
		--grpc_python_out=. model/api/forecast.proto
	@tree -hrtC model/api

cert:
	# create CA
	certstrap --depot-path cert init --common-name "gRPC Root CA"
	# create server cert request
	certstrap --depot-path cert request-cert --domain localhost
	# create client cert request
	certstrap --depot-path cert request-cert --cn grpc_client
	# sign server and client cert requests
	certstrap --depot-path cert sign --CA "gRPC Root CA" localhost
	certstrap --depot-path cert sign --CA "gRPC Root CA" grpc_client
	@tree -hrC cert

.PHONY: gen
gen: mock proto

.PHONY: cover
cover:
	echo "mode: count" > .coverage.out
	go test -coverprofile .coverage.tmp ./...
	tail -n +2 .coverage.tmp >> .coverage.out
	go tool cover -html=.coverage.out

IMAGE := build
VERSION := $(shell git describe --tags $(git rev-list --tags --max-count=1) | sed 's/^v//g')
OS := linux
ARCH := amd64

.PHONY: docker
docker:
	docker build \
		--progress=plain \
		--build-arg "PYTHON_VERSION=$(PYTHON_VERSION)" \
		--build-arg "GO_VERSION=$(GO_VERSION)" \
		--build-arg "APP=$(APP)" \
		--build-arg "VERSION=$(VERSION)" \
		--build-arg "OS=$(OS)" \
		--build-arg "ARCH=$(ARCH)" \
		-t $(APP)/$(IMAGE) \
		-f $(IMAGE).Dockerfile \
		.
	docker run --rm \
		-e "PROM_CLIENT_URL=$(PROM_CLIENT_URL)" \
		-e "PROM_EXPORTER_ADDR=$(PROM_EXPORTER_ADDR)" \
		-e "WATCH_LIST=$(WATCH_LIST)" \
		-e "SCHEDULE=$(SCHEDULE)" \
		-e "GRPC_SERVER_ADDRESS=$(GRPC_SERVER_ADDRESS)" \
		-e "GRPC_ROOT_CA=$(GRPC_ROOT_CA)" \
		-e "GRPC_SERVER_CERT=$(GRPC_SERVER_CERT)" \
		-e "GRPC_SERVER_KEY=$(GRPC_SERVER_KEY)" \
		-e "GRPC_CLIENT_CERT=$(GRPC_CLIENT_CERT)" \
		-e "GRPC_CLIENT_KEY=$(GRPC_CLIENT_KEY)" \
		-p $(PORT) \
		$(APP)/$(IMAGE)

.PHONY: debian
debian:
	make docker IMAGE=debian

# no alpine because https://pythonspeed.com/articles/alpine-docker-python/

.PHONY: release
release:
	make release-base RELEASE_ARGS="release --rm-dist"

.PHONY: release-dryrun
release-dryrun:
	make release-base RELEASE_ARGS="--rm-dist --skip-validate --skip-publish"

.PHONY: release-base
release-base:
	docker build \
		--progress=plain \
		--build-arg "PYTHON_VERSION=$(PYTHON_VERSION)" \
		--build-arg "GO_VERSION=$(GO_VERSION)" \
		--build-arg "GORELEASER_VERSION=$(GORELEASER_VERSION)" \
		--build-arg "OSX_SDK_VERSION=$(OSX_SDK_VERSION)" \
		-t $(APP)/release \
		-f release.Dockerfile \
		.
	docker run --rm \
		--env-file .release-env \
		-e "PYTHON_VERSION=$(PYTHON_VERSION)" \
		-e "GO_VERSION=$(GO_VERSION)" \
		-e "GORELEASER_VERSION=$(GORELEASER_VERSION)" \
		-v $$PWD:/go/src/$(APP) \
		-w /go/src/$(APP) \
		$(APP)/release \
		$(RELEASE_ARGS)

.PHONY: test-release
test-release:
	make release-base RELEASE_ARGS="--snapshot --skip-publish --rm-dist"
	make docker IMAGE=test-release VERSION=$$(ls -t dist/reckon_*_linux_amd64.tar.gz | head | sed 's|dist\/reckon_\(.*\)_linux_amd64.tar.gz|\1|')
