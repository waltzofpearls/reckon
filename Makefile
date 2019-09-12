all: build

build:
	go build -mod=vendor

proto:
	protoc --go_out=plugins=grpc:. api/*.proto
	python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. api/*.proto
	@tree -hrtC api

DOCKER_IMG := reckon
docker:
	docker build -t $(DOCKER_IMG) .
	docker run -it --rm \
		-e "TLS_ROOT_CA=$$(cat ../out/Reckon_Root_CA.crt)" \
		-e "TLS_SERVER_CERT=$$(cat ../out/localhost.crt)" \
		-e "TLS_SERVER_KEY=$$(cat ../out/localhost.key)" \
		-e "TLS_CLIENT_CERT=$$(cat ../out/StatsModel.crt)" \
		-e "TLS_CLIENT_KEY=$$(cat ../out/StatsModel.key)" \
		-e "GRPC_SERVER_ADDRESS=localhost:3003" \
		$(DOCKER_IMG) /bin/bash

server:
	TLS_ROOT_CA=$$(cat ../out/Reckon_Root_CA.crt) \
	TLS_SERVER_CERT=$$(cat ../out/localhost.crt) \
	TLS_SERVER_KEY=$$(cat ../out/localhost.key) \
	GRPC_SERVER_ADDRESS=localhost:3003 \
	PROM_CLIENT_URL=http://prometheus.rpi.topbass.studio:9090 \
		./reckon

client:
	TLS_ROOT_CA=$$(cat ../out/Reckon_Root_CA.crt) \
	TLS_CLIENT_CERT=$$(cat ../out/StatsModel.crt) \
	TLS_CLIENT_KEY=$$(cat ../out/StatsModel.key) \
	GRPC_SERVER_ADDRESS=localhost:3003 \
		python main.py
