all: build

build:
	go build -mod=vendor

proto:
	protoc --go_out=plugins=grpc:. api/*.proto
	python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. api/*.proto
	@tree -hrtC api

DEV_IMG := reckon
dev:
	docker build -t $(DEV_IMG) .
	docker run -it --rm $(DEV_IMG) /bin/bash
