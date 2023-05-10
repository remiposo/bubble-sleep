.PHONEY: build
.DEFAULT_GOAL := build

APP_NAME := "sleep"

build:
	@echo "Building..."
	@go build -o bin/$(APP_NAME) .

