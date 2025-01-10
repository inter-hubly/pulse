SSH_KEY_PATH=/home/saimon/.ssh/id_ed25519_no_passphrase

build:
	@SSH_KEY=$(shell cat $(SSH_KEY_PATH) | base64 -w 0) && \
	docker build --build-arg SSH_KEY="$${SSH_KEY}" -t ghcr.io/inter-hubly/pulse:development . && \
	docker push ghcr.io/inter-hubly/pulse:development
