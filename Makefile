default:
	# `make dev` starts server in localhost:8080
	# `make style` builds style

dev:
	goreload \
		--all \
		-x node_modules

.PHONY: style
style:
	gulp

clean:
	rm -f static/style.*.css
	rm -f acourse
	rm -f static.yaml
	rm -rf .build

encrypt-deploy-token:
	gcloud kms encrypt \
		--project acourse-d9d0a \
		--location global \
		--keyring builder \
		--key key \
		--plaintext-file deploy-token \
		--ciphertext-file deploy-token.encrypted
