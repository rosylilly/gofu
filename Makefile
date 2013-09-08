.PHONEY: all clean

all:
	@go build -o bin/gofu

run: all
	@./bin/gofu

clean:
	@rm -rf bin/gofu
