.PHONY: build

build:
	sam build

local-init-storage:
	docker stop dynamodb || true && docker rm dynamodb || true
	docker run -d -p 8000:8000 --name dynamodb amazon/dynamodb-local
	aws dynamodb create-table --cli-input-json file://.aws-sam/storage-table.json --endpoint-url http://localhost:8000 2>&1 > /dev/null

start-api: build local-init-storage
	sam local start-api \
	--parameter-overrides $(shell cat .aws-sam/development-params)

deploy: build
	sam deploy \
	--parameter-overrides $(shell cat .aws-sam/production-params) \
	--stack-name=url-shortener \
	--resolve-s3 \
	--capabilities=CAPABILITY_IAM
