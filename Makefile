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

.PHONY:

build:
	go build ./...

errcheck:
	errcheck ./...

morph-gnt:
	$(GORUN) ./cmd/morph -mode gnt -target jsonl

morph-wlc:
	$(GORUN) ./cmd/morph -mode wlc -style english -target jsonl

render-gnt:
	$(GORUN) ./cmd/render -mode gnt -target text

render-wlc:
	$(GORUN) ./cmd/render -mode wlc -style english -target text

counts-gnt:
	$(GORUN) ./cmd/counts -mode gnt -target text

counts-wlc:
	$(GORUN) ./cmd/counts -mode wlc -style english -target text

counts-gnt-jsonl:
	$(GORUN) ./cmd/counts -mode gnt -target jsonl

counts-wlc-jsonl:
	$(GORUN) ./cmd/counts -mode wlc -target jsonl

strongs-greek-gob:
	$(GORUN) ./cmd/morph -mode strongs-greek -target gob

strongs-greek:
	$(GORUN) ./cmd/morph -mode strongs-greek -target text

strongs-hebrew-gob:
	$(GORUN) ./cmd/morph -mode strongs-hebrew -target gob

strongs-hebrew:
	$(GORUN) ./cmd/morph -mode strongs-hebrew -target text

macula-hebrew-jsonl:
	$(GORUN) ./cmd/morph -mode macula-hebrew -target jsonl

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -o $(APP_NAME)-linux ./cmd/$(APP_NAME)

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -installsuffix cgo -v -ldflags '-w -s' -o $(APP_NAME).exe ./cmd/$(APP_NAME)

clean:
	rm -rf output/
	rm -f $(APP_NAME)-linux
	rm -f $(APP_NAME).exe
