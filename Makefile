build:
	go build ./cmd/consumer

wire:
	wire  ./internal/config/injector

bdd: 
	go install github.com/google/wire/cmd/wire
	go test -v ./test/bdd/steps/step_definitions_test.go