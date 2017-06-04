SERVICE=acourse-dev
REGISTRY=gcr.io/acoshift-1362
COMMIT_SHA=$(shell git rev-parse HEAD)

default:
	# `make deploy` build and deploy to production
	# `make stag` build and deploy to staging
	# `make dev` starts server in localhost:8080
	# `make style` builds style

dev:
	go run -tags dev cmd/acourse/main.go

stag:
	TAG=-dev make deploy

deploy: clean style build docker cluster patch

.PHONY: style
style: css = $(shell node_modules/.bin/node-sass --output-style compressed style/main.scss)
style: hash = $(shell echo "${css}" | md5)
style:
	@echo "${css}" > static/style.${hash}.css
	@echo "style.css: style.${hash}.css" >> static.yaml

clean:
	rm -f static/style.*.css
	rm -f entrypoint
	rm -f static.yaml

.PHONY: migrate
migrate:
	go run migrate/main.go

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o entrypoint -a -ldflags '-w -s' cmd/acourse/main.go

docker:
	gcloud docker -- build -t $(REGISTRY)/$(SERVICE) .
	docker tag $(REGISTRY)/$(SERVICE) $(REGISTRY)/$(SERVICE):$(COMMIT_SHA)
	gcloud docker -- push $(REGISTRY)/$(SERVICE):$(COMMIT_SHA)
	gcloud docker -- push $(REGISTRY)/$(SERVICE):latest

cluster:
	gcloud container clusters get-credentials cluster-sg-1 --zone asia-southeast1-b --project acoshift-1362

patch:
	kubectl set image deployment/$(SERVICE)$(TAG) $(SERVICE)$(TAG)=$(REGISTRY)/$(SERVICE):$(COMMIT_SHA)
	kubectl rollout status deployment/$(SERVICE)$(TAG)
