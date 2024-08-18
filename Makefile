build:
	go build ./cmd/consumer

wire:
	go install github.com/google/wire/cmd/wire
	wire  ./internal/config/injector

wire_mock:
	go install github.com/google/wire/cmd/wire
	wire  ./test/bdd/injector

bdd: 
	go test -v ./test/bdd/steps/step_definitions_test.go