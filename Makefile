.PHONY: build client

all: clean dep client build

dep:
	go get

build:
	go build

clean: clean-client
	rm -f acourse

clean-client:
	rm -rf public
	rm -rf templates

run:
	go run main.go

client: clean-client
	$(MAKE) -C client build
	mv public/static/* public/
	rm -rf public/static
