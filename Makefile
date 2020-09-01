GO_LOCATION=/usr/local/go/bin
PREFIX_NAME=morph

.PHONY:

all: print json aws azure gcp mssql

linux: linux-print linux-json linux-aws linux-azure linux-gcp linux-mssql

windows: windows-print windows-json windows-aws windows-azure windows-gcp windows-mssql

print:
	GOOS=linux $(GO_LOCATION)/go build -tags print

json:
	GOOS=linux $(GO_LOCATION)/go build -tags json

aws:
	GOOS=linux $(GO_LOCATION)/go build -tags aws

azure:
	GOOS=linux $(GO_LOCATION)/go build -tags azure

gcp:
	GOOS=linux $(GO_LOCATION)/go build -tags gcp

mssql:
	GOOS=linux $(GO_LOCATION)/go build -tags mssql

linux-print:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags print -o $(PREFIX_NAME)-print

linux-json:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags json -o $(PREFIX_NAME)-json

linux-aws:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags aws -o $(PREFIX_NAME)-aws

linux-azure:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags azure -o $(PREFIX_NAME)-azure

linux-gcp:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags gcp -o $(PREFIX_NAME)-gcp

linux-mssql:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags mssql -o $(PREFIX_NAME)-mssql

windows-print:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags print -o $(PREFIX_NAME)-print.exe

windows-json:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags json -o $(PREFIX_NAME)-json.exe

windows-aws:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags aws -o $(PREFIX_NAME)-aws.exe

windows-azure:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags azure -o $(PREFIX_NAME)-azure.exe

windows-gcp:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags gcp -o $(PREFIX_NAME)-gcp.exe

windows-mssql:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_LOCATION)/go build -installsuffix cgo -v -ldflags '-w -s' -tags mssql -o $(PREFIX_NAME)-mssql.exe

clean:
	rm morph-* main
