APP := reckon
IMAGE := $(APP)
PROMETHEUS := http://prometheus.rpi.topbass.studio:9090
METRICS := sensehat_temperature,sensehat_humidity

.PHONY: all
all: build

.PHONY: build
build:
	go build

.PHONY: run
run: build
	PROM_CLIENT_URL=$(PROMETHEUS) \
	WATCH_LIST=$(METRICS) \
		./$(APP)

.PHONY: venv
venv:
	$(shell mkvirtualenv $(APP))
	pip install -r ./model/requirements.txt

.PHONY: vrun
vrun:
	$(shell workon $(APP))
	PYTHONPATH=~/.virtualenvs/$(APP)/lib/python3.7/site-packages/:$$PYTHONPATH make run
	$(shell deactivate)

.PHONY: docker
docker:
	docker build -t $(IMAGE) .
	docker run -it --rm \
		-e "PROM_CLIENT_URL=$(PROMETHEUS)" \
		-e "WATCH_LIST=$(METRICS)" \
		$(IMAGE)
