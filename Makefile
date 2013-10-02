.PHONEY: build clean

export CC = llvm-gcc

all: build

clean:
	@rm -rf bin

build:
	go build -o ./bin/gofu

format:
	@find **/*.go | xargs gofmt --tabs=false --tabwidth=2 -w -l

run: build
	@./bin/gofu
