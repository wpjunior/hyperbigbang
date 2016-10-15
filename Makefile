.PHONY: setup run
setup:
	go get github.com/codegangsta/gin

run:
	gin -b hyperbigbang run

test:
	go test ./...
