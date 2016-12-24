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
	cp -rf private build/
	cp -rf public build/
	cp -rf templates build/
	cp Dockerfile build/

deploy: clean-build pre-build build-server
	cd build && docker build -t acourse .
	docker tag acourse b.gcr.io/acoshift/acourse
	gcloud docker -- push b.gcr.io/acoshift/acourse
	curl http://52.77.156.5:8080/acourse.sh?key=XfwU58kx4Re5TGjck99CWV7DXrz27hYy

gae: clean-build pre-build build-server project
	cp app.yaml build/
	cd build && gcloud app deploy
