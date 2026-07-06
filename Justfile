VERSION := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
GIT_COMMIT := `git rev-parse --short HEAD 2>/dev/null || echo "unknown"`
BUILD_DATE := `date -u +"%Y-%m-%dT%H:%M:%SZ"`

LDFLAGS := "-s -w -X github.com/y3owk1n/uts/cmd.Version=" + VERSION + " -X github.com/y3owk1n/uts/cmd.GitCommit=" + GIT_COMMIT + " -X github.com/y3owk1n/uts/cmd.BuildDate=" + BUILD_DATE

default: build

build:
    @echo "Building uts..."
    @echo "Version: {{ VERSION }}"
    env CGO_ENABLED=0 go build -ldflags="{{ LDFLAGS }}" -trimpath -o bin/uts ./main.go
    @echo "✓ Build complete: bin/uts"

release-ci-darwin ARCH VERSION_OVERRIDE:
    @echo "Building release artifact (darwin/{{ ARCH }}) for CI..."
    @echo "Version: {{ VERSION_OVERRIDE }}"
    @echo "Commit: {{ GIT_COMMIT }}"
    @echo "Date: {{ BUILD_DATE }}"
    mkdir -p bin
    CGO_ENABLED=0 GOOS=darwin GOARCH={{ ARCH }} go build -ldflags="-s -w -X github.com/y3owk1n/uts/cmd.Version={{ VERSION_OVERRIDE }} -X github.com/y3owk1n/uts/cmd.GitCommit={{ GIT_COMMIT }} -X github.com/y3owk1n/uts/cmd.BuildDate={{ BUILD_DATE }}" -trimpath -o bin/uts-darwin-{{ ARCH }} ./main.go
    @echo "✓ Release artifact for darwin/{{ ARCH }} built successfully"

release-ci-linux ARCH VERSION_OVERRIDE:
    @echo "Building release artifact (linux/{{ ARCH }}) for CI..."
    @echo "Version: {{ VERSION_OVERRIDE }}"
    @echo "Commit: {{ GIT_COMMIT }}"
    @echo "Date: {{ BUILD_DATE }}"
    mkdir -p bin
    CGO_ENABLED=0 GOOS=linux GOARCH={{ ARCH }} go build -ldflags="-s -w -X github.com/y3owk1n/uts/cmd.Version={{ VERSION_OVERRIDE }} -X github.com/y3owk1n/uts/cmd.GitCommit={{ GIT_COMMIT }} -X github.com/y3owk1n/uts/cmd.BuildDate={{ BUILD_DATE }}" -trimpath -o bin/uts-linux-{{ ARCH }} ./main.go
    @echo "✓ Release artifact for linux/{{ ARCH }} built successfully"

test: test-unit test-integration

test-unit:
    go test ./... -v

test-integration:
    go test -tags=integration ./... -v

test-race: test-race-unit test-race-integration

test-race-unit:
    go test -race ./... -v

test-race-integration:
    go test -tags=integration -race ./... -v

test-coverage:
    go test -coverprofile=coverage.txt ./...

test-coverage-all:
    go test -tags=integration -coverprofile=coverage-all.txt ./...

test-coverage-html:
    just test-coverage
    go tool cover -html=coverage.txt -o coverage.html

test-all: test test-race

vet:
    go vet ./...

fmt:
    golangci-lint fmt
    golangci-lint run --fix

lint:
    golangci-lint run

# Generate man pages
genman OUTPUT_DIR="build/man":
    @echo "Generating man pages..."
    @mkdir -p {{ OUTPUT_DIR }}
    env CGO_ENABLED=0 go run ./cmd/genman {{ OUTPUT_DIR }}
    @echo "Man pages generated in {{ OUTPUT_DIR }}/"

clean:
    rm -rf bin/ build/ coverage*.txt coverage*.html

deps:
    go mod download
    go mod tidy
