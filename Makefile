VERSION=`git describe --tags --abbrev=0`
COMMIT=`git rev-parse HEAD`

.PHONY: init
init:
	# pip install pre-commit
	# pre-commit install
	# go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.46.2
	# for ginkgo
	go install github.com/onsi/ginkgo/v2/ginkgo@latest
	# for make mock
	go install github.com/golang/mock/mockgen@v1.6.0
	# for gofumpt
	go install mvdan.cc/gofumpt@latest
	# for golines
	go install github.com/segmentio/golines@latest
	# for goimports
	go install -v github.com/incu6us/goimports-reviser/v3@latest
	# for swag
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: dep
dep: 
	go mod tidy
	go mod vendor

.PHONY: doc
doc:
	swag init

.PHONY: mock
mock:
	go generate ./...

.PHONY: lint
lint:
	export GOFLAGS=-mod=vendor
	golangci-lint run --fast --skip-dirs mock

.PHONY: fmt
fmt:
	golines ./ -m 120 -w --base-formatter gofmt --no-reformat-tags
	gofumpt  -l -w .
	goimports-reviser -rm-unused -set-alias -format ./...

.PHONY: test
test:
	go test -mod=vendor -gcflags=all=-l $(shell go list ./... | grep -v mock | grep -v docs) -covermode=count -coverprofile .coverage.cov

.PHONY: cov
cov:
	go tool cover -html=.coverage.cov

.PHONY: bench
bench:
	go test -run=nonthingplease -benchmem -bench=. $(shell go list ./... | grep -v /vendor/)

.PHONY: build
build:
	go build -mod=vendor -ldflags "\
		-X mygo/template/pkg/version.Version=${VERSION} \
		-X mygo/template/pkg/version.Commit=${COMMIT} \
		-X mygo/template/pkg/version.BuildTime=`date +%Y-%m-%d_%I:%M:%S` \
		-X mygo/template/pkg/version.GoVersion=`go version` \
		-w" .

.PHONY: serve
serve: build
	./mygo -c config.yaml

.PHONY: dev-image
dev-image:
	docker build --build-arg VERSION=${VERSION} --build-arg COMMIT=${COMMIT}  -f  Dockerfile . -t mygo:development

