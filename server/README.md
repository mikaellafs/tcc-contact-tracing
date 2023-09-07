Gerar stubs para grpc service:
<!-- go generate -->
protoc  pkg/grpc/proto/contact_tracing.proto --go-grpc_out=pkg/grpc --go_out=pkg/grpc

sudo apt install protobuf-compiler
sudo apt install golang-goprotobuf-dev