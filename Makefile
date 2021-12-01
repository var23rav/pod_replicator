all: prepare build
prepare:
	@echo "prepare"
	go mod download
build:
	@echo "build"
	go build
