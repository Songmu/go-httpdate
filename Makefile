deps:
	go get -d -v ./...

test-deps:
	go get -d -v -t ./...

devel-deps:
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls

test: test-deps
	go test ./...

lint: devel-deps
	go vet ./...
	golint -set_exit_status ./...

cover: devel-deps
	goveralls

.PHONY: test deps lint cover
