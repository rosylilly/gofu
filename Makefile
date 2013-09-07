.PHONEY: all clean

all:
	@go build -o bin/gofu

clean:
	@rm -rf bin/gofu
