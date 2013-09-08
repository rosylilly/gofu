.PHONEY: all clean

all:
	@go build -o bin/gofu

run:
	@go run *.go

clean:
	@rm -rf bin/gofu
