SRC := $(wildcard *.go)

music: $(SRC)
	go build
