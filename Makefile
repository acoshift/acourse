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

project:
	gcloud config set project acourse-d9d0a

clean-build:
	rm -rf build

build-server:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/acourse -a -ldflags '-s' main.go

pre-build: dep
	mkdir -p build
	curl https://curl.haxx.se/ca/cacert.pem > build/cacert.pem

deploy: clean-build pre-build build-server project
	cp -rf .private build/private
	cp -rf public build/
	cp -rf templates build/
	cp Dockerfile build/
	cp app.yaml build/
	# cd build && docker build -t acourse .
	cd build && gcloud app deploy
