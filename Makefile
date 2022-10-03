GOCMD=go
GOBUILD=$(GOCMD) build
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)

LOMBA=bin/lomba
LOMBA-DARWIN=bin/lomba-darwin


all: target target-darwin

clean:
	rm -rf ${LOMBA} 

target:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${LOMBA} github.com/zawachte/lomba
target-darwin:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${LOMBA-DARWIN} github.com/zawachte/lomba


unit-test:
	go test ./...
