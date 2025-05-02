default: fmt check-current-version lint install generate

build:
  go build -v ./...

install: build
  go install -v ./...

lint:
  golangci-lint run

generate:
  cd tools; go generate ./...

fmt:
  gofmt -s -w -e .

bump-version new_version:
  cd tools; go run version-util.go bump {{ new_version }}
  just generate

check-version version:
  cd tools; go run version-util.go check {{ version }}

check-current-version:
  #!/usr/bin/env bash
  set -ue
  CURRENT_VERSION=$(cat .version)
  cd tools; go run version-util.go check $CURRENT_VERSION
