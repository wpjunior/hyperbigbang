.PHONY: setup run
setup:
	go get github.com/codegangsta/gin
	make deps

run:
	gin -b hyperbigbang run

include go.mk
