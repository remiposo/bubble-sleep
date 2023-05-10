.PHONEY: build
.DEFAULT_GOAL := build

APP_NAME := "bubble-sleep"

clean:
	@echo "Cleaning..."
	@rm -rf bin/*

build: clean
	@echo "Building..."
	@go build -o bin/$(APP_NAME) .

