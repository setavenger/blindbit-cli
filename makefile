generate:
	protoc --proto_path=lib/BlindBit-Protos lib/BlindBit-Protos/*.proto --go_out=lib --go-grpc_out=lib