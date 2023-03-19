.PHONY: eyeroute-embed
eyeroute-embed: front-build
	go build ./cmd/eyeroute

.PHONY: eyeroute
eyeroute:
	go build ./cmd/eyeroute

.PHONY: front-build
front-build:
	cd front && npm run build
