Gerar stubs para grpc service:
<!-- go generate -->
protoc --go_out=src/grpc src/grpc/proto/contact_tracing.proto

sudo apt install protobuf-compiler
sudo apt install golang-goprotobuf-dev