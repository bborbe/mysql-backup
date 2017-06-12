all: test install run
install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/mysql_backup_cron/*.go
test:
	GO15VENDOREXPERIMENT=1 go test -cover `glide novendor`
vet:
	go tool vet .
	go tool vet --shadow .
lint:
	golint -min_confidence 1 ./...
errcheck:
	errcheck -ignore '(Close|Write)' ./...
check: lint vet errcheck
run:
	mysql_backup_cron \
	-logtostderr \
	-v=2 \
	-host=localhost \
	-port=5432 \
	-lock=/tmp/lock \
	-username=testuser \
	-password=S3CR3T \
	-database=testdatabase \
	-targetdir=/tmp \
	-name=mysql \
	-one-time
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	npm install
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install
update:
	glide up
clean:
	rm -rf vendor
