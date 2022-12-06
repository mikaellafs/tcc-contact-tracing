Gerar stubs para grpc service:
<!-- go generate -->
protoc  src/grpc/proto/contact_tracing.proto --go-grpc_out=src/grpc --go_out=src/grpc

sudo apt install protobuf-compiler
sudo apt install golang-goprotobuf-dev