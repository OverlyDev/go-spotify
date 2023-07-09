
format:
	go fmt

run: format
	go run main.go

build: format
	go build -o build/ .

test: build
	./build/go-spotify-saver