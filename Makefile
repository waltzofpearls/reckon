APP := reckon
PYTHON_VERSION := 3.7.11
GO_VERSION := 1.16.5
GORELEASER_VERSION := 0.174.0
OSX_SDK_VERSION := 10.12
PORT := 8080:8080
PROM_CLIENT_URL ?= http://prometheus.rpi.topbass.studio:9090
PROM_EXPORTER_ADDR ?= :8080
WATCH_LIST ?= sensehat_temperature,sensehat_humidity
SCHEDULE ?= @every 10m

.PHONY: all
all: build

~/.virtualenvs/$(APP):
	( \
		source $$VIRTUALENVWRAPPER_SCRIPT; \
		mkvirtualenv $(APP); \
		pip install -r ./model/requirements.txt; \
	)

.PHONY: venv
venv: ~/.virtualenvs/$(APP)

.PHONY: build
build:
	go build

.PHONY: run
run: venv build
	PYTHONPATH=~/.virtualenvs/$(APP)/lib/python3.7/site-packages/:$$PYTHONPATH \
	PROM_CLIENT_URL=$(PROM_CLIENT_URL) \
	PROM_EXPORTER_ADDR=$(PROM_EXPORTER_ADDR) \
	WATCH_LIST=$(WATCH_LIST) \
	SCHEDULE="$(SCHEDULE)" \
		./$(APP)

.PHONY: test
test:
	go test -cover -race ./...

gen:
	mockgen -package=mocks -mock_names=Logger=Logger \
		-destination=mocks/logger.go github.com/waltzofpearls/reckon/logs Logger

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
