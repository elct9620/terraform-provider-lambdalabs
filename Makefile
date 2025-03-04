default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	@TF_ACC=1 go test ./... -cover $(TESTARGS) -timeout 120m

install:
	go install .

doc:
	go generate ./...
