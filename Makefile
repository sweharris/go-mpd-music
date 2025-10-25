SRC := $(wildcard *.go)

music: $(SRC)
	go build -tags netgo
