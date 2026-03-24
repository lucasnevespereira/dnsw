BINARY = dnsw

.PHONY: build run clean

build:
	go build -o $(BINARY) .

run: build
	go run main.go

clean:
	rm -f $(BINARY)
