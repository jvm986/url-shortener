.PHONY: build

build:
	sam build

start-api: build
	sam local start-api \
	--parameter-overrides $(shell cat .aws-sam/development-params)