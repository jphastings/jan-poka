package pb

//go:generate protoc --go_out=paths=source_relative,plugin=grpc:. config.proto deliveroo.proto
