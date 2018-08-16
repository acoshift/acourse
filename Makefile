default:
	# `make dev` starts server in localhost:8080
	# `make style` builds style

dev:
	goreload --all -x vendor

.PHONY: style
style:
	gulp

clean:
	rm -f static/style.*.css
	rm -f acourse
	rm -f static.yaml
	rm -rf .build

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o acourse -ldflags '-w -s' main.go
