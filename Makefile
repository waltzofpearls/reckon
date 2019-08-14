all: build

build:
	go build

docker:
	docker build -t reckon .
