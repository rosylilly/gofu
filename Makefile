.PHONEY: all clean run fmt

all:
	@go build -o bin/gofu

run:
	@go run *.go

fmt:
	@gofmt -tabs=false -tabwidth=2 -w -l *.go

clean:
	@rm -rf bin/gofu
