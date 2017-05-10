default:
	# `make dev` starts server in localhost:8080

dev:
	go run -tags dev cmd/acourse/main.go
