.PHONY: build

build:
	export GO111MODULE=on
	env GOOS=${OS} GOARCH=${ARCH} go build -ldflags="-s -w" -o bin/taxi-cli-${OS}-${ARCH} cmd/*.go
