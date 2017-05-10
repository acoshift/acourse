default:
	# `make dev` starts server in localhost:8080
	# `make style` builds style

dev:
	go run -tags dev cmd/acourse/main.go

.PHONY: style
style: clean
	node_modules/.bin/node-sass --output-style compressed style/main.scss > static/style.css

clean:
	rm -f static/style.css
