# commands

GO=go
PROJECT=acourse-156413
# PROJECT=acourse-d9d0a

deploy: deploy-docker

deploy-docker: clean dep ui pre-build copy-dockerfile copy-ui config-prod build docker push

clean: clean-ui clean-build

dev:
	env CONFIG=private/config.stag.yaml $(GO) run cmd/acourse/main.go

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
	$(GO) get -v github.com/acoshift/acourse/cmd/acourse

re-dep:
	$(GO) get -u -v github.com/acoshift/acourse/cmd/acourse

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
	gcloud config set project $(PROJECT)

.PHONY: build
build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -o .build/acourse -a -ldflags '-s' github.com/acoshift/acourse/cmd/acourse

config-stag:
	cp private/config.stag.yaml .build/config.yaml

config-prod:
	cp private/config.prod.yaml .build/config.yaml

pre-build:
	mkdir -p .build
	curl https://curl.haxx.se/ca/cacert.pem > .build/cacert.pem
	cp Dockerfile .build/
	cp private/acourse_io.crt .build/
	cp private/acourse_io.key .build/

copy-dockerfile:
	cp Dockerfile .build/

copy-ui:
	mkdir -p .build
	cp -rf public .build/
	cp -rf templates .build/

docker:
	cd .build && docker build -t acourse .

push:
	docker tag acourse gcr.io/acourse-d9d0a/acourse
	gcloud docker -- push gcr.io/acourse-d9d0a/acourse

hook:
	./private/hook.sh

rolling-update:
	kubectl rolling-update acourse --image gcr.io/acourse-d9d0a/acourse --image-pull-policy Always

kube:
	gcloud container clusters get-credentials cluster-1 --zone asia-northeast1-a --project acourse-d9d0a
	kubectl proxy
