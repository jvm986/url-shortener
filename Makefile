.PHONY: build

build:
	sam build

start-api: build
	sam local start-api