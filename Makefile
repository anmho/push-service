
.PHONY: air
air:
	@go install github.com/air-verse/air@latest


.PHONY: api
api:
	@cd lambda && go build -o ./bin/api ./cmd

.PHONY: deploy
deploy: api check
	@cdk deploy

.PHONY: staging
staging: api
	@cdk deploy --context environment=staging


.PHONY: watch
watch: air
	@air --build.cmd "cd lambda && go build -o ./bin/api ./cmd" --build.bin "./lambda/bin/api -local"

.PHONY: check
check:
	@cd lambda && go vet ./...
	@echo "synthesizing CloudFormation"
	@cdk synth

.PHONY: test
test:
	@cd lambda && go test ./...