
.PHONY: air
air:
	@go install github.com/air-verse/air@latest


.PHONY: api
api:
	@cd lambda && go build -o ./bin/api ./cmd

.PHONY: deploy
deploy: api


.PHONY: watch
watch: air
	@air --build.cmd "cd lambda && go build -o ./bin/api ./cmd" --build.bin "./lambda/bin/api -local"
