
.PHONY: test
test: 
	mkdir -p build
	go test ./... -v -race -coverprofile=build/coverage.out