default:
	# `make dev` starts server in localhost:8000
	# `make style` builds style

dev:
	goreload --all -x vendor

.PHONY: style
style:
	node build.js

clean:
	rm -f static/style.*.css
	rm -f acourse
	rm -f static.yaml

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o acourse -ldflags '-w -s' main.go
