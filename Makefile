.PHONY: eyeroute-embed
eyeroute-embed: front-build
	go build ./cmd/eyeroute

.PHONY: eyeroute
eyeroute:
	go build ./cmd/eyeroute

.PHONY: eyeroute-cli
eyeroute-cli:
	go build ./cmd/eyeroute-cli

.PHONY: front-build
front-build:
	cd front && npm run build

.PHONY: buf-generate
buf-generate:
	rm -rf ./gen ./front/src/gen && buf generate
