.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	go run github.com/syumai/workers/cmd/workers-assets-gen@v0.18.0
	tinygo build -o ./build/app.wasm -target wasm -no-debug ./...

.PHONY: deploy
deploy:
	wrangler deploy

.PHONY: init-db
init-db:
	wrangler d1 execute d1vtfloat --file=./db/schema.sql

.PHONY: init-db-local
init-db-local:
	wrangler d1 execute d1vtfloat --local --file=./db/schema.sql

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go clean -testcache
	GOOS=js GOARCH=wasm go test -cover ./...

