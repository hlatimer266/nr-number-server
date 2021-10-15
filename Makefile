.PHONY: test
test:
	@go test -mod=vendor -v -race ./...

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: vendor
vendor:
	@go mod vendor

.PHONY: build
build: fmt test
build:
	@go build -o bin/main main.go
	
.PHONY: run
run:
	@./bin/main

.PHONY: build-pvs
build-pvs:
	@go build -o performance/bin performance/cmd/performance.go

.PHONY: test-pvs
test-pvs:
	@./performance/bin/performance