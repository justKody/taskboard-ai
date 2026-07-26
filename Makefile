.PHONY: build

build:
	go build -o build/taskboard-go-api cmd/main.go

run: build
	./build/taskboard-go-api

clean:
	rm -rf build


