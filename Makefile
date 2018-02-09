default:
	# `make dev` starts server in localhost:8000
	# `make style` builds style

dev:
	gin -p 8000 -a 8080 -x vendor --all -i

.PHONY: style
style:
	node build.js

clean:
	rm -f static/style.*.css
	rm -f acourse
	rm -f static.yaml

build:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o acourse -a -ldflags '-w -s' main.go
