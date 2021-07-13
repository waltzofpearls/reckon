APP := reckon
IMAGE := $(APP)
PROMETHEUS ?= http://prometheus.rpi.topbass.studio:9090
WATCH_LIST ?= sensehat_temperature,sensehat_humidity
SCHEDULE ?= @every 5m

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
	PROM_CLIENT_URL=$(PROMETHEUS) \
	WATCH_LIST=$(WATCH_LIST) \
	SCHEDULE="$(SCHEDULE)" \
		./$(APP)

.PHONY: test
test:
	go test -cover ./...

.PHONY: cover
cover:
	echo "mode: count" > .coverage.out
	go test -coverprofile .coverage.tmp ./...
	tail -n +2 .coverage.tmp >> .coverage.out
	go tool cover -html=.coverage.out

.PHONY: docker
docker:
	docker build -t $(IMAGE) .
	docker run -it --rm \
		-e "PROM_CLIENT_URL=$(PROMETHEUS)" \
		-e "WATCH_LIST=$(WATCH_LIST)" \
		-e "SCHEDULE=$(SCHEDULE)" \
		$(IMAGE)
