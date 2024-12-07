all: clean build run

build:
	go build -o bin/main .
run: 
	bin/main
clean:
	go mod tidy
	rm bin/* || true


