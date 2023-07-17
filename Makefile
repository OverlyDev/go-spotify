HELPER_PATH  = ./cmd/helper
HELPER_NAME  = go-spotify-helper
HELPER_BUILD = build/$(HELPER_NAME)

SAVER_PATH  = ./cmd/saver
SAVER_NAME  = go-spotify-saver
SAVER_BUILD = build/$(SAVER_NAME)

# General targets
clean:
	rm -rf build

format:
	go fmt ./...

build: format
	go build -o $(HELPER_BUILD)-linux-x64 $(HELPER_PATH)
	go build -o $(SAVER_BUILD)-linux-x64 $(SAVER_PATH)
	GOOS=windows GOARCH=amd64 go build -o $(HELPER_BUILD)-win-x64.exe $(HELPER_PATH)
	GOOS=windows GOARCH=amd64 go build -o $(SAVER_BUILD)-win-x64.exe $(SAVER_PATH)
	ls -lh ./build

# Helper specific targets
run_helper: format
	API_ID=$(shell cat settings/.clientid) API_SECRET=$(shell cat settings/.clientsecret) go run $(HELPER_PATH)

test_helper: build
	./$(HELPER_BUILD)

# Saver specific targets
run_saver: format
	go run $(SAVER_PATH)
	
test_saver:
	./$(SAVER_BUILD)