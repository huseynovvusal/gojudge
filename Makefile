proto:
	protoc --go_out=. --go-grpc_out=. proto/executor.proto

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

clean:
	rm -f proto/*.pb.go