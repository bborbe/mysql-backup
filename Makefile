REGISTRY ?= docker.io
IMAGE ?= bborbe/mysql-backup
ifeq ($(VERSION),)
	VERSION := $(shell git fetch --tags; git describe --tags `git rev-list --tags --max-count=1`)
endif

all: test install

install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install *.go

test:
	go test -cover -race $(shell go list ./... | grep -v /vendor/)

vet:
	go tool vet .
	go tool vet --shadow .

lint:
	golint -min_confidence 1 ./...

errcheck:
	errcheck -ignore '(Close|Write)' ./...

check: lint vet errcheck

goimports:
	go get golang.org/x/tools/cmd/goimports

format: goimports
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists
	go get -u github.com/golang/dep/cmd/dep

clean:
	docker rmi $(REGISTRY)/$(IMAGE)-build:$(VERSION)
	docker rmi $(REGISTRY)/$(IMAGE):$(VERSION)

buildgo:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o mysql-backup ./go/src/github.com/$(IMAGE)

.PHONY: check-version-tag
check-version-tag:
	@if [ -n "$(ALLOW_UNTAGGED_BUILD)" ]; then \
		echo "ALLOW_UNTAGGED_BUILD set — skipping version/tag check"; \
	else \
		head_tag=$$(git describe --tags --exact-match HEAD 2>/dev/null); \
		if [ "$$head_tag" != "$(VERSION)" ]; then \
			echo "ERROR: refusing to build $(VERSION) from this tree." >&2; \
			echo "  HEAD is at tag: $${head_tag:-<untagged>}" >&2; \
			echo "  building as:    $(VERSION)" >&2; \
			echo "  An image stamped vX.Y.Z must be built from the vX.Y.Z tag." >&2; \
			echo "  Fix: git checkout $(VERSION)   (or set ALLOW_UNTAGGED_BUILD=1 for a scratch build)" >&2; \
			exit 1; \
		fi; \
	fi

build: check-version-tag
	docker build --build-arg VERSION=$(VERSION) --no-cache --rm=true -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) -f ./Dockerfile.build .
	docker run -t $(REGISTRY)/$(IMAGE)-build:$(VERSION) /bin/true
	docker cp `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited`:/mysql-backup .
	docker rm `docker ps -q -n=1 -f ancestor=$(REGISTRY)/$(IMAGE)-build:$(VERSION) -f status=exited`
	docker build --no-cache --rm=true --tag=$(REGISTRY)/$(IMAGE):$(VERSION) -f Dockerfile.static .
	rm mysql-backup

upload:
	docker push $(REGISTRY)/$(IMAGE):$(VERSION)

trigger:
	@go get github.com/bborbe/docker-utils/cmd/docker-remote-tag-exists
	@exists=`docker-remote-tag-exists \
		-registry=${REGISTRY} \
		-repository="${IMAGE}" \
		-credentialsfromfile \
		-tag="${VERSION}" \
		-alsologtostderr \
		-v=0`; \
	trigger="build"; \
	if [ "$${exists}" = "true" ]; then \
		trigger="skip"; \
	fi; \
	echo $${trigger}
