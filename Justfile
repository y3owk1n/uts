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

release-ci VERSION_OVERRIDE:
    @echo "Building release artifacts..."
    mkdir -p build
    env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/y3owk1n/uts/cmd.Version={{ VERSION_OVERRIDE }}" -trimpath -o build/uts-darwin-arm64 ./main.go
    env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/y3owk1n/uts/cmd.Version={{ VERSION_OVERRIDE }}" -trimpath -o build/uts-darwin-amd64 ./main.go
    env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/y3owk1n/uts/cmd.Version={{ VERSION_OVERRIDE }}" -trimpath -o build/uts-linux-arm64 ./main.go
    env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/y3owk1n/uts/cmd.Version={{ VERSION_OVERRIDE }}" -trimpath -o build/uts-linux-amd64 ./main.go
    @echo "✓ Release artifacts built"

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
