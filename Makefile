default:
	# `make dev` starts server in localhost:8080
	# `make style` builds style

dev:
	go run -tags dev cmd/acourse/main.go

.PHONY: style
style:
	node build.js

clean:
	rm -f static/style.*.css
	rm -f acourse
	rm -f static.yaml

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o acourse -a -ldflags '-w -s' cmd/acourse/main.go
