# commands

deploy: clean dep ui pre-build config-prod build docker push hook

clean: clean-ui clean-build

dev:
	env CONFIG=private/config.stag.yaml go run cmd/acourse/main.go

indexes: project
	gcloud datastore create-indexes index.yaml

cleanup-indexes: project
	gcloud datastore cleanup-indexes index.yaml

.PHONY: proto
proto:
	protoc -I proto/ proto/acourse.proto --go_out=plugins=grpc:pkg/acourse

local: clean-ui
	$(MAKE) -C ui local
	mv public/static/* public/
	rm -rf public/static

dep:
	go get -v github.com/acoshift/acourse/cmd/acourse

re-dep:
	go get -u -v github.com/acoshift/acourse/cmd/acourse

# steps
# do not manually call step

clean-ui:
	rm -rf public
	rm -rf templates

clean-build:
	rm -rf .build

.PHONY: ui
ui: clean-ui
	$(MAKE) -C ui build
	mv public/static/* public/
	rm -rf public/static

project:
	gcloud config set project acourse-d9d0a

.PHONY: build
build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o .build/acourse -a -ldflags '-s' github.com/acoshift/acourse/cmd/acourse

config-stag:
	cp private/config.stag.yaml .build/

config-prod:
	cp private/config.prod.yaml .build/

pre-build:
	mkdir -p .build
	curl https://curl.haxx.se/ca/cacert.pem > .build/cacert.pem
	cp -rf public .build/
	cp -rf templates .build/
	cp Dockerfile .build/

docker:
	cd .build && docker build -t acourse .

push:
	docker tag acourse b.gcr.io/acoshift/acourse
	gcloud docker -- push b.gcr.io/acoshift/acourse

hook:
	./private/hook.sh
