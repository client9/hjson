
all: install lint test

install:
	go install ./...

lint:
	golint ./...
	go vet ./...
	find . -name '*.go' | xargs gofmt -w -s

test:
	go test .
	find . -name '*' -type f | xargs misspell

clean:
	rm -f *~
	go clean ./...
	git gc

ci: install lint test

docker-ci:
	docker run --rm \
		-e COVERALLS_REPO_TOKEN=$COVERALLS_REPO_TOKEN \
		-v $(PWD):/go/src/github.com/client9/hjson \
		-w /go/src/github.com/client9/hjson \
		nickg/golang-dev-docker \
		make ci

.PHONY: ci docker-ci
