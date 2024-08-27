ifneq (,$(wildcard ./.env))
	include .env
	export
endif

ifeq ($(OS),Windows_NT)
	GOCMD=go
else
	GOCMD=/usr/local/go/bin/go
endif

GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
APP_NAME=morph
BINARY_LINUX=$(APP_NAME)-cli_linux
BINARY_WINDOWS=$(APP_NAME)-cli.exe
UI_DIR=./ui

.PHONY:

build:
	go build -tags json ./...

all: print json aws azure gcp mssql

linux: linux-print linux-json linux-aws linux-azure linux-gcp linux-mssql

windows: windows-print windows-json windows-aws windows-azure windows-gcp windows-mssql

print:
	GOOS=linux $(GOBUILD) -tags print ./cmd/$(APP_NAME)

json:
	GOOS=linux $(GOBUILD) -tags json ./cmd/$(APP_NAME)

aws:
	GOOS=linux $(GOBUILD) -tags aws ./cmd/$(APP_NAME)

azure:
	GOOS=linux $(GOBUILD) -tags azure ./cmd/$(APP_NAME)

gcp:
	GOOS=linux $(GOBUILD) -tags gcp ./cmd/$(APP_NAME)

mssql:
	GOOS=linux $(GOBUILD) -tags mssql ./cmd/$(APP_NAME)

linux-print:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags print -o $(APP_NAME)-print ./cmd/$(APP_NAME)

linux-json:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags json -o $(APP_NAME)-json ./cmd/$(APP_NAME)

linux-aws:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags aws -o $(APP_NAME)-aws ./cmd/$(APP_NAME)

linux-azure:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags azure -o $(APP_NAME)-azure ./cmd/$(APP_NAME)

linux-gcp:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags gcp -o $(APP_NAME)-gcp ./cmd/$(APP_NAME)

linux-mssql:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags mssql -o $(APP_NAME)-mssql ./cmd/$(APP_NAME)

windows-print:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags print -o $(APP_NAME)-print.exe ./cmd/$(APP_NAME)

windows-json:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags json -o $(APP_NAME)-json.exe ./cmd/$(APP_NAME)

windows-aws:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags aws -o $(APP_NAME)-aws.exe ./cmd/$(APP_NAME)

windows-azure:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags azure -o $(APP_NAME)-azure.exe ./cmd/$(APP_NAME)

windows-gcp:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags gcp -o $(APP_NAME)-gcp.exe ./cmd/$(APP_NAME)

windows-mssql:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -tags mssql -o $(APP_NAME)-mssql.exe ./cmd/$(APP_NAME)

clean:
	rm morph-* main
