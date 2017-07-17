SERVICE=acourse
REGISTRY=gcr.io/acoshift-1362
COMMIT_SHA=$(shell git rev-parse HEAD)

default:
	# `make deploy` build and deploy to production
	# `make stag` build and deploy to staging
	# `make dev` starts server in localhost:8080
	# `make style` builds style

dev:
	go run -tags dev cmd/acourse/main.go

deploy: cluster patch

stag: cluster
	TAG=-dev make patch

# deploy: clean style build docker cluster patch

.PHONY: style
style:
	node build.js

clean:
	rm -f static/style.*.css
	rm -f entrypoint
	rm -f static.yaml

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o entrypoint -a -ldflags '-w -s' cmd/acourse/main.go

docker: clean style build
	gcloud docker -- build -t $(REGISTRY)/$(SERVICE):$(COMMIT_SHA) .
	gcloud docker -- push $(REGISTRY)/$(SERVICE):$(COMMIT_SHA)

cluster:
	gcloud container clusters get-credentials cluster-sg-1 --zone asia-southeast1-b --project acoshift-1362

patch:
	kubectl set image deployment/$(SERVICE)$(TAG) $(SERVICE)$(TAG)=$(REGISTRY)/$(SERVICE):$(COMMIT_SHA)
	kubectl rollout status deployment/$(SERVICE)$(TAG)
