

.PHONY: api
api:
	@go build -o api main

.PHONY: deploy
deploy: api