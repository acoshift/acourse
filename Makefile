.PHONY: build ui dev

all: clean dep ui build

dep:
	go get

build:
	go build github.com/acoshift/acourse/cmd/acourse

clean: clean-ui
	rm -f acourse

clean-ui:
	rm -rf public
	rm -rf templates

run:
	go run main.go

dev:
	go run cmd/acourse-dev/main.go

ui: clean-ui
	$(MAKE) -C ui build
	mv public/static/* public/
	rm -rf public/static

local: clean-ui
	$(MAKE) -C ui local
	mv public/static/* public/
	rm -rf public/static

project:
	gcloud config set project acourse-d9d0a

clean-build:
	rm -rf build

build-server:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/acourse -a -ldflags '-s' github.com/acoshift/acourse/cmd/acourse

pre-build: dep ui
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
	./private/hook.sh

indexes: project
	gcloud preview datastore create-indexes index.yaml

cleanup-indexes: project
	gcloud preview datastore cleanup-indexes index.yaml
