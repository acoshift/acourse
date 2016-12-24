.PHONY: build client

all: clean dep client build

dep:
	go get

build:
	go build

clean:
	rm -rf public
	rm -f acourse

run:
	go run main.go

client:
	$(MAKE) -C client build
