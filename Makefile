default:
	# `make dev` starts server in localhost:8080
	# `make style` builds style

dev:
	goreload \
		--all \
		-x vendor \
		-x node_modules

.PHONY: style
style:
	gulp

clean:
	rm -f static/style.*.css
	rm -f acourse
	rm -f static.yaml
	rm -rf .build
