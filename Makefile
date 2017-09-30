.PHONY: build build-darwin build-linux

build:
	go build github.com/hayeah/solar/cli/solar

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o solar-darwin-amd64 github.com/hayeah/solar/cli/solar

build-linux:
	GOOS=linux GOARCH=amd64 go build -o solar-linux-amd64 github.com/hayeah/solar/cli/solar